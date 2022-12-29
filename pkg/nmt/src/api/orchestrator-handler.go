package api

import (
	log "10.254.188.33/matyspi5/erd/pkg/nmt/src/logger"
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
		mecHost.InstantiateApp(mecApp)
		w.WriteHeader(http.StatusOK)
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
		newMecHost.InstantiateApp(mecApp)
		oldMecHost.UninstallApp(mecApp)
		w.WriteHeader(http.StatusOK)
		log.Infof("App relocated succesfully ! \n")
	}
}

func (h *apiHandler) GenerateInitialClusters(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	/*params := mux.Vars(r)
	appNumber := params["app-number"]
	h.graphClient.DeclareApplications(appNumber)*/

	var mecHostSource []model.MecHost
	var cells = map[int]int{}

	//PREREQUESTIES
	//in order not to work on real clusters ( cause we need to update resources, etc for looking for initial configutration)
	for _, v := range h.graphClient.MecHosts {
		mecHostSource = append(mecHostSource, *v)
	}
	//another copy for each serach iteration ( if one search will fail, wee need to repeat on still fresh data)
	var mecHostsSourcesTmp = make([]model.MecHost, len(mecHostSource))

	var cnt = 0
	search := true
	for search {
		fmt.Println("[DEBUG] Starting search.")
		cnt++
		if cnt > 5 {
			fmt.Printf("Cannot identify initial clusters!\n")
			w.WriteHeader(http.StatusGone)
			return
		}
		copy(mecHostsSourcesTmp, mecHostSource)
		cells = generateRandomCells()

		search = false
		for index, edgeApp := range h.graphClient.Application {
			startCell := h.graphClient.GetCell(strconv.Itoa(cells[index+1]))
			cmh, err := findCanidateMec(*edgeApp, startCell, mecHostsSourcesTmp, &h.graphClient)
			if err != nil {
				search = true
				fmt.Printf("Could not find candidate mec for App[%v]. Search failed.\n", edgeApp.Id)
				break
			}
			mecHostsSourcesTmp = updateMecResourcesInfo(mecHostsSourcesTmp, cmh, *edgeApp)
			edgeApp.ClusterId = cmh.Identity.Cluster
		}
	}

	fmt.Printf("Found after %v iterations", cnt)
	w.WriteHeader(http.StatusOK)

	fmt.Printf("Apps with clusters:\n")
	for i := 0; i < len(h.graphClient.Application); i++ {
		h.graphClient.Application[i].PrintApplication()
	}

	//TODO: instantiate apps

}

func (h *apiHandler) OnboardApplications(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	h.graphClient.DeclareApplications(params["applications"])
	w.WriteHeader(http.StatusOK)
}
