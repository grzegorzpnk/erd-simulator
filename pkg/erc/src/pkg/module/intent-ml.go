package module

import (
	"10.254.188.33/matyspi5/erd/pkg/erc/src/config"
	log "10.254.188.33/matyspi5/erd/pkg/erc/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/model"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/topology"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func (i *SmartPlacementIntentClient) ServeSmartPlacementIntentML(checkIfMasked bool, intent model.SmartPlacementIntent) (model.MecHost, error) {

	var err error
	var bestMec model.MecHost

	tc := topology.NewTopologyClient()
	tc.CurrentCell = intent.Spec.SmartPlacementIntentData.TargetCell

	log.Infof("Received request to find new cluster for App [%v] located at Cluster [%v], that moves towards Cell [%v]", intent.Spec.AppName, intent.CurrentPlacement.Cluster, intent.Spec.SmartPlacementIntentData.TargetCell)
	log.Infof("Smart Placement Intent: %+v", intent)
	if checkIfMasked {
		log.Infof("Searching Type: ML MASKED")
	} else {
		log.Infof("Searching Type: ML NON MASKED")
	}

	MLspi, err := GenerateMLSmartPlacementIntent(intent, checkIfMasked)
	if err != nil {
		log.Errorf("Cannot generate SPI: %v", err.Error())
		return model.MecHost{}, err
	}

	//send request to ML ctrl to select new position
	bestMec, err = CallMLClient(MLspi)
	if err != nil {
		log.Warnf(" Could not find optimal cluster for given APP[%v]. Error[%v]", intent.Spec.AppName, err)
		return model.MecHost{}, err
	}

	// check if relocation redundant
	if bestMec.Identity.Provider == intent.CurrentPlacement.Provider &&
		bestMec.Identity.Cluster == intent.CurrentPlacement.Cluster {
		log.Infof("[POSITIVE] Relocation redundant. Selected the same mec host.")
		return bestMec, nil
	}
	// Collect info from topology
	bestMec.Resources.Latency, err = tc.GetShortestPath(tc.CurrentCell, bestMec)
	if err != nil {
		log.Warnf("Could not collect MEC [%v+%v] LATENCY info", bestMec.Identity.Provider, bestMec.Identity.Cluster)
		return model.MecHost{}, err
	}

	bestMec.Resources.Cpu, err = tc.GetMecResource(model.MecCpu, bestMec)
	if err != nil {
		log.Warnf("Could not collect MEC [%v+%v] CPU info", bestMec.Identity.Provider, bestMec.Identity.Cluster)
		return model.MecHost{}, err
	}

	bestMec.Resources.Memory, err = tc.GetMecResource(model.MecMem, bestMec)
	if err != nil {
		log.Warnf("Could not collect MEC [%v+%v] MEMORY info", bestMec.Identity.Provider, bestMec.Identity.Cluster)
		return model.MecHost{}, err
	}

	// Check if ok

	cpuUtilAfterRel := 100 * (bestMec.GetCpuUsed() + intent.Spec.SmartPlacementIntentData.AppCpuReq) / bestMec.GetCpuCapacity()
	memUtilAfterRel := 100 * (bestMec.GetMemUsed() + intent.Spec.SmartPlacementIntentData.AppMemReq) / bestMec.GetMemCapacity()

	reqLatency := intent.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax

	threshold := 80.0
	//log.Infof("[DEBUG], bestMec.GetCpuUsed: %v. AppCpuReq: %v, bestMec.getCpuCapacity: %v.", bestMec.GetCpuUsed(), intent.Spec.SmartPlacementIntentData.AppCpuReq, bestMec.GetCpuCapacity())

	if bestMec.GetLatency() > reqLatency {
		log.Errorf("Bad Lat cluster was selected by ML model: ")
		err = errors.New(fmt.Sprintf("mec latency offered [%v] is unfortunately higher than required [%v]", bestMec.GetLatency(), reqLatency))
		return model.MecHost{}, err
	} else if cpuUtilAfterRel > threshold || memUtilAfterRel > threshold {
		log.Errorf("Bad RES cluster was selected by ML model: ")
		err = errors.New(fmt.Sprintf("mec resource utilization after relocation [cpu: %v, memory: %v] would be higher than allowed threshold [%v]", cpuUtilAfterRel, memUtilAfterRel, threshold))
		return model.MecHost{}, err
	}

	log.Infof("Seems good cluster was identified, after relocation: ")
	log.Infof("[%v] CPU Util [%v], Mem Util [%v], threshold [%v]", bestMec.Identity.Cluster, cpuUtilAfterRel, memUtilAfterRel, threshold)
	log.Infof("ML CLient has found cluster[%v] for given APP[%v]\n", bestMec.Identity.Cluster, intent.Spec.AppName)
	return bestMec, nil

}

