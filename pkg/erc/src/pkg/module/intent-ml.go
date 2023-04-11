package module

import (
	"10.254.188.33/matyspi5/erd/pkg/erc/src/config"
	log "10.254.188.33/matyspi5/erd/pkg/erc/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func (i *SmartPlacementIntentClient) ServeSmartPlacementIntentML(checkIfMasked bool, intent model.SmartPlacementIntent) (model.MecHost, error) {

	var bestMec model.MecHost

	log.Infof("\n\nReceived request about finding new cluster for app: %v located at Cluster: %v, that moves towards cell: %v\n", intent.Spec.AppName, intent.CurrentPlacement.Cluster, intent.Spec.SmartPlacementIntentData.TargetCell)
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
		log.Warnf(" Could not find optimal cluster for given APP[%v]", intent.Spec.AppName)
		return model.MecHost{}, err
	} else {
		log.Infof("ML CLient has found cluster[%v] for given APP[%v]", bestMec.Identity.Cluster, intent.Spec.AppName)
		return bestMec, nil
	}

}

func CallMLClient(intent model.MLSmartPlacementIntent) (model.MecHost, error) {
	//	log.Printf("CallPlacementController: function start\n")
	var respClusterID string

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
			Cluster:  "orange",
			Provider: "mec" + respClusterID,
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
	//log.Printf("GenerateSmartPlacementIntent: activity start\n")

	var spIntent model.MLSmartPlacementIntent
	var state model.State

	clusterID, _ := convertMECNameToID(intent.CurrentPlacement.Cluster)
	userLocation, _ := strconv.Atoi(string(intent.Spec.SmartPlacementIntentData.TargetCell))

	appState := [1][5]int{{
		determineReqRes(int(intent.Spec.SmartPlacementIntentData.AppCpuReq)),
		determineReqRes(int(intent.Spec.SmartPlacementIntentData.AppMemReq)),
		determineStateofAppLatReq(int(intent.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax)),
		clusterID,
		userLocation}}

	url := buildNMTCurrentStateEndpoint()
	fmt.Printf("asking for MECs config url:, %v     |", url)

	mecState, err := GetMECsStateFromNMT(url)
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
		fmt.Errorf("Invalid type of experiment")
		return spIntent, err
	}

	log.Infof("GenerateMLSmartPlacementIntent: intent = %+v\n", spIntent)

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
	fmt.Printf("asking for Mask url:, %v     |", url)

	mask, err := GetMaskFromNMT(url, app)
	if err != nil {
		return mask, err
	}

	log.Infof("GeneratedMask: %+v\n", mask)

	return mask, nil
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

func buildNMTMaskEndpoint() string {
	url := config.GetConfiguration().NMTEndpoint
	url += "/v1/topology/ml/get-mask"

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
