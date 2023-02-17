package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"simu/src/config"
	log "simu/src/logger"
	"simu/src/pkg/model"
	"strconv"
)

func GenerateSmartPlacementIntent(app model.MECApp, weights model.Weights) (model.SmartPlacementIntent, error) {
	//log.Printf("GenerateSmartPlacementIntent: activity start\n")

	targetAppName := app.Id

	spIntent := model.SmartPlacementIntent{
		Metadata: model.Metadata{
			Name:        targetAppName + "-er-intent",
			Description: fmt.Sprintf("Edge Relocation Intent for app: %s", targetAppName),
		},
		CurrentPlacement: model.Cluster{
			Provider: "orange",
			Cluster:  app.ClusterId,
		},
		Spec: model.SmartPlacementIntentSpec{
			AppName: targetAppName,
			SmartPlacementIntentData: model.SmartPlacementIntentStruct{
				TargetCell: app.UserLocation,
				AppCpuReq:  app.Requirements.RequestedCPU,
				AppMemReq:  app.Requirements.RequestedMEMORY,
				ConstraintsList: model.Constraints{
					LatencyMax: app.Requirements.RequestedLatency,
				},
				ParametersWeights: model.Weights{
					LatencyWeight:        weights.LatencyWeight,
					ResourcesWeight:      weights.ResourcesWeight,
					CpuUtilizationWeight: weights.CpuUtilizationWeight,
					MemUtilizationWeight: weights.MemUtilizationWeight,
				},
			},
		},
	}
	//log.Printf("GenerateSmartPlacementIntent: intent = %+v\n", spIntent)

	return spIntent, nil
}

type ExperimentIntent struct {
	ExperimentType    string            `json:"experiment-type"`
	ExperimentDetails ExperimentDetails `json:"experiment-details"`
	Weights           model.Weights     `json:"Weights"`
}

type ExperimentDetails struct {
	ExperimentIterations string `json:"experiments-number"`
	AppNumber            string `json:"app-number"`
}

func CallPlacementController(intent model.SmartPlacementIntent, experimentType string) (*model.Cluster, error) {
	//	log.Printf("CallPlacementController: function start\n")
	var resp model.Cluster

	plcCtrlUrl := buildPlcCtrlURL(experimentType)
	data := intent

	responseBody, err := postHttpRespBody(plcCtrlUrl, data)
	if err != nil {
		log.Warnf("[ERROR] Placement Controller returned status: %v. Relocation for APP[%v] not done.", err, intent.Metadata.Name)
		return nil, err
	}

	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		log.Warnf("error occured while unmarshaling: %v. Resp body: %v", err, string(responseBody))
		return nil, err
	}

	var cluster model.Cluster
	cluster = resp

	return &cluster, nil
}

func buildPlcCtrlURL(experimentType string) string {

	url := config.GetConfiguration().ERCEndpoint
	url += "/v2/erc/smart-placement-intents/optimal-mec/"
	url += experimentType

	return url

}

func postHttpRespBody(url string, data interface{}) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error: marshaling failed")
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Errorf("Could not make post request. reason: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		getErr := fmt.Errorf("HTTP GET returned status code %s for URL %s.\n",
			resp.Status, url)
		fmt.Fprintf(os.Stderr, getErr.Error())
		return nil, getErr
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return b, nil
}

func sendRelocationRequest(app model.MECApp, newCluster model.Cluster) error {

	orchestratorEndpoint := buildOrchestratorURL(app, newCluster)

	_, err := postHttpRespBody(orchestratorEndpoint, app)
	if err != nil {
		log.Errorf("[ERROR] Orchestrator returned error: %v. Relocation for APP[%v] failed.", err, app.Id)
		return err
	} else {
		//	log.Infof("Orchestrator returned confirmation: %v. Relocation for APP[%v] done.", responseBody, app.Id)
		return nil
	}
}

func resetResultsAtERC() error {

	url := config.GetConfiguration().ERCEndpoint
	url += "/v2/erc/results/reset"

	_, err := postHttpRespBody(url, nil)
	if err != nil {
		log.Errorf("[ERROR] ERC returned error: %v ", err)
		return err
	} else {
		return nil
	}

}