func CallMLClient(intent model.MLSmartPlacementIntent) (model.MecHost, error) {
	//	log.Printf("CallPlacementController: function start")
	var respClusterID json.Number

	plcCtrlUrl := buildMLPlcCtrlURL()
	data := intent

	responseBody, err := postHttpRespBody(plcCtrlUrl, data)
	if err != nil {
		return model.MecHost{}, fmt.Errorf("relocation for APP not done: %v", err)
	}

	err = json.Unmarshal(responseBody, &respClusterID)
	if err != nil {
		return model.MecHost{}, fmt.Errorf("relocation for APP not done: %v", err)
	}

	cluster := model.MecHost{
		Identity: model.MecIdentity{
			Provider: "orange",
			Cluster:  "mec" + respClusterID.String(),
		},
	}

	return cluster, nil
}

func buildMLPlcCtrlURL() string {
	url := config.GetConfiguration().MLClientEndpoint
	url += "/ermodel/get-prediction"

	return url
}

// SpaceAPP (for single app)  : 1) Required mvCPU 2) required Memory 3) Required Latency 4) Current MEC 5) Current RAN
func GenerateMLSmartPlacementIntent(intent model.SmartPlacementIntent, checkIfMasked bool) (model.MLSmartPlacementIntent, error) {
	//log.Printf("GenerateSmartPlacementIntent: activity start")

	var spIntent model.MLSmartPlacementIntent
	var state model.State

	clusterID, _ := convertMECNameToID(intent.CurrentPlacement.Cluster)
	userLocation, _ := strconv.Atoi(string(intent.Spec.SmartPlacementIntentData.TargetCell))

	var appState [][]int

	if checkIfMasked {
		appState = append(appState, []int{
			determineReqRes(int(intent.Spec.SmartPlacementIntentData.AppCpuReq)),
			determineReqRes(int(intent.Spec.SmartPlacementIntentData.AppMemReq)),
			determineStateofAppLatReq(int(intent.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax)),
			clusterID,
			userLocation})
	} else {
		//appState = append(appState, []int{
		//	determineReqResInEdgeContext(int(intent.Spec.SmartPlacementIntentData.AppCpuReq), intent.CurrentPlacement.Cluster),
		//	determineReqResInEdgeContext(int(intent.Spec.SmartPlacementIntentData.AppMemReq), intent.CurrentPlacement.Cluster),
		//	int(intent.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax),
		//	clusterID})
		appState = append(appState, []int{int(intent.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax)})
	}

	url := buildNMTCurrentStateEndpoint()
	mecState, err := GetMECsStateFromNMT(url, intent, checkIfMasked)
	if err != nil {
		return spIntent, err
	}

	state = model.State{
		SpaceAPP:  appState,
		SpaceMECs: mecState,
	}

	app := model.MECApp{Id: intent.Spec.AppName,
		ClusterId:    intent.CurrentPlacement.Cluster,
		UserLocation: string(intent.Spec.SmartPlacementIntentData.TargetCell),
		Requirements: model.RequestedResources{
			RequestedCPU:     intent.Spec.SmartPlacementIntentData.AppCpuReq,
			RequestedMEMORY:  intent.Spec.SmartPlacementIntentData.AppMemReq,
			RequestedLatency: intent.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax,
		},
	}

	if checkIfMasked {
		mask, err := GenerateMLMask(app)
		if err != nil {
			log.Errorf("Cannot generate Mask: %v", err.Error())
		}
		spIntent = model.MLSmartPlacementIntent{State: state, CurrentMask: mask}
	} else if !checkIfMasked {
		spIntent = model.MLSmartPlacementIntent{State: state}
	} else {
		log.Errorf("Invalid type of experiment: %v", err.Error())
		return spIntent, err
	}

	log.Infof("GenerateMLSmartPlacementIntent: intent = %+v", spIntent)

	return spIntent, nil
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

func GenerateMLMask(app model.MECApp) ([]int, error) {

	mask := make([]int, 22)

	url := buildNMTMaskEndpoint()
	//fmt.Printf("Asking for Mask url:, %v     |", url)

	mask, err := GetMaskFromNMT(url, app)
	if err != nil {
		return mask, err
	}

	log.Infof("Mask: %+v", mask)

	return mask, nil
}

func CheckMecCapacity(mecID string) (int, error) {

	url := buildNMTEndpointGetClusterCapacity(mecID)

	capacity, err := GetCapacityCLusterFromNMT(url)
	if err != nil {
		return capacity, err
	}

	log.Infof("Checked Capacity of MEC: %v ", capacity)

	return capacity, nil
}

func determineReqRes(reqRes int) int {
	//resMap := map[int]int{
	//	500:  1,
	//	600:  2,
	//	700:  3,
	//	800:  4,
	//	900:  5,
	//	1000: 6,
	//}
	//if val, ok := resMap[reqRes]; ok {
	//	return val
	//}
	//return 0

	return (reqRes - 500)
}

func determineReqResInEdgeContext(reqRes int, cluster string) int {

	clusterCapacity, err := CheckMecCapacity(cluster)
	if err != nil {
		log.Errorf("CANNOT CREATE CURRENT STATE FOR NONMAKSED RL :(")
		return 0
	}

	log.Infof("APP req res: %v, while cluster capacity is %v", reqRes, clusterCapacity)
	return int(reqRes * 100 / clusterCapacity)

}

func reverseDetermineReqResFloat64(reqRes int) float64 {
	resMap := map[int]float64{
		1: 500.0,
		2: 600.0,
		3: 700.0,
		4: 800.0,
		5: 900.0,
		6: 1000.0,
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

func reverseDetermineAppLatReqFloat64(latValue int) float64 {
	latMap := map[int]float64{
		1: 10.0,
		2: 15.0,
		3: 30.0,
	}
	if val, ok := latMap[latValue]; ok {
		return val
	}
	return 0
}

func buildNMTMaskEndpoint() string {
	url := config.GetConfiguration().NMTEndpoint
	url += "/topology/ml/get-mask"

	return url
}

func buildNMTEndpointGetClusterCapacity(clusterID string) string {
	url := config.GetConfiguration().NMTEndpoint
	url += "/topology/mecHosts/provider/orange/cluster/" +
		clusterID + "/capacity"

	return url
}

func buildNMTCurrentStateEndpoint() string {
	///topology/ml/get-state
	url := config.GetConfiguration().NMTEndpoint
	url += "/topology/ml/get-state"
	return url
}

func GetMECsStateFromNMT(endpoint string, intent model.SmartPlacementIntent, isMasked bool) ([][]int, error) {

	if isMasked {
		intent.Spec.SmartPlacementIntentData.Masked = "masked"
	}

	resp, err := postHttpRespBody(endpoint, intent)
	if err != nil {
		return nil, err
	}
	//defer resp.Body.Close()

	//Convert the body to type [][]int
	stateOfMECs := make([][]int, 22)
	for i := 0; i < len(stateOfMECs); i++ {
		stateOfMECs[i] = make([]int, 4)
	}

	json.Unmarshal(resp, &stateOfMECs)
	return stateOfMECs, nil

}

func GetMaskFromNMT(endpoint string, app model.MECApp) ([]int, error) {

	mask := make([]int, 22)

	responseBody, err := postHttpRespBody(endpoint, app)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch mask from nmt: %v", err)
	}

	err = json.Unmarshal(responseBody, &mask)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshall mask: %v", err)
	}

	return mask, nil
}

func GetCapacityCLusterFromNMT(endpoint string) (int, error) {

	var capacity int

	responseBody, err := getHTTPRespBody(endpoint)
	if err != nil {
		return 0, fmt.Errorf("cannot fetch mask from nmt: %v", err)
	}

	err = json.Unmarshal(responseBody, &capacity)
	if err != nil {
		return 0, fmt.Errorf("cannot unmarshall mask: %v", err)
	}

	return capacity, nil
}

func postHttpRespBody(url string, data interface{}) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error: marshaling failed")
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Errorf("Could not make POST request. reason: %v", err)
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

// func getHTTPRespBody(url string) (io.ReadCloser, error) {
func getHTTPRespBody(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		err := fmt.Errorf("HTTP GET failed for URL %s.\nError: %s\n", url, err)
		log.Errorf("%v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("HTTP GET returned status code %s for URL %s.\n", resp.Status, url)
		log.Errorf("%v", err)
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	return b, nil
}
