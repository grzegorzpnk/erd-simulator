package api

import (
	log "10.254.188.33/matyspi5/erd/pkg/nmt/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
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

	h.graphClient.DeleteAllDeclaredApps()
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	h.graphClient.DeclareApplications(params["applications"])
	w.WriteHeader(http.StatusOK)
}

func (h *apiHandler) GenerateInitialClusters(w http.ResponseWriter, r *http.Request) {

	h.graphClient.UninstallAllApps()
	w.Header().Set("Content-Type", "application/json")

	status := h.graphClient.FindInitialClusters()
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

func (h *apiHandler) Prerequesties(w http.ResponseWriter, r *http.Request) {

	h.graphClient.DeleteAllDeclaredApps()
	h.graphClient.UninstallAllApps()

	//declare
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	h.graphClient.DeclareApplications(params["applications"])

	//find candidates mec and assign
	status := h.graphClient.FindInitialClusters()
	if status == true {
		fmt.Printf("Apps with clusters:\n")
		for i := 0; i < len(h.graphClient.Application); i++ {
			h.graphClient.Application[i].PrintApplication()
		}
	} else {
		w.WriteHeader(http.StatusConflict)
		log.Errorf("Cannot find clsuters for declared apps\n")
		return
	}

	err := h.graphClient.InstantiateAllDefinedApps()
	if err != nil {
		log.Errorf("Error has been raised: ", err)
		w.WriteHeader(http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusOK)

}
