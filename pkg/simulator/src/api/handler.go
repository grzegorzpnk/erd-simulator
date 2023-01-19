package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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

	for i := 0; i < experimentsNumber; i++ {

		//generate number of user to move
		id := h.generateUserToMove()

		// select new position for selected user and add new position to UserPath
		app := h.SimuClient.GetApps(id)
		h.generateTargetCellId(app)

		//create smart placement intent

		spi, err := GenerateSmartPlacementIntent(*app, weights)
		if err != nil {
			log.Fatal("Cannot generate SPI: %v", err.Error())
		}

		//send request to ERC to select new position
		cluster, err := CallPlacementController(spi)
		if err != nil {
			log.Fatal("Call Placement ctrl has failed : %v", err.Error())
		}

		//generate request to orchestrator
		err2 := sendRelocationRequest(*app, *cluster)
		if err2 != nil {
			log.Fatal("Cannot relocate app! Error: %v", err2.Error())
		}

		//update cluster in app list internallyt
		app.ClusterId = cluster.Cluster

	}
	w.WriteHeader(http.StatusOK)
}
