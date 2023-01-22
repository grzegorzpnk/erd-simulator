package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	log "simu/src/logger"
	"simu/src/pkg/model"
	"strconv"
)

type apiHandler struct {
	SimuClient model.SimuClient
}

func (h *apiHandler) SetClients(simulatotClient model.SimuClient) {
	h.SimuClient = simulatotClient
}

func (h *apiHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.SimuClient.Apps)
	w.WriteHeader(http.StatusOK)
}

func (h *apiHandler) conductExperiment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	experimentsNumber, _ := strconv.Atoi(params["mobility-number"])

	var weights model.Weights
	_ = json.NewDecoder(r.Body).Decode(&weights)
	log.Infof("")
	log.Infof("Started new experiment, with %v relocations", experimentsNumber)

	for i := 0; i < experimentsNumber; i++ {

		//log.Infof("Experiment numer: %v", i+1)

		experimentN := "[EXPERIMENT " + strconv.Itoa(i+1) + "] "
		//generate number of user to move
		//id := h.generateUserToMove()
		id := "10"

		// select new position for selected user and add new position to UserPath
		app := h.SimuClient.GetApps(id)
		h.generateTargetCellId(app)
		log.Infof(experimentN+"User(app) with ID: %v [current mec: %v] moved FROM cell: %v, towards cell: %v", app.Id, app.ClusterId, app.UserPath[len(app.UserPath)-2], app.UserLocation)

		//create smart placement intent

		spi, err := GenerateSmartPlacementIntent(*app, weights)
		if err != nil {
			log.Errorf("Cannot generate SPI: %v", err.Error())
		}

		//send request to ERC to select new position
		cluster, err := CallPlacementController(spi)
		if err != nil {
			log.Errorf("Call Placement ctrl has failed : %v", err.Error())
		}

		if cluster.Cluster == app.ClusterId {
			log.Infof(experimentN+"Selected redundant cluster: %v -> missing relocation", cluster.Cluster)
		} else {
			log.Infof(experimentN+"Selected new cluster: %v", cluster.Cluster)

			//generate request to orchestrator
			err2 := sendRelocationRequest(*app, *cluster)
			if err2 != nil {
				log.Errorf("Cannot relocate app! Error: %v", err2.Error())
			} else {
				log.Infof(experimentN + "Application has been relocated in nmt")
			}

			//update cluster in app list internallyt
			app.ClusterId = cluster.Cluster
		}
	}
	w.WriteHeader(http.StatusOK)
}
