package api

import (
	"encoding/json"
	"net/http"
	log "simu/src/logger"
	"strconv"
)

func (h *apiHandler) conductMLExperiment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var intent ExperimentIntent

	err0 := json.NewDecoder(r.Body).Decode(&intent)
	if err0 != nil {
		log.Errorf("Cannot parse experiment intent. Error: %v", err0.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//appsNumber := intent.ExperimentDetails.AppsNumber

	experimentIterations, err := strconv.Atoi(intent.ExperimentDetails.ExperimentIterations)
	if err != nil {
		log.Errorf("Could not proceed with experiment. Reason: [experiments-iterations] %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Infof("Started new experiment: %v, with %v relocations", intent.ExperimentType, experimentIterations)

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

	////todo: reset results at ML placemen ctrl///to refactor
	//err = resetResultsAtERC()
	//if err != nil {
	//	log.Errorf("Cannot reset the results at NMT. Error: %v", err.Error())
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//} else {
	//	log.Infof("Results module ready -> cache cleared at NMT")
	//}
	//

	//This is already refactored
	//loop for each experimentIterations defined in intent
	for i := 0; i < experimentIterations; i++ {
		status := executeMLExperiment(intent, h, 1, i)
		if status != true {
			log.Error("Experiment cannot be coninued due to error in one of the iterations")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		////todo: to refactor
		//err = h.ResultClient.CollectExperimentStats(experimentType, "null", appsNumber, experimentIterations)
		//if err != nil {
		//	log.Errorf("Error: %v. Status code: %v", err, http.StatusInternalServerError)
		//	w.WriteHeader(http.StatusInternalServerError)
		//	return
		//}
	}

	h.ResultClient.IncExpId()

	w.WriteHeader(http.StatusOK)
}

//
//func (h *apiHandler) stateTest (w http.ResponseWriter, r *http.Request) {
//
//	GenerateMLSmartPlacementIntent(app
//	model.MECApp)
//
//}
