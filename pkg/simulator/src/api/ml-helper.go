package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"simu/src/config"
	log "simu/src/logger"
	"simu/src/pkg/model"
	"strconv"
	"strings"
)

func executeMLExperiment(experiment ExperimentIntent, h *apiHandler, expIndex, subExpIndex int) bool {

	experimentN := "[EXPERIMENT " + strconv.Itoa(expIndex+1) + "." + strconv.Itoa(subExpIndex+1) + "] "
	//generate number of user to move
	id := h.generateUserToMove() //USER==APP

	// select new position for selected user and add new position to UserPath
	app := h.SimuClient.GetApps(id)
	h.generateTargetCellId(app)
	log.Infof(experimentN+"User(app) with ID: %v [current mec: %v] moved FROM cell: %v, towards cell: %v", app.Id, app.ClusterId, app.UserPath[len(app.UserPath)-2], app.UserLocation)

	spi, err := GenerateMLSmartPlacementIntent(*app)
	if err != nil {
		log.Errorf("Cannot generate SPI: %v", err.Error())
		return false
	}

	//send request to ML ctrl to select new position
	cluster, err := CallMLClient(spi)

	if err != nil {
		log.Warnf("ML Placement ctrl has returned status : %v", err.Error())
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

// SpaceAPP (for single app)  : 1) Required mvCPU 2) required Memory 3) Required Latency 4) Current MEC 5) Current RAN
func GenerateMLSmartPlacementIntent(app model.MECApp) (model.MLSmartPlacementIntent, error) {
	//log.Printf("GenerateSmartPlacementIntent: activity start\n")

	var spIntent model.MLSmartPlacementIntent
	clusterID, _ := convertMECNameToID(app.ClusterId)
	userLocation, _ := strconv.Atoi(app.UserLocation)

	appState := [5]int{
		determineReqRes(int(app.Requirements.RequestedCPU)),
		determineReqRes(int(app.Requirements.RequestedMEMORY)),
		determineStateofAppLatReq(int(app.Requirements.RequestedLatency)),
		clusterID,
		userLocation}

	url := buildNMTCurrentStateEndpoint()
	fmt.Printf("asking for MECs config url:, %v     |", url)

	mecState, err := GetMECsStateFromNMT(url)
	if err != nil {
		return spIntent, err
	}

	spIntent = model.MLSmartPlacementIntent{
		StateApp: model.SpaceAPP{
			AppCharacteristic: appState,
		},
		StateMECS: model.SpaceMECs{
			MecCharacteristics: mecState,
		},
	}
	log.Infof("GenerateSmartPlacementIntent: intent = %+v\n", spIntent)

	return spIntent, nil
}

func determineReqRes(reqRes int) int {
	resMap := map[int]int{
		500:  1,
		600:  2,
		700:  3,
		800:  4,
		900:  5,
		1000: 6,
	}
	if val, ok := resMap[reqRes]; ok {
		return val
	}
	return 0
}

func determineStateofAppLatReq(latValue int) int {
	latMap := map[int]int{
		10: 1,
		15: 2,
		30: 3,
	}
	if val, ok := latMap[latValue]; ok {
		return val
	}
	return 0
}

func CallMLClient(intent model.MLSmartPlacementIntent) (*model.Cluster, error) {
	//	log.Printf("CallPlacementController: function start\n")
	var respClusterID string

	plcCtrlUrl := buildMLPlcCtrlURL()
	data := intent

	responseBody, err := postHttpRespBody(plcCtrlUrl, data)
	if err != nil {
		return nil, fmt.Errorf("relocation for APP not done: %v", err)
	}

	err = json.Unmarshal(responseBody, &respClusterID)
	if err != nil {
		return nil, fmt.Errorf("relocation for APP not done: %v", err)
	}

	cluster := model.Cluster{
		Cluster:  "orange",
		Provider: "mec" + respClusterID,
	}

	return &cluster, nil
}

//todo: write proper url
func buildMLPlcCtrlURL() string {
	url := "ml"
	return url

}

func buildNMTCurrentStateEndpoint() string {
	///topology/ml/get-state
	url := config.GetConfiguration().NMTEndpoint
	url += "/v1/topology/ml/get-state"
	return url
}

func GetMECsStateFromNMT(endpoint string) ([][]int, error) {

	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	//Convert the body to type [][]int
	stateOfMECs := make([][]int, 22)
	for i := 0; i < len(stateOfMECs); i++ {
		stateOfMECs[i] = make([]int, 5)
	}

	json.Unmarshal(body, &stateOfMECs)
	return stateOfMECs, nil

}

func convertMECNameToID(s string) (int, error) {
	if !strings.HasPrefix(s, "mec") {
		return 0, fmt.Errorf("invalid input format: %s", s)
	}
	numStr := strings.TrimPrefix(s, "mec")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, fmt.Errorf("failed to convert %s to int: %v", numStr, err)
	}
	return num, nil
}
