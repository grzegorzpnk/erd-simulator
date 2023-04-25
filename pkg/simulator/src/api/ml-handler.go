package api

//
//func (h *apiHandler) conductMLExperiment(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//
//	var intent model.ExperimentIntent
//
//	err0 := json.NewDecoder(r.Body).Decode(&intent)
//	if err0 != nil {
//		log.Errorf("Cannot parse experiment intent. Error: %v", err0.Error())
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	experimentType, err := checkMLExperimentType(intent.ExperimentType)
//	if err != nil {
//		log.Errorf("Could not proceed with experiment. Reason: %v", err)
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	experimentIterations, err := strconv.Atoi(intent.ExperimentDetails.MovementsInExperiment)
//	if err != nil {
//		log.Errorf("Could not proceed with experiment. Reason: [experiments-iterations] %v", err)
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	log.Infof("Started new experiment: %v, with %v relocations", intent.ExperimentType, experimentIterations)
//
//	//at the beggining let's synchro latest placement at nmt
//	err = GenerateInitialAppPlacementAtNMT(intent.ExperimentDetails.InitialAppsNumber)
//	if err != nil {
//		log.Errorf("Cannot make initial placement of app at NMT. Error: %v", err.Error())
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	} else {
//		log.Infof("NMT has just randomly deployed %v apps. NMT ready to start experiment", intent.ExperimentDetails.InitialAppsNumber.GetTotalAsString())
//	}
//
//	err = h.SimuClient.FetchAppsFromNMT()
//	if err != nil {
//		log.Errorf("Cannot fetch current app list from NMT. Error: %v", err.Error())
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	} else {
//		log.Infof("Initial app list fetched from NMT")
//	}
//
//	////todo: reset results at ML placemen ctrl///to refactor
//	//err = resetResultsAtERC()
//	//if err != nil {
//	//	log.Errorf("Cannot reset the results at NMT. Error: %v", err.Error())
//	//	w.WriteHeader(http.StatusInternalServerError)
//	//	return
//	//} else {
//	//	log.Infof("Results module ready -> cache cleared at NMT")
//	//}
//	//
//
//	//This is already refactored
//	//loop for each experimentIterations defined in intent
//	for i := 0; i < experimentIterations; i++ {
//		status := executeMLExperiment(h, 1, i, experimentType)
//		if status != true {
//			log.Error("Experiment cannot be coninued due to error in one of the iterations")
//			w.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//		////todo: to refactor
//		//err = h.ResultClient.CollectExperimentStats(experimentType, "null", appsNumber, experimentIterations)
//		//if err != nil {
//		//	log.Errorf("Error: %v. Status code: %v", err, http.StatusInternalServerError)
//		//	w.WriteHeader(http.StatusInternalServerError)
//		//	return
//		//}
//	}
//
//	h.ResultClient.IncExpId()
//
//	w.WriteHeader(http.StatusOK)
//}
//
///// function for debugging purposes only
//func (h *apiHandler) stateTest(w http.ResponseWriter, r *http.Request) {
//
//	var intent model.ExperimentIntent
//
//	err0 := json.NewDecoder(r.Body).Decode(&intent)
//	if err0 != nil {
//		log.Errorf("Cannot parse experiment intent. Error: %v", err0.Error())
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	//appsNumber := intent.ExperimentDetails.InitialAppsNumber
//
//	experimentIterations, err := strconv.Atoi(intent.ExperimentDetails.MovementsInExperiment)
//	if err != nil {
//		log.Errorf("Could not proceed with experiment. Reason: [experiments-iterations] %v", err)
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	experimentType, err := checkMLExperimentType(intent.ExperimentType)
//	if err != nil {
//		log.Errorf("Could not proceed with experiment. Reason: %v", err)
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	log.Infof("Started new experiment: %v, with %v relocations", intent.ExperimentType, experimentIterations)
//
//	//at the beggining let's synchro latest placement at nmt
//	err = GenerateInitialAppPlacementAtNMT(intent.ExperimentDetails.InitialAppsNumber)
//	if err != nil {
//		log.Errorf("Cannot make initial placement of app at NMT. Error: %v", err.Error())
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	} else {
//		log.Infof("NMT has just randomly deployed %v apps. NMT ready to start experiment", intent.ExperimentDetails.InitialAppsNumber.GetTotalAsString())
//	}
//
//	err = h.SimuClient.FetchAppsFromNMT()
//	if err != nil {
//		log.Errorf("Cannot fetch current app list from NMT. Error: %v", err.Error())
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	} else {
//		log.Infof("Initial app list fetched from NMT")
//	}
//
//	//generate number of user to move
//	id := h.generateUserToMove() //USER==APP
//
//	// select new position for selected user and add new position to UserPath
//	app := h.SimuClient.GetApps(id)
//	h.generateTargetCellId(app)
//
//	spi, err := GenerateMLSmartPlacementIntent(*app, experimentType)
//	if err != nil {
//		log.Errorf("Cannot generate SPI: %v", err.Error())
//	}
//	json.NewEncoder(w).Encode(spi)
//	w.WriteHeader(http.StatusOK)
//}
