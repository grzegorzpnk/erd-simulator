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

func CallPlacementController(intent model.SmartPlacementIntent) (*model.Cluster, error) {
	//	log.Printf("CallPlacementController: function start\n")
	var resp model.Cluster

	plcCtrlUrl := buildPlcCtrlURL()
	data := intent

	responseBody, err := postHttpRespBody(plcCtrlUrl, data)
	if err != nil {
		log.Errorf("[ERROR] Placement Controller returned error: %v. Relocation for APP[%v] failed.", err, intent.Metadata.Name)
		return nil, err
	}

	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		log.Errorf("error occured while unmarshaling: %v. Resp body: %v", err, string(responseBody))
		return nil, err
	}

	var cluster model.Cluster
	cluster = resp

	return &cluster, nil
}

func buildPlcCtrlURL() string {

	url := config.GetConfiguration().ERCEndpoint
	url += "/v2/erc/smart-placement-intents/optimal-mec/optimal"

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

func buildOrchestratorURL(app model.MECApp, cluster model.Cluster) string {
	//orchestrator/relocate/old-cluster/{old-cluster}/new-cluster/{new-cluster}/application"

	url := config.GetConfiguration().NMTEndpoint
	url += "/v1/orchestrator/relocate/old-cluster/"
	url += app.ClusterId + "/new-cluster/"
	url += cluster.Cluster + "/application"

	return url

}
