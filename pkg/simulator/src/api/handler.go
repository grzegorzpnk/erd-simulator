package api

import (
	"encoding/json"
	"net/http"
	log "simu/src/logger"
	"simu/src/pkg/model"
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

func (h *apiHandler) conductSingleExperiment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var intent ExperimentIntent

	err0 := json.NewDecoder(r.Body).Decode(&intent)
	if err0 != nil {
		log.Errorf("Cannot parse experiment intent. Error: %v", err0.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	status := executeExperiment(intent, h, 1)
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}

	//todo: get results from nmt and return as a response
	//Fetch stats from ERC
	//process results

	w.WriteHeader(http.StatusOK)
}

func (h *apiHandler) conductExperiment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var experimentDetails ExperimentDetails
	err0 := json.NewDecoder(r.Body).Decode(&experimentDetails)
	if err0 != nil {
		log.Errorf("Cannot parse experiment intent. Error: %v", err0.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var experiments []ExperimentIntent
	experiments = declareExperiments(experimentDetails)

	log.Infof("Started new full experiment with all 5 types")

	for z, experiment := range experiments {

		status := executeExperiment(experiment, h, z)
		if status != http.StatusOK {
			w.WriteHeader(status)
			return
		}

		//todo: get results from nmt and return as a response
		//Fetch stats from ERC
		//process results

	}

	log.Infof("Finished full experiment with all %v variations", len(experiments))
	w.WriteHeader(http.StatusOK)
}
