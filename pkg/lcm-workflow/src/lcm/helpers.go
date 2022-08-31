// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package lcm

import (
	"10.254.188.33/matyspi5/erd/pkg/lcm-workflow/src/types"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const MigTaskQueue = "LCM_TASK_Q"

const (
	ApplyPhase UpdatePhase = iota
	DeletePhase
)

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

type MigParam struct {
	InParams map[string]string
	// map indexed by generic placement intent name
	AppsNameDetails map[string][]AppNameDetails
	InnotUrl        string
	NotifyUrl       string
	Router          *http.Server
}

// GetOrchestratorGrpcEndpoint gRPC endpoint for Orchestrator
func GetOrchestratorGrpcEndpoint(mp MigParam) string {
	return mp.InParams["emcoOrchStatusEndpoint"]
}

// GetClmEndpoint is endpoint for cluster manager microservice
func GetClmEndpoint(mp MigParam) string {
	return mp.InParams["emcoClmEndpoint"]
}

func buildDigURL(params map[string]string) string {
	url := params["emcoOrchEndpoint"]
	url += "/v2/projects/" + params["project"]
	url += "/composite-apps/" + params["compositeApp"]
	url += "/" + params["compositeAppVersion"]
	url += "/deployment-intent-groups/" + params["deploymentIntentGroup"]

	return url
}

func buildGenericPlacementIntentsURL(params map[string]string) string {
	url := buildDigURL(params)
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

type testStruct struct {
	value string
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

func postHttp(url string, data types.AmfEventSubscription) error {
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
		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		fmt.Println(res)
		return err
	}
	return err
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func RandString() string {
	rand.Seed(time.Now().UnixNano())
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
