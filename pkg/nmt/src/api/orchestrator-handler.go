package api

import (
	log "10.254.188.33/matyspi5/erd/pkg/nmt/src/logger"
	mec_topology "10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/mec-topology"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (h *apiHandler) InstantiateApplication(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var mecApp model.MECApp
	_ = json.NewDecoder(r.Body).Decode(&mecApp)
	log.Infof("Client tries to deploy new mecApp ID: %v at MEC Host: %v \n", mecApp.Id, params["cluster"])

	var mecHost *model.MecHost
	mecHost = h.graphClient.GetMecHost(params["cluster"], "orange")

	if mecHost == nil {
		err := fmt.Errorf("MecHost %v does not exists", params["cluster"])
		log.Errorf(err.Error())
		w.WriteHeader(http.StatusConflict)
	}
	if !(mecHost.CheckEnoughResources(mecApp)) {
		err := fmt.Errorf("Mec App %v,cannot be instantiated beacuse of not enough resources", mecApp.Id)
		log.Errorf(err.Error())
		w.WriteHeader(http.StatusConflict)
	} else {
		err := mecHost.InstantiateApp(mecApp)
		if err != nil {
			log.Errorf("Error has been raised: ", err)
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

func (h *apiHandler) DeleteApplication(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var mecApp model.MECApp
	_ = json.NewDecoder(r.Body).Decode(&mecApp)
	log.Infof("Client tries to uninstall mecApp ID: %v from MEC Host: %v \n", mecApp.Id, params["cluster"])

	var mecHost *model.MecHost
	mecHost = h.graphClient.GetMecHost(params["cluster"], "orange")

	if !(mecHost.CheckAppExists(mecApp)) {
		err := fmt.Errorf("Mec App %v,cannot be deleted beacuse it does not exists", mecApp.Id)
		log.Errorf(err.Error())
		w.WriteHeader(http.StatusConflict)
	} else {
		mecHost.UninstallApp(mecApp)
		w.WriteHeader(http.StatusOK)
	}
}

/*
func (h *apiHandler) RelocateApplication(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var oldMecHost, newMecHost *model.MecHost
	oldMecHost = h.graphClient.GetMecHost(params["old-cluster"], "orange")
	newMecHost = h.graphClient.GetMecHost(params["new-cluster"], "orange")

	if oldMecHost.Identity.Cluster == newMecHost.Identity.Cluster {
		log.Infof("Application already on cluster - redundant relocation")
		w.WriteHeader(http.StatusOK)
		return
	}

	var mecApp model.MECApp
	_ = json.NewDecoder(r.Body).Decode(&mecApp)
	log.Infof("Client tries to RELOCATE mecApp ID: %v from MEC Host: %v to new MEC Host: %v \n", mecApp.Id, oldMecHost.Identity.Cluster, newMecHost.Identity.Cluster)

	if !(oldMecHost.CheckAppExists(mecApp)) {
		err := fmt.Errorf("Mec App %v,cannot be relocated beacuse it does not exists at source cluster", mecApp.Id)
		log.Errorf(err.Error())
		w.WriteHeader(http.StatusConflict)
		return
	}

	//instantiate
	if !(newMecHost.CheckEnoughResources(mecApp)) {
		err := fmt.Errorf("Mec App %v,cannot be instantiated beacuse of not enough resources", mecApp.Id)
		log.Errorf(err.Error())
		w.WriteHeader(http.StatusConflict)
		return
	} else {
		err := newMecHost.InstantiateApp(mecApp)
		if err != nil {
			log.Errorf("Error has been raised: ", err)
			w.WriteHeader(http.StatusConflict)
			return
		}
		oldMecHost.UninstallApp(mecApp)
		w.WriteHeader(http.StatusOK)
		log.Infof("App relocated succesfully ! \n")
	}
}*/

func (h *apiHandler) RelocateApplication2(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var oldMecHost, newMecHost *model.MecHost
	oldMecHost = h.graphClient.GetMecHost(params["old-cluster"], "orange")
	newMecHost = h.graphClient.GetMecHost(params["new-cluster"], "orange")

	if oldMecHost.Identity.Cluster == newMecHost.Identity.Cluster {
		log.Infof("Application already on cluster - redundant relocation")
		w.WriteHeader(http.StatusConflict)
		return
	}

	var mecApp model.MECApp
	_ = json.NewDecoder(r.Body).Decode(&mecApp)
	log.Infof("Client tries to RELOCATE mecApp ID: %v from MEC Host: %v to new MEC Host: %v \n", mecApp.Id, oldMecHost.Identity.Cluster, newMecHost.Identity.Cluster)
	log.Infof("Old mec resources: %v, %v", oldMecHost.CpuResources, oldMecHost.MemoryResources)
	log.Infof("New mec resources: %v, %v", newMecHost.CpuResources, newMecHost.MemoryResources)

	if !(oldMecHost.CheckAppExists(mecApp)) {
		err := fmt.Errorf("Mec App %v,cannot be relocated beacuse it does not exists at source cluster", mecApp.Id)
		log.Errorf(err.Error())
		w.WriteHeader(http.StatusConflict)
		return
	}

	//instantiate
	if !(newMecHost.CheckEnoughResources(mecApp)) {
		err := fmt.Errorf("Mec App %v,cannot be instantiated beacuse of not enough resources", mecApp.Id)
		log.Errorf(err.Error())
		w.WriteHeader(http.StatusConflict)
		return
	} else {
		err := newMecHost.InstantiateApp(mecApp)
		if err != nil {
			log.Errorf("Error has been raised: ", err)
			w.WriteHeader(http.StatusConflict)
			return
		}
		oldMecHost.UninstallApp(mecApp)

		//todo: riski point? how to handle if app will be relocated and not  find on a list?
		//update app clister on declaration list and ue positioning
		err2 := h.graphClient.UpdateAppCluster(mecApp, newMecHost)
		if err2 != nil {
			log.Errorf("Cannot find app on declaration list! cannot update cluster! %v", err.Error())
		}

		w.WriteHeader(http.StatusOK)
		log.Infof("App relocated succesfully ! \n")
		log.Infof("Old mec resources: %v, %v", oldMecHost.CpuResources, oldMecHost.MemoryResources)
		log.Infof("New mec resources: %v, %v", newMecHost.CpuResources, newMecHost.MemoryResources)
	}
}

func (h *apiHandler) DefineApplications(w http.ResponseWriter, r *http.Request) {
	var ac mec_topology.AppCounter

	h.graphClient.DeleteAllDeclaredApps()
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&ac)
	if err != nil {
		log.Errorf("Error: %v", err)
	}

	h.graphClient.DeclareApplications(ac)
	w.WriteHeader(http.StatusOK)
}

func (h *apiHandler) GenerateInitialClusters(w http.ResponseWriter, r *http.Request) {

	h.graphClient.UninstallAllApps()
	w.Header().Set("Content-Type", "application/json")

	status, _ := h.graphClient.FindInitialClusters()
	if status == true {
		w.WriteHeader(http.StatusOK)
		fmt.Printf("Apps with clusters:\n")
		for i := 0; i < len(h.graphClient.Application); i++ {
			h.graphClient.Application[i].PrintApplication()
		}
	} else {
		w.WriteHeader(http.StatusGone)
		log.Errorf("Cannot find clsuters for declared apps\n")
	}
}

func (h *apiHandler) InstantiateAllDefinedApps(w http.ResponseWriter, r *http.Request) {

	h.graphClient.UninstallAllApps()
	w.Header().Set("Content-Type", "application/json")
	err := h.graphClient.InstantiateAllDefinedApps()
	if err != nil {
		log.Errorf("Error has been raised: ", err)
		w.WriteHeader(http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *apiHandler) Prerequisites(w http.ResponseWriter, r *http.Request) {

	log.Infof("Prerequesties method")
	var ac mec_topology.AppCounter
	h.graphClient.DeleteAllDeclaredApps()
	h.graphClient.UninstallAllApps()

	//declare
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&ac)
	if err != nil {
		log.Errorf("Error: %v", err)
	}
	h.graphClient.DeclareApplications(ac)

	//find candidates mec and assign
	status, _ := h.graphClient.FindInitialClusters()
	if status == true {
		fmt.Printf("Apps with clusters:\n")
		for i := 0; i < len(h.graphClient.Application); i++ {
			var app model.MECApp
			app = *h.graphClient.Application[i]
			h.graphClient.ImmutableApplicationList = append(h.graphClient.ImmutableApplicationList, app)
			h.graphClient.ImmutableApplicationList[i].PrintApplication()
		}
	} else {
		w.WriteHeader(http.StatusConflict)
		log.Errorf("Cannot find cluster for declared apps\n")
		return
	}

	err = h.graphClient.InstantiateAllDefinedApps()
	if err != nil {
		log.Errorf("Error has been raised: ", err)
		w.WriteHeader(http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (h *apiHandler) PrerequisitesTunning(w http.ResponseWriter, r *http.Request) {

	var ac mec_topology.AppCounter

	w.Header().Set("Content-Type", "application/json")
	log.Infof("Started tunning max number of deployable applications.")

	appNumber := 44

	ac = mec_topology.AppCounter{Cg: 15, V2x: 15, Uav: 14}
	h.graphClient.DeclareApplications(ac)

	for x := 0; x < 50; x++ {

		appNumber += 1
		log.Infof("Started testing placement of %v applications", appNumber)
		//log.Infof("V2X: %v, CG: %v, UAV: %v", ac.V2x, ac.Cg, ac.Uav)

		//supportive flags
		alreadyFirstSuccess := false
		succCnt := 0
		firstSuccessCnt := 0

		for i := 0; i < 100; i++ {

			h.graphClient.UninstallAllApps()
			for i := 0; i < len(h.graphClient.Application); i++ {
				h.graphClient.Application[i].ClusterId = ""
				h.graphClient.Application[i].UserLocation = ""
			}

			//add additional app
			var app model.MECApp
			if appNumber%3 == 0 {
				app.Requirements.RequestedLatency = 10
			} else if appNumber%3 == 1 {
				app.Requirements.RequestedLatency = 15
			} else if appNumber%3 == 2 {
				app.Requirements.RequestedLatency = 30
			}

			app.Id = strconv.Itoa(len(h.graphClient.Application) + 1)
			app.GeneratreResourceRequirements()
			h.graphClient.Application = append(h.graphClient.Application, &app)

			status := h.graphClient.FindInitialClustersTunning()
			if status == true {
				succCnt++
				if !alreadyFirstSuccess {
					firstSuccessCnt = i + 1
					alreadyFirstSuccess = true
				}
			}
		}

		log.Infof("%v. Placement of %v apps, has been identified %v / 100 times!. First sucessfull search at: %v attempts!", x, appNumber, succCnt, firstSuccessCnt)
		if succCnt == 0 {
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (h *apiHandler) Recreate(w http.ResponseWriter, r *http.Request) {

	//declare
	w.Header().Set("Content-Type", "application/json")

	log.Infof("Received request to recreate initial app placement")
	//first let's remove all apps from clusters
	h.graphClient.UninstallAllApps()
	h.graphClient.Application = nil

	for i := 0; i < len(h.graphClient.ImmutableApplicationList); i++ {
		app := h.graphClient.ImmutableApplicationList[i]
		h.graphClient.Application = append(h.graphClient.Application, &app)
	}

	err := h.graphClient.InstantiateAllDefinedApps()
	if err != nil {
		log.Errorf("Error has been raised: ", err)
		w.WriteHeader(http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusOK)

}
