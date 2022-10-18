// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package lcm

import (
	spi "10.254.188.33/matyspi5/erd/pkg/lcm-workflow/src/model"
	"10.254.188.33/matyspi5/erd/pkg/lcm-workflow/src/types"

	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	ti "gitlab.com/project-emco/core/emco-base/src/workflowmgr/pkg/module"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
)

const TaskQueue = "LCM_TASK_Q"

const (
	ApplyPhase UpdatePhase = iota
	DeletePhase
)

var portList = GenerateNotificationPorts()

func GetPorts(np []NotificationPort) (int, int) {
	for {
		index := rand.Intn(len(np) - 1)
		if !np[index].Used {
			np[index].Used = true
			return np[index].Port, np[index].NodePort
		}
	}

}

// TODO REVISIT Copied from EMCO as import leads to conflicts
type GenericPlacementIntent struct {
	MetaData GenIntentMetaData `json:"metadata"`
}

// GenIntentMetaData has name, description, userdata1, userdata2
type GenIntentMetaData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	UserData1   string `json:"userData1"`
	UserData2   string `json:"userData2"`
}

// AppIntent has two components - metadata, spec
type AppIntent struct {
	MetaData MetaData `json:"metadata,omitempty"`
	Spec     SpecData `json:"spec,omitempty"`
}

// MetaData has - name, description, userdata1, userdata2
type MetaData struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	UserData1   string `json:"userData1,omitempty"`
	UserData2   string `json:"userData2,omitempty"`
}

// SpecData consists of appName and intent
type SpecData struct {
	AppName string      `json:"app,omitempty"`
	Intent  IntentStruc `json:"intent,omitempty"`
}

type IntentStruc struct {
	AllOfArray []AllOf `json:"allOf,omitempty"`
	AnyOfArray []AnyOf `json:"anyOf,omitempty"`
}

// AllOf consists if ProviderName, ClusterName, ClusterLabelName and AnyOfArray. Any of   them can be empty
type AllOf struct {
	ProviderName     string  `json:"clusterProvider,omitempty"`
	ClusterName      string  `json:"cluster,omitempty"`
	ClusterLabelName string  `json:"clusterLabel,omitempty"`
	AnyOfArray       []AnyOf `json:"anyOf,omitempty"`
}

// AnyOf consists of Array of ProviderName & ClusterLabelNames
type AnyOf struct {
	ProviderName     string `json:"clusterProvider,omitempty"`
	ClusterName      string `json:"cluster,omitempty"`
	ClusterLabelName string `json:"clusterLabel,omitempty"`
}
type UpdatePhase int8

type AppNameDetails struct {
	AppName       string
	AppIntentName string
	Phase         UpdatePhase
	PrimaryIntent IntentStruc
}

type WorkflowParams struct {
	InParams map[string]string
	// map indexed by generic placement intent name
	AppsNameDetails      map[string][]AppNameDetails
	NotifyUrl            string
	NewCellId            types.CellId
	SmartPlacementIntent spi.SmartPlacementIntent
	OptimalCluster       Cluster
	ErWfIntent           ti.WorkflowIntent
	ListenerPort         int
	ListenerNodePort     int
}

func (mp *WorkflowParams) GetParamByKey(key string) string {
	return mp.InParams[key]
}

func (mp *WorkflowParams) GetInnotUrl() string {
	return mp.InParams["innotUrl"]
}

// GetOrchestratorGrpcEndpoint gRPC endpoint for Orchestrator
func (mp *WorkflowParams) GetOrchestratorGrpcEndpoint() string {
	return mp.InParams["emcoOrchStatusEndpoint"]
}

// GetClmEndpoint is endpoint for cluster manager microservice
func (mp *WorkflowParams) GetClmEndpoint() string {
	return mp.InParams["emcoClmEndpoint"]
}

func (mp *WorkflowParams) buildWfMgrURL() string {
	url := mp.InParams["emcoWfMgrURL"]
	url += "/v2/projects/" + mp.InParams["project"]
	url += "/composite-apps/" + mp.InParams["compositeApp"]
	url += "/" + mp.InParams["compositeAppVersion"]
	url += "/deployment-intent-groups/" + mp.InParams["deploymentIntentGroup"]
	url += "/temporal-workflow-intents"

	return url
}

func (mp *WorkflowParams) buildDigURL() string {
	url := mp.InParams["emcoOrchEndpoint"]
	url += "/v2/projects/" + mp.InParams["project"]
	url += "/composite-apps/" + mp.InParams["compositeApp"]
	url += "/" + mp.InParams["compositeAppVersion"]
	url += "/deployment-intent-groups/" + mp.InParams["deploymentIntentGroup"]

	return url
}

func (mp *WorkflowParams) buildGenericPlacementIntentsURL() string {
	url := mp.buildDigURL()
	url += "/generic-placement-intents"

	return url
}

