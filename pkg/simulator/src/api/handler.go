package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"simu/src/config"
	log "simu/src/logger"
	"simu/src/pkg/model"
	"simu/src/pkg/results"
	"strconv"
	"strings"
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

func checkExperimentType(inputType string) (results.ExperimentType, error) {
	if strings.ToLower(inputType) == "optimal" {
		return results.ExpOptimal, nil

	} else if strings.ToLower(inputType) == "heuristic" {
		return results.ExpHeuristic, nil
	} else if strings.ToLower(inputType) == "ear-heuristic" || strings.ToLower(inputType) == "earheuristic" || strings.ToLower(inputType) == "ear" {
		return results.ExpEarHeuristic, nil
	}

	return results.ExpNotExists, fmt.Errorf("provided experiment type [%v] in not an option: %v", inputType, results.GetExpTypes())
}

func (h *apiHandler) conductExperiment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var intent ExperimentIntent

	err0 := json.NewDecoder(r.Body).Decode(&intent)
	if err0 != nil {
		log.Errorf("Cannot parse experiment intent. Error: %v", err0.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	experimentType, err := checkExperimentType(intent.ExperimentType)
	if err != nil {
		log.Errorf("Could not proceed with experiment. Reason: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	appsNumber, err := strconv.Atoi(intent.AppNumber)
	if err != nil {
		log.Errorf("Could not proceed with experiment. Reason: [apps-number] %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	experimentNumber, err := strconv.Atoi(intent.ExperimentsNumber)
	if err != nil {
		log.Errorf("Could not proceed with experiment. Reason: [experiments-number] %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Infof("Started new experiment: %v, with %v relocations", experimentType, experimentNumber)

	//at the beggining let's synchro latest placement at nmt
	//todo: run initial placement generator in NMT

	err = GenerateInitialAppPlacementAtNMT(intent.AppNumber)
	if err != nil {
		log.Errorf("Cannot make initial placement of app at NMT. Error: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		log.Infof("NMT has just randomly deployed %v apps. NMT ready to start experiment", intent.AppNumber)
	}
	//take initial topology and apps from NMT - done

	err2 := h.SimuClient.FetchAppsFromNMT()
	if err2 != nil {
		log.Errorf("Cannot fetch current app list from NMT. Error: %v", err2.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		log.Infof("Initial app list fetched from NMT")
	}

	//todo: reset results status at ERC before starting new expe
	err3 := resetResultsAtERC()
	if err3 != nil {
		log.Errorf("Cannot reset the results at NMT. Error: %v", err3.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		log.Infof("Results module ready -> cache cleared at NMT")
	}

	//check type of experiment
	//take statistics every M repetition

	for i := 0; i < experimentNumber; i++ {

		//log.Infof("Experiment numer: %v", i+1)

		experimentN := "[EXPERIMENT " + strconv.Itoa(i+1) + "] "
		//generate number of user to move
		id := h.generateUserToMove() //USER==APP
		//id := "10"

		// select new position for selected user and add new position to UserPath
		app := h.SimuClient.GetApps(id)
		h.generateTargetCellId(app)
		log.Infof(experimentN+"User(app) with ID: %v [current mec: %v] moved FROM cell: %v, towards cell: %v", app.Id, app.ClusterId, app.UserPath[len(app.UserPath)-2], app.UserLocation)

		//create smart placement intent

		spi, err := GenerateSmartPlacementIntent(*app, intent.Weights)
		if err != nil {
			log.Errorf("Cannot generate SPI: %v", err.Error())
		}

		// TODO: make sure that experimentType exists (use constants?)
		//send request to ERC to select new position
		cluster, err := CallPlacementController(spi, string(experimentType))
		if err != nil {
			log.Warnf(experimentN+"Error: %v", err.Error())
			log.Warnf(experimentN + "Stopped, NO RELOCATION, going to next iteration")
			continue
		}

		if cluster.Cluster == app.ClusterId {
			log.Warnf(experimentN+"Selected redundant cluster: %v -> missing relocation", cluster.Cluster)
			continue
		}

		log.Infof(experimentN+"Selected new cluster: %v", cluster.Cluster)

		//generate request to orchestrator
		err2 := sendRelocationRequest(*app, *cluster)
		if err2 != nil {
			log.Errorf("Cannot relocate app! Error: %v", err2.Error())
		} else {
			log.Infof(experimentN + "Application has been relocated in nmt")

			//update cluster in internal app list
			app.ClusterId = cluster.Cluster
		}

	}

	//todo: get results from nmt and return as a response

	var ercResults results.ErdResults
	var topoResults results.TopoResults

	ercUrl := config.GetConfiguration().ERCEndpoint + "/v2/erc/results"
	topoUrl := config.GetConfiguration().NMTEndpoint + "/v1/topology/mecHosts/metrics"

	ercBody, err := getHttpRespBody(ercUrl)
	topoBody, err := getHttpRespBody(topoUrl)

	err = json.Unmarshal(ercBody, &ercResults)
	if err != nil {
		log.Errorf("Error: %v. Status code: %v", err, http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	err = json.Unmarshal(topoBody, &topoResults.MecHostsResults)
	if err != nil {
		log.Errorf("Error: %v. Status code: %v", err, http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := results.ExpResult{
		Metadata: results.ExpResultsMeta{
			Type:      experimentType,
			Apps:      appsNumber,
			Movements: experimentNumber,
		},
		Data: results.ExpResultsData{
			Erd:  ercResults,
			Topo: topoResults,
		},
	}

	h.ResultClient.AppendResult(res)

	w.WriteHeader(http.StatusOK)
}
