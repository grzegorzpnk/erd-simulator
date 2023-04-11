package api

//
//func executeMLExperiment(h *apiHandler, expIndex, subExpIndex int, experimentType results.ExperimentType) bool {
//
//	experimentN := "[EXPERIMENT " + strconv.Itoa(expIndex+1) + "." + strconv.Itoa(subExpIndex+1) + "] "
//	//generate number of user to move
//	id := h.generateUserToMove() //USER==APP
//
//	// select new position for selected user and add new position to UserPath
//	app := h.SimuClient.GetApps(id)
//	h.generateTargetCellId(app)
//	log.Infof(experimentN+"User(app) with ID: %v [current mec: %v] moved FROM cell: %v, towards cell: %v", app.Id, app.ClusterId, app.UserPath[len(app.UserPath)-2], app.UserLocation)
//
//	spi, err := GenerateMLSmartPlacementIntent(*app, experimentType)
//	if err != nil {
//		log.Errorf("Cannot generate SPI: %v", err.Error())
//		return false
//	}
//
//	//send request to ML ctrl to select new position
//	cluster, err := CallMLClient(spi)
//
//	if err != nil {
//		log.Warnf("ML Placement ctrl has returned status : %v", err.Error())
//		log.Warnf(experimentN + "stopped, NO RELOCATION, going to next iteration")
//		return true
//	}
//
//	if cluster.Cluster == app.ClusterId {
//		log.Infof(experimentN+"Selected redundant cluster: %v -> missing relocation", cluster.Cluster)
//		return true
//	}
//
//	log.Infof(experimentN+"Selected new cluster: %v", cluster.Cluster)
//
//	//generate request to orchestrator
//	err2 := sendRelocationRequest(*app, *cluster)
//	if err2 != nil {
//		log.Errorf("Cannot relocate app! Error: %v", err2.Error())
//	} else {
//		log.Infof(experimentN + "Application has been relocated in nmt")
//
//		//update cluster in internal app list
//		app.ClusterId = cluster.Cluster
//	}
//
//	return true
//}

//
//// SpaceAPP (for single app)  : 1) Required mvCPU 2) required Memory 3) Required Latency 4) Current MEC 5) Current RAN
//func GenerateMLSmartPlacementIntent(app model.MECApp, experimentType results.ExperimentType) (model.MLSmartPlacementIntent, error) {
//	//log.Printf("GenerateSmartPlacementIntent: activity start\n")
//
//	var spIntent model.MLSmartPlacementIntent
//	var state model.State
//
//	clusterID, _ := convertMECNameToID(app.ClusterId)
//	userLocation, _ := strconv.Atoi(app.UserLocation)
//
//	appState := [1][5]int{{
//		determineReqRes(int(app.Requirements.RequestedCPU)),
//		determineReqRes(int(app.Requirements.RequestedMEMORY)),
//		determineStateofAppLatReq(int(app.Requirements.RequestedLatency)),
//		clusterID,
//		userLocation}}
//
//	url := buildNMTCurrentStateEndpoint()
//	fmt.Printf("asking for MECs config url:, %v     |", url)
//
//	mecState, err := GetMECsStateFromNMT(url)
//	if err != nil {
//		return spIntent, err
//	}
//
//	state = model.State{
//		SpaceAPP:  appState,
//		SpaceMECs: mecState,
//	}
//
//	if experimentType == results.ExpMLNonMasked {
//		spIntent = model.MLSmartPlacementIntent{State: state}
//
//	} else if experimentType == results.ExpMLMasked {
//		mask, err := GenerateMLMask(app)
//		if err != nil {
//			log.Errorf("Cannot generate Mask: %v", err.Error())
//		}
//		spIntent = model.MLSmartPlacementIntent{State: state, CurrentMask: mask}
//	} else {
//		fmt.Errorf("Invalid type of experiment: %v ", experimentType)
//		return spIntent, err
//	}
//
//	log.Infof("GenerateSmartPlacementIntent: intent = %+v\n", spIntent)
//
//	return spIntent, nil
//}
//
//func GenerateMLMask(app model.MECApp) ([]int, error) {
//
//	mask := make([]int, 22)
//
//	url := buildNMTMaskEndpoint()
//	fmt.Printf("asking for Mask url:, %v     |", url)
//
//	mask, err := GetMaskFromNMT(url, app)
//	if err != nil {
//		return mask, err
//	}
//
//	log.Infof("GeneratedMask: %+v\n", mask)
//
//	return mask, nil
//}
//
//func determineReqRes(reqRes int) int {
//	resMap := map[int]int{
//		500:  1,
//		600:  2,
//		700:  3,
//		800:  4,
//		900:  5,
//		1000: 6,
//	}
//	if val, ok := resMap[reqRes]; ok {
//		return val
//	}
//	return 0
//}
//
//func determineStateofAppLatReq(latValue int) int {
//	latMap := map[int]int{
//		10: 1,
//		15: 2,
//		30: 3,
//	}
//	if val, ok := latMap[latValue]; ok {
//		return val
//	}
//	return 0
//}
//
//func CallMLClient(intent model.MLSmartPlacementIntent) (*model.Cluster, error) {
//	//	log.Printf("CallPlacementController: function start\n")
//	var respClusterID string
//
//	plcCtrlUrl := buildMLPlcCtrlURL()
//	data := intent
//
//	responseBody, err := postHttpRespBody(plcCtrlUrl, data)
//	if err != nil {
//		return nil, fmt.Errorf("relocation for APP not done: %v", err)
//	}
//
//	err = json.Unmarshal(responseBody, &respClusterID)
//	if err != nil {
//		return nil, fmt.Errorf("relocation for APP not done: %v", err)
//	}
//
//	cluster := model.Cluster{
//		Cluster:  "orange",
//		Provider: "mec" + respClusterID,
//	}
//
//	return &cluster, nil
//}
//
//func buildMLPlcCtrlURL() string {
//	url := config.GetConfiguration().MLClientEndpoint
//	url += "/ermodel/get-prediction"
//
//	return url
//}