func GenerateInitialAppPlacementAtNMT(appQuantity string) error {

	url := config.GetConfiguration().NMTEndpoint
	url += "/v1/topology/prerequesties/"
	url += appQuantity

	_, err := postHttpRespBody(url, nil)
	if err != nil {
		log.Errorf("[ERROR] NMT returned error: %v ", err)
		return err
	}

	return nil
}

func buildOrchestratorURL(app model.MECApp, cluster model.Cluster) string {
	//orchestrator/relocate/old-cluster/{old-cluster}/new-cluster/{new-cluster}/application"

	url := config.GetConfiguration().NMTEndpoint
	url += "/v1/orchestrator/relocate/old-cluster/"
	url += app.ClusterId + "/new-cluster/"
	url += cluster.Cluster + "/application"

	return url

}

func specifyStrategy(weights model.Weights) string {

	var strategy string

	if weights.LatencyWeight == 1 {
		strategy = "latency min"
	} else if weights.ResourcesWeight == 1 {
		strategy = "resources LB"
	} else {
		strategy = "hybrid"
	}
	return strategy
}

func executeExperiment(experiment ExperimentIntent, h *apiHandler, expIndex, subExpIndex int) bool {

	//log.Infof("Experiment numer: %v", i+1)

	experimentN := "[EXPERIMENT " + strconv.Itoa(expIndex+1) + "." + strconv.Itoa(subExpIndex+1) + "] "
	//generate number of user to move
	id := h.generateUserToMove() //USER==APP
	//id := "10"

	// select new position for selected user and add new position to UserPath
	app := h.SimuClient.GetApps(id)
	h.generateTargetCellId(app)
	log.Infof(experimentN+"User(app) with ID: %v [current mec: %v] moved FROM cell: %v, towards cell: %v", app.Id, app.ClusterId, app.UserPath[len(app.UserPath)-2], app.UserLocation)

	//create smart placement intent

	spi, err := GenerateSmartPlacementIntent(*app, experiment.Weights)
	if err != nil {
		log.Errorf("Cannot generate SPI: %v", err.Error())
		return false
	}

	//send request to ERC to select new position
	cluster, err := CallPlacementController(spi, experiment.ExperimentType)

	if err != nil {
		log.Warnf("Call Placement ctrl has returned status : %v", err.Error())
		log.Warnf(experimentN + "stopped, NO RELOCATION, going to next iteration")
		return true
	}

	if cluster.Cluster == app.ClusterId {
		log.Infof(experimentN+"Selected redundant cluster: %v -> missing relocation", cluster.Cluster)
		return true
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

	return true
}

func declareExperiments(details ExperimentDetails) []ExperimentIntent {

	experiments := []ExperimentIntent{}
	experiment1 := ExperimentIntent{
		ExperimentType:    "optimal",
		ExperimentDetails: details,
		Weights: model.Weights{
			LatencyWeight:        0.5,
			ResourcesWeight:      0.5,
			CpuUtilizationWeight: 0.5,
			MemUtilizationWeight: 0.5,
		},
	}
	experiment2 := ExperimentIntent{
		ExperimentType:    "optimal",
		ExperimentDetails: details,
		Weights: model.Weights{
			LatencyWeight:        1,
			ResourcesWeight:      0,
			CpuUtilizationWeight: 0,
			MemUtilizationWeight: 0,
		},
	}
	experiment3 := ExperimentIntent{
		ExperimentType:    "optimal",
		ExperimentDetails: details,
		Weights: model.Weights{
			LatencyWeight:        0,
			ResourcesWeight:      1,
			CpuUtilizationWeight: 0.5,
			MemUtilizationWeight: 0.5,
		},
	}
	experiment4 := ExperimentIntent{
		ExperimentType:    "heuristic",
		ExperimentDetails: details,
		Weights: model.Weights{
			LatencyWeight:        0.5,
			ResourcesWeight:      0.5,
			CpuUtilizationWeight: 0.5,
			MemUtilizationWeight: 0.5,
		},
	}
	experiment5 := ExperimentIntent{
		ExperimentType:    "ear-heuristic",
		ExperimentDetails: details,
		Weights: model.Weights{
			LatencyWeight:        0.5,
			ResourcesWeight:      0.5,
			CpuUtilizationWeight: 0.5,
			MemUtilizationWeight: 0.5,
		},
	}
	experiments = append(experiments, experiment1, experiment2, experiment3, experiment4, experiment5)

	return experiments
}
