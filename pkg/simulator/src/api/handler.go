package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	log "simu/src/logger"
	"simu/src/pkg/model"
	"simu/src/pkg/results"
	"strconv"
)

type apiHandler struct {
	SimuClient   model.SimuClient
	ResultClient results.Client
}

func (h *apiHandler) SetClients(sClient model.SimuClient, rClient results.Client) {
	h.SimuClient = sClient
	h.ResultClient = rClient
}

func (h *apiHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.SimuClient.Apps)
	w.WriteHeader(http.StatusOK)
}

// getAllResults
func (h *apiHandler) getAllResults(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	respBody := h.ResultClient.GetResults()

	err := json.NewEncoder(w).Encode(respBody)
	if err != nil {
		log.Errorf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *apiHandler) generateChartPkg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	basePath := "results"
	iter := "iterations"
	aggr := "aggregated"
	mecs := "mec-hosts"

	err := h.ResultClient.GenerateChartPkgApps(results.RelocationRates, basePath+"/"+iter)
	if err != nil {
		log.Errorf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.ResultClient.GenerateChartPkgApps(results.SkippedRates, basePath+"/"+iter)
	if err != nil {
		log.Errorf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.ResultClient.GenerateChartPkgApps(results.RedundantRates, basePath+"/"+iter)
	if err != nil {
		log.Errorf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.ResultClient.GenerateChartPkgApps(results.FailedRates, basePath+"/"+iter)
	if err != nil {
		log.Errorf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.ResultClient.GenerateChartPkgApps(results.RelocationTriggeringRates, basePath+"/"+aggr)
	if err != nil {
		log.Errorf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.ResultClient.GenerateChartPkgApps(results.RelocationSuccessfulSearchRates, basePath+"/"+aggr)
	if err != nil {
		log.Errorf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.ResultClient.GenerateChartPkgApps(results.RelocationRejectionRates, basePath+"/"+aggr)
	if err != nil {
		log.Errorf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.ResultClient.GenerateChartPkgMecs(results.ResCpu, basePath+"/"+mecs)
	if err != nil {
		log.Errorf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.ResultClient.GenerateChartPkgMecs(results.ResMemory, basePath+"/"+mecs)
	if err != nil {
		log.Errorf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

	strategy := specifyStrategy(intent.Weights)

	log.Infof("Starting Experiment [%v] type: %v, strategy: %v ", 1, intent.ExperimentType, strategy)
	experimentType, err := checkExperimentType(intent.ExperimentType)
	if err != nil {
		log.Errorf("Could not proceed with experiment. Reason: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	appsNumber := intent.ExperimentDetails.AppsNumber

	if err != nil {
		log.Errorf("Could not proceed with experiment. Reason: [apps-number] %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	experimentIterations, err := strconv.Atoi(intent.ExperimentDetails.ExperimentIterations)
	if err != nil {
		log.Errorf("Could not proceed with experiment. Reason: [experiments-iterations] %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Infof("Started new experiment: %v, with %v relocations", experimentType, experimentIterations)

	//at the beggining let's synchro latest placement at nmt
	err = GenerateInitialAppPlacementAtNMT(intent.ExperimentDetails.AppsNumber)
	if err != nil {
		log.Errorf("Cannot make initial placement of app at NMT. Error: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		log.Infof("NMT has just randomly deployed %v apps. NMT ready to start experiment", intent.ExperimentDetails.AppsNumber.GetTotalAsString())
	}

	err = h.SimuClient.FetchAppsFromNMT()
	if err != nil {
		log.Errorf("Cannot fetch current app list from NMT. Error: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		log.Infof("Initial app list fetched from NMT")
	}

	err = resetResultsAtERC()
	if err != nil {
		log.Errorf("Cannot reset the results at NMT. Error: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		log.Infof("Results module ready -> cache cleared at NMT")
	}

	//loop for each experimentIterations defined in intent
	for i := 0; i < experimentIterations; i++ {
		status := executeExperiment(intent, h, 1, i)
		if status != true {
			log.Error("Experiment cannot be coninued due to error in one of the iterations")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = h.ResultClient.CollectExperimentStats(experimentType, "null", appsNumber, experimentIterations)
		if err != nil {
			log.Errorf("Error: %v. Status code: %v", err, http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	h.ResultClient.IncExpId()

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

	fmt.Printf("%+v, total: %v", experimentDetails, experimentDetails.AppsNumber.GetTotal())

	var experiments []ExperimentIntent

	experiments = declareExperiments(experimentDetails)

	log.Infof("Started new full experiment with all 5 types")

	//loop for each experiment defined in method declareExperiments()
	for z, experiment := range experiments {

		experimentType, err := checkExperimentType(experiment.ExperimentType)
		if err != nil {
			log.Errorf("Could not proceed with experiment. Reason: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		strategy := specifyStrategy(experiment.Weights)

		log.Infof("Experiment [%v] type: %v, strategy: %v ", z+1, experiment.ExperimentType, strategy)

		appsNumber := experiment.ExperimentDetails.AppsNumber

		if err != nil {
			log.Errorf("Could not proceed with experiment. Reason: [apps-number] %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		experimentIterations, err := strconv.Atoi(experiment.ExperimentDetails.ExperimentIterations)
		if err != nil {
			log.Errorf("Could not proceed with experiment. Reason: [experiments-iterations] %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//at the beggining let's synchro latest placement at nmt
		err = GenerateInitialAppPlacementAtNMT(experiment.ExperimentDetails.AppsNumber)
		if err != nil {
			log.Errorf("Cannot make initial placement of app at NMT. Error: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			log.Infof("NMT has just randomly deployed %v apps. NMT ready to start experiment", experiment.ExperimentDetails.AppsNumber.GetTotalAsString())
		}

		err = h.SimuClient.FetchAppsFromNMT()
		if err != nil {
			log.Errorf("Cannot fetch current app list from NMT. Error: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			log.Infof("Initial app list fetched from NMT")
		}

		err = resetResultsAtERC()
		if err != nil {
			log.Errorf("Cannot reset the results at NMT. Error: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			log.Infof("Results module ready -> cache cleared at NMT")
		}

		//loop for each sub-experiment defined in method declareExperiments()
		for i := 0; i < experimentIterations; i++ {
			status := executeExperiment(experiment, h, z, i)
			if status != true {
				log.Error("Experiment cannot be coninued due to error in one of the iterations, skip this and let's go to next experiment")
				break
			}

		}

		err = h.ResultClient.CollectExperimentStats(experimentType, strategy, appsNumber, experimentIterations)
		if err != nil {
			log.Errorf("Error: %v. Status code: %v", err, http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		h.ResultClient.IncExpId()
		log.Infof("Finished Experiment [%v] type: %v, strategy: %v ", z+1, experiment.ExperimentType, specifyStrategy(experiment.Weights))
	}
	log.Infof("Finished all Experiments (%v), each %v iterations", len(experiments), experiments[0].ExperimentDetails.ExperimentIterations)
	w.WriteHeader(http.StatusOK)
}
