package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	log "simu/src/logger"
	"simu/src/pkg/model"
)

//prereqquesties types and function
type apiHandler struct {
	simuClient model.SimuClient
}

func (h *apiHandler) SetClients(simulatotClient model.SimuClient) {
	h.simuClient = simulatotClient
}

//main functions

func (h *apiHandler) createMecHostHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var mecHost model.MecHost
	_ = json.NewDecoder(r.Body).Decode(&mecHost)
	log.Infof("Client tries to add new mecHost ID: %v, provider: %v\n", mecHost.Identity.Cluster, mecHost.Identity.Provider)
	if h.graphClient.CheckGraphContainsVertex(mecHost) {
		err := fmt.Errorf("Mec Host %v, %v not added beacuse it already exists", mecHost.Identity.Cluster, mecHost.Identity.Provider)
		log.Errorf(err.Error())
		w.WriteHeader(http.StatusConflict)
	} else {
		//this methods is only for create Mec Host, this is blocker for not to create list of edges at that time
		if containsAnyEdge(mecHost) {
			mecHost.Neighbours = nil
		}
		h.graphClient.AddMecHost(mecHost)
		w.WriteHeader(http.StatusOK)
	}
}

func (h *apiHandler) fetchCurrentFromNMT(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var apps []model.MECApp

}
