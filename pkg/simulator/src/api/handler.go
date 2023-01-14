package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
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

	for i := 0; i < experimentsNumber; i++ {

		//generate number of user to move
		id := h.generateUserToMove()

		// select new position for selected user
		h.generateTargetCellId(id)

		//update New Target Cell At user list (application list)

		//create smart placement intent

		//send request to ERC to select new position

		//generate request to orchestrator

		//send POST to orchestrator

		//update user list internally

	}
	w.WriteHeader(http.StatusOK)
}