//
//func buildNMTMaskEndpoint() string {
//	url := config.GetConfiguration().NMTEndpoint
//	url += "/v1/topology/ml/get-mask"
//
//	return url
//}
//
//func buildNMTCurrentStateEndpoint() string {
//	///topology/ml/get-state
//	url := config.GetConfiguration().NMTEndpoint
//	url += "/v1/topology/ml/get-state"
//	return url
//}
//
//func GetMECsStateFromNMT(endpoint string) ([][]int, error) {
//
//	resp, err := http.Get(endpoint)
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		log.Fatalln(err)
//		return nil, err
//	}
//
//	//Convert the body to type [][]int
//	stateOfMECs := make([][]int, 22)
//	for i := 0; i < len(stateOfMECs); i++ {
//		stateOfMECs[i] = make([]int, 5)
//	}
//
//	json.Unmarshal(body, &stateOfMECs)
//	return stateOfMECs, nil
//
//}
//
//func GetMaskFromNMT(endpoint string, app model.MECApp) ([]int, error) {
//
//	mask := make([]int, 22)
//
//	responseBody, err := postHttpRespBody(endpoint, app)
//	if err != nil {
//		return nil, fmt.Errorf("cannot fetch mask from nmt: %v", err)
//	}
//
//	err = json.Unmarshal(responseBody, &mask)
//	if err != nil {
//		return nil, fmt.Errorf("cannot unmarshall mask: %v", err)
//	}
//
//	return mask, nil
//}

//
//func convertMECNameToID(s string) (int, error) {
//	if !strings.HasPrefix(s, "mec") {
//		return 0, fmt.Errorf("invalid input format: %s", s)
//	}
//	numStr := strings.TrimPrefix(s, "mec")
//	num, err := strconv.Atoi(numStr)
//	if err != nil {
//		return 0, fmt.Errorf("failed to convert %s to int: %v", numStr, err)
//	}
//	return num, nil
//}

//func checkMLExperimentType(inputType string) (results.ExperimentType, error) {
//	if strings.ToLower(inputType) == "ml-masked" {
//		return results.ExpMLMasked, nil
//
//	} else if strings.ToLower(inputType) == "ml-non-masked" {
//		return results.ExpMLNonMasked, nil
//	}
//
//	return results.ExpNotExists, fmt.Errorf("provided experiment type [%v] in not an option: %v", inputType, results.GetExpTypes())
//}