func buildAppIntentsURL(gpiURL string, gpiName string) string {
	url := gpiURL + "/" + gpiName + "/app-intents"
	return url
}

// func getHttpRespBody(url string) (io.ReadCloser, error) {
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

type SubMsg struct {
	Subscription struct {
		EventList []struct {
			Type          string `json:"type"`
			ImmediateFlag bool   `json:"immediateFlag"`
			AreaList      []struct {
				PresenceInfo struct {
					TrackingAreaList []struct {
						PlmnID struct {
							Mcc string `json:"mcc"`
							Mnc string `json:"mnc"`
						} `json:"plmnId"`
						Tac string `json:"tac"`
					} `json:"trackingAreaList"`
				} `json:"presenceInfo"`
				SNssai []struct {
					Sd  string `json:"sd"`
					Sst int    `json:"sst"`
				} `json:"sNssai"`
			} `json:"areaList"`
			UdmDetectInd     bool `json:"udmDetectInd"`
			PresenceInfoList struct {
			} `json:"presenceInfoList"`
			MaxResponseTime int `json:"maxResponseTime"`
			TargetArea      struct {
				TaList []struct {
					PlmnID struct {
						Mcc string `json:"mcc"`
						Mnc string `json:"mnc"`
					} `json:"plmnId"`
					Tac string `json:"tac"`
				} `json:"taList"`
			} `json:"targetArea"`
		} `json:"eventList"`
		AnyUE          bool   `json:"anyUE"`
		EventNotifyURI string `json:"eventNotifyUri"`
		Options        struct {
			Trigger string `json:"trigger"`
		} `json:"options"`
	} `json:"subscription"`
}

func postHttpRespBody(url string, data interface{}) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error: marshaling failed")
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Could not make post request. reason: %v\n", err)
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

func postHttp(url string, data interface{}) (string, error) {
	body, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error: marshaling failed")
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Could not make post request. reason: %v\n", err)
	}

	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		var res string
		json.NewDecoder(resp.Body).Decode(&res)
		return res, err
	}
	return "", err
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func RandString() string {
	b := make([]byte, 10)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func GenerateListenerEndpoint(baseUrl, notifyType string) string {
	randStr := RandString()
	url := fmt.Sprintf("/%s/%s/%s/notify", baseUrl, notifyType, randStr)
	return url
}

type apiHandler struct {
	channel chan types.CellId
}

func (h *apiHandler) serveNotification(w http.ResponseWriter, r *http.Request) {
	var info types.CellChangedInfo
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		fmt.Errorf("error while decoding: %v\n", err)
	}

	log.Printf("Notification reason: %v, cell id: %v", info.Reason, info.Cell)

	h.channel <- info.Cell

}

func generateSubscriptionBody() types.AmfEventSubscription {
	boole := true

	amfEventArea := types.AmfEventArea{
		LadnInfo: &types.LadnInfo{
			Ladn:     "",
			Presence: nil,
		},
		PresenceInfo: &types.PresenceInfo{
			EcgiList:            nil,
			GlobalRanNodeIdList: nil,
			NcgiList:            nil,
			PraId:               nil,
			PresenceState:       nil,
			TrackingAreaList: &[]types.Tai{{
				PlmnId: types.PlmnId{"208", "93"},
				Tac:    "000001",
			}},
		},
	}
	amfEvent := types.AmfEvent{
		AreaList:                 &[]types.AmfEventArea{amfEventArea},
		ImmediateFlag:            &boole,
		LocationFilterList:       nil,
		SubscribedDataFilterList: nil,
		Type:                     "LOCATION_REPORT",
	}

	body := types.AmfEventSubscription{
		AnyUE:                         &boole,
		EventList:                     &[]types.AmfEvent{amfEvent},
		EventNotifyUri:                "http://localhost/workflow-listener/cell-changed/ABCDEFGHIJ/notify",
		Gpsi:                          nil,
		GroupId:                       nil,
		NfId:                          uuid.UUID{},
		NotifyCorrelationId:           "",
		Options:                       &types.AmfEventMode{Trigger: "ONE_TIME"},
		Pei:                           nil,
		SubsChangeNotifyCorrelationId: nil,
		SubsChangeNotifyUri:           nil,
		Supi:                          nil,
	}
	return body
}

type Cluster struct {
	Provider string `json:"provider"`
	Cluster  string `json:"cluster"`
}

type NotificationPort struct {
	Port     int
	NodePort int
	Used     bool
}

func GenerateNotificationPorts() []NotificationPort {
	var portList []NotificationPort
	for i := 1; i <= 200; i++ {
		portList = append(portList, NotificationPort{
			Port:     8200 + i,
			NodePort: 32200 + i,
			Used:     false,
		})
	}
	return portList
}
