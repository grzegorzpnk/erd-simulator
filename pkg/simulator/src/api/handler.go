package api

import (
	"encoding/json"
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

func (h *apiHandler) conductSingleExperiment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var intent ExperimentIntent

	err0 := json.NewDecoder(r.Body).Decode(&intent)
	if err0 != nil {
		log.Errorf("Cannot parse experiment intent. Error: %v", err0.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Infof("Starting Experiment [%v] type: %v, strategy: %v ", 1, intent.ExperimentType, specifyStrategy(intent.Weights))

	//at the beggining let's synchro latest placement at nmt
	err := GenerateInitialAppPlacementAtNMT(intent.ExperimentDetails.AppNumber)
	if err != nil {
		log.Errorf("Cannot make initial placement of app at NMT. Error: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		log.Infof("NMT has just randomly deployed %v apps. NMT ready to start experiment", intent.ExperimentDetails.AppNumber)
	}

	err2 := h.SimuClient.FetchAppsFromNMT()
	if err2 != nil {
		log.Errorf("Cannot fetch current app list from NMT. Error: %v", err2.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		log.Infof("Initial app list fetched from NMT")
	}

	err3 := resetResultsAtERC()
	if err3 != nil {
		log.Errorf("Cannot reset the results at NMT. Error: %v", err3.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		log.Infof("Results module ready -> cache cleared at NMT")
	}

	//loop for each sub-experiment defined in method declareExperiments()
	experimentIterations, _ := strconv.Atoi(intent.ExperimentDetails.ExperimentIterations)

	for i := 0; i < experimentIterations; i++ {
		status := executeExperiment(intent, h, 1, i)
		if status != true {
			log.Error("Experiment cannot be coninued due to error in one of the iterations")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
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

	//loop for each experiment defined in method declareExperiments()
	for z, experiment := range experiments {

		log.Infof("Experiment [%v] type: %v, strategy: %v ", z+1, experiment.ExperimentType, specifyStrategy(experiment.Weights))

		//at the beggining let's synchro latest placement at nmt
		err := GenerateInitialAppPlacementAtNMT(experiment.ExperimentDetails.AppNumber)
		if err != nil {
			log.Errorf("Cannot make initial placement of app at NMT. Error: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			log.Infof("NMT has just randomly deployed %v apps. NMT ready to start experiment", experiment.ExperimentDetails.AppNumber)
		}

		err2 := h.SimuClient.FetchAppsFromNMT()
		if err2 != nil {
			log.Errorf("Cannot fetch current app list from NMT. Error: %v", err2.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			log.Infof("Initial app list fetched from NMT")
		}

		err3 := resetResultsAtERC()
		if err3 != nil {
			log.Errorf("Cannot reset the results at NMT. Error: %v", err3.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			log.Infof("Results module ready -> cache cleared at NMT")
		}

		//loop for each sub-experiment defined in method declareExperiments()
		experimentIterations, _ := strconv.Atoi(experiment.ExperimentDetails.ExperimentIterations)
		for i := 0; i < experimentIterations; i++ {
			status := executeExperiment(experiment, h, z, i)
			if status != true {
				break
				log.Error("Experiment cannot be coninued due to error in one of the iterations, skip this and let's go to next experiment")
			}
		}

		//todo: get results from nmt and return as a response
		//Fetch stats from ERC
		//process results

		log.Infof("Finished Experiment [%v] type: %v, strategy: %v ", z+1, experiment.ExperimentType, specifyStrategy(experiment.Weights))

	}

	log.Infof("Finished all Experiments (%v), each %v iterations", len(experiments), experiments[0].ExperimentDetails.ExperimentIterations)
	w.WriteHeader(http.StatusOK)
}
