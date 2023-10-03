package api

import (
	"bytes"
	"encoding/json"
	"errors"
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

func GenerateSmartPlacementMLIntent(app model.MECApp) (model.SmartPlacementIntent, error) {
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
			},
		},
	}
	//log.Printf("GenerateSmartPlacementIntent: intent = %+v\n", spIntent)

	return spIntent, nil
}

//	type AppCounter struct {
//		Cg  int `json:"cg"`
//		V2x int `json:"v2x"`
//		Uav int `json:"uav"`
//	}

func CallPlacementController(intent model.SmartPlacementIntent, experimentType model.ExperimentType) (*model.Cluster, error) {
	//	log.Printf("CallPlacementController: function start\n")
	var resp model.Cluster

	plcCtrlUrl := buildPlcCtrlURL(string(experimentType))
	data := intent

	responseBody, err := postHttpRespBody(plcCtrlUrl, data)
	if err != nil {
		return nil, fmt.Errorf("relocation for APP[%v] not done: %v", intent.Metadata.Name, err)
	}

	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		return nil, fmt.Errorf("relocation for APP[%v] not done: %v", intent.Metadata.Name, err)
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
		return nil, fmt.Errorf("response code [%s] for URL [%s]", resp.Status, url)
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

func GenerateInitialAppPlacementAtNMT(appQuantity model.AppCounter) error {

	url := config.GetConfiguration().NMTEndpoint
	url += "/v1/topology/prerequesties/"
	url += "generate-apps"

	_, err := postHttpRespBody(url, appQuantity)
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

func checkExperimentType(et model.ExperimentType) error {

	for _, expType := range model.ExperimentTypes {
		if et == expType {
			return nil
		}
	}

	return errors.New(fmt.Sprintf("Expriment Type [%v] does not exist. Valid options: %v", et, model.ExperimentTypes))
}

func getHttpRespBody(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		getErr := fmt.Errorf("HTTP GET failed for URL %s.\nError: %s\n",
			url, err)
		fmt.Fprintf(os.Stderr, getErr.Error())
		return nil, getErr
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
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

func specifyStrategy(weights model.Weights) string {

	var strategy string

	if weights.LatencyWeight == 1 {
		strategy = "latency"
	} else if weights.ResourcesWeight == 1 {
		strategy = "lb"
	} else {
		strategy = "hybrid"
	}
	return strategy
}

func (h *apiHandler) executeExperiment(exp model.ExperimentIntent, expIndex, subExpIndex int) bool {

	experimentN := "[EXPERIMENT " + strconv.Itoa(expIndex) + "." + strconv.Itoa(subExpIndex+1) + "] "
	//generate number of user to move
	id := h.generateUserToMove() //USER==APP

	// select new position for selected user and add new position to UserPath
	app := h.SimuClient.GetApps(id)
	h.generateTargetCellId(app)
	log.Infof(experimentN+"User(app) with ID: %v [current mec: %v] moved FROM cell: %v, towards cell: %v", app.Id, app.ClusterId, app.UserPath[len(app.UserPath)-1], app.UserLocation)

	var spi model.SmartPlacementIntent
	var err error

	if exp.ExperimentType == model.ExpOptimal || exp.ExperimentType == model.ExpEarHeuristic || exp.ExperimentType == model.ExpHeuristic {
		spi, err = GenerateSmartPlacementIntent(*app, exp.Weights)
	} else {
		spi, err = GenerateSmartPlacementMLIntent(*app)
	}
	if err != nil {
		log.Errorf("Cannot generate SPI: %v", err.Error())
		return false
	}

	//send request to ERC to select new position
	cluster, err := CallPlacementController(spi, exp.ExperimentType)

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

func (h *apiHandler) executeGlobcomExperiment(exp model.ExperimentIntent, expIndex, subExpIndex int, userID, userPosition int) bool {

	experimentN := "[EXPERIMENT " + strconv.Itoa(expIndex) + "." + strconv.Itoa(subExpIndex+1) + "] "

	app := h.SimuClient.GetApps(strconv.Itoa(userID))
	app.UserLocation = strconv.Itoa(userPosition)

	log.Infof(experimentN+"User(app) with ID: %v [current mec: %v] moved FROM cell: %v, towards cell: %v", app.Id, app.ClusterId, app.UserPath[len(app.UserPath)-1], app.UserLocation)

	var spi model.SmartPlacementIntent
	var err error

	if exp.ExperimentType == model.ExpOptimal || exp.ExperimentType == model.ExpEarHeuristic || exp.ExperimentType == model.ExpHeuristic {
		spi, err = GenerateSmartPlacementIntent(*app, exp.Weights)
	} else {
		spi, err = GenerateSmartPlacementMLIntent(*app)
	}
	if err != nil {
		log.Errorf("Cannot generate SPI: %v", err.Error())
		return false
	}

	//send request to ERC to select new position
	cluster, err := CallPlacementController(spi, exp.ExperimentType)

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

//copy of executeGlobcomExperiment()
func (h *apiHandler) executePhDExperiment(exp model.ExperimentIntent, expIndex, subExpIndex int, userID, userPosition int) bool {

	experimentN := "[EXPERIMENT " + strconv.Itoa(expIndex) + "." + strconv.Itoa(subExpIndex+1) + "] "

	app := h.SimuClient.GetApps(strconv.Itoa(userID))
	//log.Infof(experimentN+"User(app) with ID: %v (%v ms) [current mec: %v] has moved FROM cell: %v, towards cell: %v", app.Id, app.Requirements.RequestedLatency, app.ClusterId, app.UserLocation, strconv.Itoa(userPosition))
	log.Infof(experimentN)
	app.UserLocation = strconv.Itoa(userPosition)

	var spi model.SmartPlacementIntent
	var err error

	if exp.ExperimentType == model.ExpOptimal || exp.ExperimentType == model.ExpEarHeuristic || exp.ExperimentType == model.ExpHeuristic {
		spi, err = GenerateSmartPlacementIntent(*app, exp.Weights)
	} else {
		spi, err = GenerateSmartPlacementMLIntent(*app)
	}
	if err != nil {
		log.Errorf("Cannot generate SPI: %v", err.Error())
		return false
	}

	//send request to ERC to select new position
	cluster, err := CallPlacementController(spi, exp.ExperimentType)

	if err != nil {
		//log.Warnf("Call Placement ctrl has returned status : %v", err.Error())
		//log.Warnf(experimentN + "stopped, NO RELOCATION, going to next iteration")
		return true
	}

	if cluster.Cluster == app.ClusterId {
		//log.Infof(experimentN+"Selected redundant cluster: %v -> missing relocation", cluster.Cluster)
		return true
	}

	//log.Infof(experimentN+"Selected new cluster: %v", cluster.Cluster)

	//generate request to orchestrator
	err2 := sendRelocationRequest(*app, *cluster)
	if err2 != nil {
		//	log.Errorf("Cannot relocate app! Error: %v", err2.Error())
	} else {
		//	log.Infof(experimentN + "Application has been relocated in nmt")

		//update cluster in internal app list
		app.ClusterId = cluster.Cluster
	}

	return true
}

func declareExperiments(details model.ExperimentDetails) []model.ExperimentIntent {

	experiments := []model.ExperimentIntent{}
	experiment1 := model.ExperimentIntent{
		ExperimentType:    "optimal",
		ExperimentDetails: details,
		Weights: model.Weights{
			LatencyWeight:        0.5,
			ResourcesWeight:      0.5,
			CpuUtilizationWeight: 0.5,
			MemUtilizationWeight: 0.5,
		},
	}
	experiment2 := model.ExperimentIntent{
		ExperimentType:    "optimal",
		ExperimentDetails: details,
		Weights: model.Weights{
			LatencyWeight:        1,
			ResourcesWeight:      0,
			CpuUtilizationWeight: 0,
			MemUtilizationWeight: 0,
		},
	}
	experiment3 := model.ExperimentIntent{
		ExperimentType:    "optimal",
		ExperimentDetails: details,
		Weights: model.Weights{
			LatencyWeight:        0,
			ResourcesWeight:      1,
			CpuUtilizationWeight: 0.5,
			MemUtilizationWeight: 0.5,
		},
	}
	experiment4 := model.ExperimentIntent{
		ExperimentType:    "heuristic",
		ExperimentDetails: details,
		Weights: model.Weights{
			LatencyWeight:        0.5,
			ResourcesWeight:      0.5,
			CpuUtilizationWeight: 0.5,
			MemUtilizationWeight: 0.5,
		},
	}
	experiment5 := model.ExperimentIntent{
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

func declareGlobcomExperiments(details model.ExperimentDetails) []model.ExperimentIntent {

	experiments := []model.ExperimentIntent{}

	experiment1 := model.ExperimentIntent{
		ExperimentType:     model.ExpOptimal,
		ExperimentStrategy: model.StrHybrid,
		ExperimentDetails:  details,
		Weights: model.Weights{
			LatencyWeight:        0.5,
			ResourcesWeight:      0.5,
			CpuUtilizationWeight: 0.5,
			MemUtilizationWeight: 0.5,
		},
	}

	experiment2 := model.ExperimentIntent{
		ExperimentType:     model.ExpHeuristic,
		ExperimentStrategy: model.StrHybrid,
		ExperimentDetails:  details,
		Weights: model.Weights{
			LatencyWeight:        0.5,
			ResourcesWeight:      0.5,
			CpuUtilizationWeight: 0.5,
			MemUtilizationWeight: 0.5,
		},
	}

	experiment3 := model.ExperimentIntent{
		ExperimentType:     model.ExpEarHeuristic,
		ExperimentStrategy: model.StrHybrid,
		ExperimentDetails:  details,
		Weights: model.Weights{
			LatencyWeight:        0.5,
			ResourcesWeight:      0.5,
			CpuUtilizationWeight: 0.5,
			MemUtilizationWeight: 0.5,
		},
	}

	experiment4 := model.ExperimentIntent{
		ExperimentType:     model.ExpMLMasked,
		ExperimentStrategy: model.StrML,
		ExperimentDetails:  details,
	}

	experiment5 := model.ExperimentIntent{
		ExperimentType:     model.ExpMLNonMasked,
		ExperimentStrategy: model.StrML,
		ExperimentDetails:  details,
	}

	experiments = append(experiments, experiment1, experiment2, experiment3, experiment4, experiment5)

	return experiments
}

func declareICCExperiments(details model.ExperimentDetails) []model.ExperimentIntent {

	experiments := []model.ExperimentIntent{}

	experiment1 := model.ExperimentIntent{
		ExperimentType:     model.ExpOptimal,
		ExperimentStrategy: model.StrHybrid,
		ExperimentDetails:  details,
		Weights: model.Weights{
			LatencyWeight:        0.5,
			ResourcesWeight:      0.5,
			CpuUtilizationWeight: 0.5,
			MemUtilizationWeight: 0.5,
		},
	}

	experiment2 := model.ExperimentIntent{
		ExperimentType:     model.ExpOptimal,
		ExperimentStrategy: model.StrLB,
		ExperimentDetails:  details,
		Weights: model.Weights{
			LatencyWeight:        0,
			ResourcesWeight:      1,
			CpuUtilizationWeight: 0.5,
			MemUtilizationWeight: 0.5,
		},
	}

	experiment3 := model.ExperimentIntent{
		ExperimentType:     model.ExpOptimal,
		ExperimentStrategy: model.StrLatency,
		ExperimentDetails:  details,
		Weights: model.Weights{
			LatencyWeight:        1,
			ResourcesWeight:      0,
			CpuUtilizationWeight: 0,
			MemUtilizationWeight: 0,
		},
	}

	experiment4 := model.ExperimentIntent{
		ExperimentType:     model.ExpHeuristic,
		ExperimentStrategy: model.StrHybrid,
		ExperimentDetails:  details,
		Weights: model.Weights{
			LatencyWeight:        0.5,
			ResourcesWeight:      0.5,
			CpuUtilizationWeight: 0.5,
			MemUtilizationWeight: 0.5,
		},
	}

	experiment5 := model.ExperimentIntent{
		ExperimentType:     model.ExpEarHeuristic,
		ExperimentStrategy: model.StrHybrid,
		ExperimentDetails:  details,
		Weights: model.Weights{
			LatencyWeight:        0.5,
			ResourcesWeight:      0.5,
			CpuUtilizationWeight: 0.5,
			MemUtilizationWeight: 0.5,
		},
	}

	experiment6 := model.ExperimentIntent{
		ExperimentType:     model.ExpEarHeuristic,
		ExperimentStrategy: model.StrLB,
		ExperimentDetails:  details,
		Weights: model.Weights{
			LatencyWeight:        0,
			ResourcesWeight:      1,
			CpuUtilizationWeight: 0.5,
			MemUtilizationWeight: 0.5,
		},
	}

	experiments = append(experiments, experiment1, experiment2, experiment3, experiment4, experiment5, experiment6)
	//experiments = append(experiments, experiment3)

	return experiments
}

func declareICCTunningExperiments(details model.ExperimentDetails) []model.ExperimentIntent {

	experiments := []model.ExperimentIntent{}

	experiment1 := model.ExperimentIntent{
		ExperimentType:     model.ExpEarHeuristic,
		ExperimentStrategy: model.StrHybrid,
		ExperimentDetails:  details,
		Weights: model.Weights{
			LatencyWeight:        0.5,
			ResourcesWeight:      0.5,
			CpuUtilizationWeight: 0.5,
			MemUtilizationWeight: 0.5,
		},
	}

	experiment2 := model.ExperimentIntent{
		ExperimentType:     model.ExpEarHeuristic,
		ExperimentStrategy: model.StrLB,
		ExperimentDetails:  details,
		Weights: model.Weights{
			LatencyWeight:        0,
			ResourcesWeight:      1,
			CpuUtilizationWeight: 0.5,
			MemUtilizationWeight: 0.5,
		},
	}

	experiment3 := model.ExperimentIntent{
		ExperimentType:     model.ExpEarHeuristic,
		ExperimentStrategy: model.StrLatency,
		ExperimentDetails:  details,
		Weights: model.Weights{
			LatencyWeight:        1,
			ResourcesWeight:      0,
			CpuUtilizationWeight: 0,
			MemUtilizationWeight: 0,
		},
	}

	experiment4 := model.ExperimentIntent{
		ExperimentType:     model.ExpEarHeuristic,
		ExperimentStrategy: model.Str7L3R,
		ExperimentDetails:  details,
		Weights: model.Weights{
			LatencyWeight:        0.7,
			ResourcesWeight:      0.3,
			CpuUtilizationWeight: 0.5,
			MemUtilizationWeight: 0.5,
		},
	}

	experiment5 := model.ExperimentIntent{
		ExperimentType:     model.ExpEarHeuristic,
		ExperimentStrategy: model.Str3L7R,
		ExperimentDetails:  details,
		Weights: model.Weights{
			LatencyWeight:        0.3,
			ResourcesWeight:      0.7,
			CpuUtilizationWeight: 0.5,
			MemUtilizationWeight: 0.5,
		},
	}

	experiments = append(experiments, experiment1, experiment2, experiment3, experiment4, experiment5)
	//experiments = append(experiments, experiment3)

	return experiments
}

func createTrajectory(movements int, h *apiHandler) ([][]int, error) {

	trajectory := make([][]int, movements)

	for i := 0; i < movements; i++ {
		UserID := h.generateUserToMove() //USER==APP

		// select new position for selected user and add new position to UserPath
		app := h.SimuClient.GetApps(UserID)
		h.generateTargetCellId(app)

		trajectory[i] = make([]int, 2)
		trajectory[i][0], _ = strconv.Atoi(UserID)
		trajectory[i][1], _ = strconv.Atoi(app.UserLocation)

	}

	return trajectory, nil
}

func DeclareApplications(ac model.AppCounter) []model.MECApp {

	var mecList []model.MECApp

	v2x, drones, video := ac.V2x, ac.Uav, ac.Cg

	for i := 0; i < v2x; i++ {
		var app model.MECApp
		app.Id = strconv.Itoa(i + 1)
		mecList = append(mecList, app)
	}
	for i := v2x; i < v2x+drones; i++ {
		var app model.MECApp
		app.Id = strconv.Itoa(i + 1)
		mecList = append(mecList, app)
	}
	for i := v2x + drones; i < v2x+drones+video; i++ {
		var app model.MECApp
		app.Id = strconv.Itoa(i + 1)
		mecList = append(mecList, app)
	}

	return mecList
}
