package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
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
	apps := "apps"
	mecs := "mecs"

	err := h.ResultClient.GenerateChartPkgApps(results.RelocationTriggeringRates, basePath+"/"+apps)
	if err != nil {
		log.Errorf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.ResultClient.GenerateChartPkgApps(results.RelocationSuccessfulSearchRates, basePath+"/"+apps)
	if err != nil {
		log.Errorf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.ResultClient.GenerateChartPkgApps(results.RelocationRejectionRates, basePath+"/"+apps)
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

func (h *apiHandler) generateICCHeuristicChart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	basePath := "results"
	apps := "apps"
	mecs := "mecs"

	err := h.ResultClient.GenerateChartPkgAppsICC(results.RelocationTriggeringRates, basePath+"/"+apps)
	if err != nil {
		log.Errorf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.ResultClient.GenerateChartPkgAppsICC(results.RelocationRejectionRates, basePath+"/"+apps)
	if err != nil {
		log.Errorf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.ResultClient.GenerateChartPkgMecsICC(results.ResCpu, basePath+"/"+mecs)
	if err != nil {
		log.Errorf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.ResultClient.GenerateChartPkgMecsICC(results.ResMemory, basePath+"/"+mecs)
	if err != nil {
		log.Errorf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *apiHandler) conductSingleExperiment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var expIntent model.ExperimentIntent

	err0 := json.NewDecoder(r.Body).Decode(&expIntent)
	if err0 != nil {
		log.Errorf("Cannot parse experiment intent. Error: %v", err0.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !reflect.ValueOf(expIntent).FieldByName("Weights").IsZero() {
		strategy := specifyStrategy(expIntent.Weights)
		log.Infof("Starting Experiment [%v] type: %v, strategy: %v ", 1, expIntent.ExperimentType, strategy)
	}

	log.Infof("Starting Experiment [%v] type: %v", 1, expIntent.ExperimentType)
	err := checkExperimentType(expIntent.ExperimentType)
	if err != nil {
		log.Errorf("Could not proceed with experiment. Reason: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Errorf("Could not proceed with experiment. Reason: [apps-number] %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	movementsInExperiment, err := strconv.Atoi(expIntent.ExperimentDetails.MovementsInExperiment)
	if err != nil {
		log.Errorf("Could not proceed with experiment. Reason: [experiments-iterations] %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Infof("Started new experiment: [T:%v][S:%v], with %v relocations", expIntent.ExperimentType, expIntent.ExperimentStrategy, movementsInExperiment)

	//at the beggining let's synchro latest placement at nmt
	err = GenerateInitialAppPlacementAtNMT(expIntent.ExperimentDetails.InitialAppsNumber)
	if err != nil {
		log.Errorf("Cannot make initial placement of app at NMT. Error: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		log.Infof("NMT has just randomly deployed %v apps. NMT ready to start experiment", expIntent.ExperimentDetails.InitialAppsNumber.GetTotalAsString())
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

	//loop for each movementsInExperiment defined in intent
	for i := 0; i < movementsInExperiment; i++ {
		status := h.executeExperiment(expIntent, 1, i)
		if status != true {
			log.Error("Experiment cannot be coninued due to error in one of the iterations")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	err = h.ResultClient.CollectExperimentStats(expIntent)
	if err != nil {
		log.Errorf("Error: %v. Status code: %v", err, http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.ResultClient.IncExpId()
	w.WriteHeader(http.StatusOK)
}

func (h *apiHandler) conductExperiment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var experimentDetails model.ExperimentDetails
	err0 := json.NewDecoder(r.Body).Decode(&experimentDetails)
	if err0 != nil {
		log.Errorf("Cannot parse experiment intent. Error: %v", err0.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("%+v, total: %v", experimentDetails, experimentDetails.InitialAppsNumber.GetTotal())

	var experiments []model.ExperimentIntent

	experiments = declareExperiments(experimentDetails)

	log.Infof("Started new full experiment with all 5 types")

	//loop for each experiment defined in method declareExperiments()
	for z, experiment := range experiments {

		err := checkExperimentType(experiment.ExperimentType)
		if err != nil {
			log.Errorf("Could not proceed with experiment. Reason: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		strategy := specifyStrategy(experiment.Weights)

		log.Infof("Experiment [%v] type: %v, strategy: %v ", z+1, experiment.ExperimentType, strategy)

		if err != nil {
			log.Errorf("Could not proceed with experiment. Reason: [apps-number] %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		experimentIterations, err := strconv.Atoi(experiment.ExperimentDetails.MovementsInExperiment)
		if err != nil {
			log.Errorf("Could not proceed with experiment. Reason: [experiments-iterations] %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//at the beggining let's synchro latest placement at nmt
		err = GenerateInitialAppPlacementAtNMT(experiment.ExperimentDetails.InitialAppsNumber)
		if err != nil {
			log.Errorf("Cannot make initial placement of app at NMT. Error: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			log.Infof("NMT has just randomly deployed %v apps. NMT ready to start experiment", experiment.ExperimentDetails.InitialAppsNumber.GetTotalAsString())
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
			status := h.executeExperiment(experiment, z, i)
			if status != true {
				log.Error("Experiment cannot be coninued due to error in one of the iterations, skip this and let's go to next experiment")
				break
			}

		}

		err = h.ResultClient.CollectExperimentStats(experiment)
		if err != nil {
			log.Errorf("Error: %v. Status code: %v", err, http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		h.ResultClient.IncExpId()
		log.Infof("Finished Experiment [%v] type: %v, strategy: %v ", z+1, experiment.ExperimentType, specifyStrategy(experiment.Weights))
	}
	log.Infof("Finished all Experiments (%v), each %v iterations", len(experiments), experiments[0].ExperimentDetails.MovementsInExperiment)
	w.WriteHeader(http.StatusOK)
}

func (h *apiHandler) conductExperimentGlobcom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var experimentDetails model.ExperimentDetails
	err0 := json.NewDecoder(r.Body).Decode(&experimentDetails)
	if err0 != nil {
		log.Errorf("Cannot parse experiment intent. Error: %v", err0.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var experiments []model.ExperimentIntent
	experiments = declareGlobcomExperiments(experimentDetails)
	log.Infof("Started new full GLOBECOM experiment with all 4 types: Optimal, EAR, RL-masked, RL-no-masked")

	movements, err := strconv.Atoi(experiments[0].ExperimentDetails.MovementsInExperiment)
	if err != nil {
		log.Errorf("Could not proceed with experiment. Reason: [movements-in-experiment] %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//in order to keep the same settings for each of experiment, let's generate common trajectory, that each of experiment will be invoked on
	err = GenerateInitialAppPlacementAtNMT(experiments[0].ExperimentDetails.InitialAppsNumber)
	if err != nil {
		log.Errorf("Cannot make initial placement of app at NMT. Error: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		log.Infof("NMT has just randomly deployed %v apps. NMT ready to start experiment", experiments[0].ExperimentDetails.InitialAppsNumber.GetTotalAsString())
	}

	err = h.SimuClient.FetchAppsFromNMT()
	if err != nil {
		log.Errorf("Cannot fetch current app list from NMT. Error: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		log.Infof("Initial app list fetched from NMT")
		for i := 0; i < len(h.SimuClient.Apps); i++ {
			log.Infof("1. Apps[%v]: %v, cluster: %v", h.SimuClient.Apps[i].Id, h.SimuClient.Apps[i].UserPath, h.SimuClient.Apps[i].ClusterId)
		}
	}

	trajectory, err := createTrajectory(movements, h)
	if err != nil {
		log.Errorf("Cannot create trajectory. Error: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		log.Infof("Trajectory has been created: %v", trajectory)
		for i := 0; i < len(h.SimuClient.Apps); i++ {
			log.Infof("2. Apps[%v] path: %v, cluster: %v ", h.SimuClient.Apps[i].Id, h.SimuClient.Apps[i].UserPath, h.SimuClient.Apps[i].ClusterId)
		}
	}

	//loop for each experiment defined in method declareExperiments()
	for z, experiment := range experiments {

		err := checkExperimentType(experiment.ExperimentType)
		if err != nil {
			log.Errorf("Could not proceed with experiment. Reason: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Infof("Experiment [%v] type: %v", z+1, experiment.ExperimentType)

		//at the beggining let's recreate initial app placement at NMT and fetch
		err = h.SimuClient.RecreateInitialPlacementAtNMT()
		if err != nil {
			log.Errorf("Cannot recreate initial placement and fetch current app list from NMT. Error: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			log.Infof("Initial placement recreated at NMT")
		}

		err = h.SimuClient.FetchAppsFromNMT()
		if err != nil {
			log.Errorf("Cannot recreate initial placement and fetch current app list from NMT. Error: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			log.Infof("Initial app list fetched from NMT")
		}
		for i := 0; i < len(h.SimuClient.Apps); i++ {
			log.Infof("3. Apps[%v]: %v, cluster: %v ", h.SimuClient.Apps[i].Id, h.SimuClient.Apps[i].UserPath, h.SimuClient.Apps[i].ClusterId)
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
		for i := 0; i < len(trajectory); i++ {
			status := h.executeGlobcomExperiment(experiment, z, i, trajectory[i][0], trajectory[i][1])
			if status != true {
				log.Error("Experiment cannot be coninued due to error in one of the iterations, skip this and let's go to next experiment")
				break
			}

		}

		err = h.ResultClient.CollectExperimentStats(experiment)
		if err != nil {
			log.Errorf("Error: %v. Status code: %v", err, http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		h.ResultClient.IncExpId()
		log.Infof("Finished Experiment [%v] type: %v", z+1, experiment.ExperimentType)
	}
	log.Infof("Finished all Experiments (%v), each %v iterations", len(experiments), experiments[0].ExperimentDetails.MovementsInExperiment)
	w.WriteHeader(http.StatusOK)
}

func (h *apiHandler) conductExperimentICC(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var experimentDetails model.ExperimentDetails
	err0 := json.NewDecoder(r.Body).Decode(&experimentDetails)
	if err0 != nil {
		log.Errorf("Cannot parse experiment intent. Error: %v", err0.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var experiments []model.ExperimentIntent
	experiments = declarePhDExperiments(experimentDetails)
	log.Infof("Started new full PhD experiment with all %v types.", len(experiments))

	movements, err := strconv.Atoi(experiments[0].ExperimentDetails.MovementsInExperiment)
	if err != nil {
		log.Errorf("Could not proceed with experiment. Reason: [movements-in-experiment] %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//in order to keep the same settings for each of experiment, let's generate common trajectory, that each of experiment will be invoked on
	err = GenerateInitialAppPlacementAtNMT(experiments[0].ExperimentDetails.InitialAppsNumber)
	if err != nil {
		log.Errorf("Cannot make initial placement of app at NMT. Error: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		log.Infof("NMT has just randomly deployed %v apps. NMT ready to start experiment", experiments[0].ExperimentDetails.InitialAppsNumber.GetTotalAsString())

	}

	err = h.SimuClient.FetchAppsFromNMT()
	if err != nil {
		log.Errorf("Cannot fetch current app list from NMT. Error: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		log.Infof("Initial app list fetched from NMT")
		for i := 0; i < len(h.SimuClient.Apps); i++ {
			log.Infof("1. Apps[%v]: %v, cluster: %v", h.SimuClient.Apps[i].Id, h.SimuClient.Apps[i].UserPath, h.SimuClient.Apps[i].ClusterId)
			h.SimuClient.Apps[i].CurrentMove = 0
		}
	}

	trajectory, err := createTrajectory(movements, h)
	if err != nil {
		log.Errorf("Cannot create trajectory. Error: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		log.Infof("Trajectory has been created: %v", trajectory)
		for i := 0; i < len(h.SimuClient.Apps); i++ {
			log.Infof("2. Apps[%v]: %v, cluster: %v ", h.SimuClient.Apps[i].Id, h.SimuClient.Apps[i].UserPath, h.SimuClient.Apps[i].ClusterId)
		}
	}

	//loop for each experiment defined in method declareExperiments()
	for z, experiment := range experiments {

		err := checkExperimentType(experiment.ExperimentType)
		if err != nil {
			log.Errorf("Could not proceed with experiment. Reason: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Infof("Experiment [%v] type: %v, strategy: %v ", z+1, experiment.ExperimentType, experiment.ExperimentStrategy)

		//at the beggining let's recreate initial app placement at NMT and fetch
		err = h.SimuClient.RecreateInitialPlacementAtNMT()
		if err != nil {
			log.Errorf("Cannot recreate initial placement and fetch current app list from NMT. Error: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			log.Infof("Initial placement recreated at NMT")
		}

		err = h.SimuClient.FetchAppsFromNMT()
		if err != nil {
			log.Errorf("Cannot recreate initial placement and fetch current app list from NMT. Error: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			log.Infof("Initial app list fetched from NMT")
		}
		for i := 0; i < len(h.SimuClient.Apps); i++ {
			log.Infof("3. Apps[%v]: %v, cluster: %v ", h.SimuClient.Apps[i].Id, h.SimuClient.Apps[i].UserPath, h.SimuClient.Apps[i].ClusterId)
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
		for i := 0; i < len(trajectory); i++ {
			status := h.executePhDExperiment(experiment, z, i, trajectory[i][0], trajectory[i][1])
			if status != true {
				log.Error("Experiment cannot be coninued due to error in one of the iterations, skip this and let's go to next experiment")
				break
			}

		}

		err = h.ResultClient.CollectExperimentStats(experiment)
		if err != nil {
			log.Errorf("Error: %v. Status code: %v", err, http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		h.ResultClient.IncExpId()
		log.Infof("Finished Experiment [%v] type: %v", z+1, experiment.ExperimentType)
	}
	log.Infof("Finished all Experiments (%v), each %v iterations", len(experiments), experiments[0].ExperimentDetails.MovementsInExperiment)
	w.WriteHeader(http.StatusOK)
}
