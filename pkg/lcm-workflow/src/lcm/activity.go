// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package lcm

import (
	"10.254.188.33/matyspi5/erd/pkg/lcm-workflow/src/model"
	spi "10.254.188.33/matyspi5/erd/pkg/lcm-workflow/src/model"
	"10.254.188.33/matyspi5/erd/pkg/lcm-workflow/src/types"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	eta "gitlab.com/project-emco/core/emco-base/src/workflowmgr/pkg/emcotemporalapi"
	ti "gitlab.com/project-emco/core/emco-base/src/workflowmgr/pkg/module"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	wf "go.temporal.io/sdk/workflow"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func SubCellChangedNotification(ctx context.Context, migParam WorkflowParams) (*WorkflowParams, error) {

	//log.Printf("SubCellChangedNotification got params: %#v\n", migParam)
	log.Printf("SubCellChangedNotification: activity start\n")

	migParam.NotifyUrl = GenerateListenerEndpoint("workflow-listener", "cell-changed")

	innotUrl := migParam.GetInnotUrl()

	log.Printf("\nSubCellChangedNotification: InnotUrl = %s\n", innotUrl)

	if migParam.ListenerPort == 0 {
		port, nodePort := GetPorts(portList)
		migParam.ListenerPort = port
		migParam.ListenerNodePort = nodePort
	}

	host := fmt.Sprintf("http://%v:%v", os.Getenv("NOTIFICATION_NODE_ADDR"), migParam.ListenerNodePort)

	data := generateSubscriptionBody()
	data.EventNotifyUri = host + migParam.NotifyUrl

	log.Printf("SubCellChangedNotification: EventNorifyUri = %v\n", data.EventNotifyUri)

	_, err := postHttp(innotUrl, data)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("SubCellChangedNotification: Subscribed for notifications.")
	}

	return &migParam, nil
}

func GetCellChangedNotification(ctx context.Context, migParam WorkflowParams) (*WorkflowParams, error) {
	log.Printf("GetCellChangedNotification: activity start\n")

	var waitNotification = make(chan types.CellId)
	router := mux.NewRouter()
	handler := apiHandler{waitNotification}

	router.HandleFunc(migParam.NotifyUrl, handler.serveNotification).Methods("POST")

	httpServer := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%v", migParam.ListenerPort),
	}

	go func() {
		migParam.NewCellId = <-waitNotification

		fmt.Printf("GetCellChangedNotification: Got notification. New CELL_ID: %v\n.", migParam.NewCellId)
		_ = httpServer.Shutdown(context.Background())
	}()

	fmt.Printf("GetCellChangedNotification: Listening for notification..\n")
	_ = httpServer.ListenAndServe()

	return &migParam, nil
}

func DiscoverCurrentCluster(ctx context.Context, migParam WorkflowParams) (*WorkflowParams, error) {
	log.Printf("DiscoverCurrentCluster: activity start\n")

	gpiUrl := buildGenericPlacementIntentsURL(migParam.InParams)

	respBody, err := getHttpRespBody(gpiUrl)
	if err != nil {
		return nil, err
	}

	var gpIntents []GenericPlacementIntent

	if err := json.Unmarshal(respBody, &gpIntents); err != nil {
		decodeErr := fmt.Errorf("Failed to decode GET responde body for URL %s.\n"+"Decoder error: %#v\n", gpiUrl, err)
		_, _ = fmt.Fprintf(os.Stderr, decodeErr.Error())
		if err != nil {
			return nil, err
		}
		return nil, decodeErr
	}

	// For now assume that there is only one GenericPlacementIntent
	appIntentsUrl := buildAppIntentsURL(gpiUrl, gpIntents[0].MetaData.Name)

	respBody, err = getHttpRespBody(appIntentsUrl)
	if err != nil {
		return nil, err
	}

	var appIntents []AppIntent
	if err := json.Unmarshal(respBody, &appIntents); err != nil {
		decodeErr := fmt.Errorf("Failed to decode GET responde body for URL %s.\nDecoder error: %#v\n", appIntentsUrl, err)
		_, _ = fmt.Fprintf(os.Stderr, decodeErr.Error())
		return nil, decodeErr
	}
	fmt.Printf("\nDiscoverCurrentCluster: body = %#v\n", appIntents)

	for _, appIntent := range appIntents {
		if strings.ToLower(appIntent.Spec.AppName) == strings.ToLower(migParam.InParams["targetAppName"]) {
			// Need to consider AllOfArray, AnyOfArray, ClusterNames, ClusterLabels
			// Assume the simplest case: we consider only 1 application, deployed using ClusterName and the Intent is
			// in the AllOf array, AnyOf array is empty
			if len(appIntent.Spec.Intent.AllOfArray) > 0 {
				for _, elem := range appIntent.Spec.Intent.AllOfArray {
					if elem.ClusterName != "" {
						migParam.CurrentCluster = model.Cluster{
							Provider: elem.ProviderName,
							Cluster:  elem.ClusterName,
						}
					}
				}
			}
			if len(appIntent.Spec.Intent.AnyOfArray) > 0 {
				for _, elem := range appIntent.Spec.Intent.AnyOfArray {
					if elem.ClusterName != "" || elem.ClusterLabelName != "" {
						log.Printf("[WARN]DiscoverCurrentCluster: There was items in the AnyOf array but skipped.")
					}
				}
			}
		}
	}

	return &migParam, nil
}

func GenerateSmartPlacementIntent(ctx context.Context, migParam WorkflowParams) (*WorkflowParams, error) {
	log.Printf("GenerateSmartPlacementIntent: activity start\n")

	targetAppName := migParam.GetParamByKey("targetAppName")

	appCpuReq, err := strconv.ParseFloat(migParam.GetParamByKey("appCpuReq"), 64)
	if err != nil {
		log.Printf("Could not parse appCpuReq[%v] into float", migParam.GetParamByKey("appCpuReq"))
		appCpuReq = 0.0
	}

	appMemReq, err := strconv.ParseFloat(migParam.GetParamByKey("appMemReq"), 64)
	if err != nil {
		log.Printf("Could not parse appMemReq[%v] into float", migParam.GetParamByKey("appMemReq"))
		appCpuReq = 0.0
	}

	latMax, err := strconv.ParseFloat(migParam.GetParamByKey("latencyMax"), 64)
	if err != nil {
		log.Printf("Could not parse latencyMax[%v] into float", migParam.GetParamByKey("latencyMax"))
		appCpuReq = 0.0
	}

	cpuUtilMax, err := strconv.ParseFloat(migParam.GetParamByKey("cpuUtilizationMax"), 64)
	if err != nil {
		log.Printf("Could not parse cpuUtilizationMax[%v] into float", migParam.GetParamByKey("cpuUtilizationMax"))
		appCpuReq = 0.0
	}

	memUtilMax, err := strconv.ParseFloat(migParam.GetParamByKey("memUtilizationMax"), 64)
	if err != nil {
		log.Printf("Could not parse memUtilizationMax[%v] into float", migParam.GetParamByKey("memUtilizationMax"))
		appCpuReq = 0.0
	}

	latWeight, err := strconv.ParseFloat(migParam.GetParamByKey("latencyWeight"), 64)
	if err != nil {
		log.Printf("Could not parse latencyWeight[%v] into float", migParam.GetParamByKey("latencyWeight"))
		appCpuReq = 0.0
	}

	resWeight, err := strconv.ParseFloat(migParam.GetParamByKey("resourcesWeight"), 64)
	if err != nil {
		log.Printf("Could not parse resourcesWeight[%v] into float", migParam.GetParamByKey("resourcesWeight"))
		appCpuReq = 0.0
	}

	cpuWeight, err := strconv.ParseFloat(migParam.GetParamByKey("cpuUtilizationWeight"), 64)
	if err != nil {
		log.Printf("Could not parse cpuUtilizationWeight[%v] into float", migParam.GetParamByKey("cpuUtilizationWeight"))
		appCpuReq = 0.0
	}

	memWeight, err := strconv.ParseFloat(migParam.GetParamByKey("memUtilizationWeight"), 64)
	if err != nil {
		log.Printf("Could not parse memUtilizationWeight[%v] into float", migParam.GetParamByKey("memUtilizationWeight"))
		appCpuReq = 0.0
	}

	spIntent := spi.SmartPlacementIntent{
		Metadata: spi.Metadata{
			Name:        targetAppName + "-er-intent",
			Description: fmt.Sprintf("Edge Relocation Intent for app: %s", targetAppName),
		},
		CurrentPlacement: model.Cluster{
			Provider: migParam.CurrentCluster.Provider,
			Cluster:  migParam.CurrentCluster.Cluster,
		},
		Spec: spi.SmartPlacementIntentSpec{
			AppName: targetAppName,
			SmartPlacementIntentData: spi.SmartPlacementIntentStruct{
				TargetCell: migParam.NewCellId,
				AppCpuReq:  appCpuReq,
				AppMemReq:  appMemReq,
				ConstraintsList: spi.Constraints{
					LatencyMax:        latMax,
					CpuUtilizationMax: cpuUtilMax,
					MemUtilizationMax: memUtilMax,
				},
				ParametersWeights: spi.Weights{
					LatencyWeight:        latWeight,
					ResourcesWeight:      resWeight,
					CpuUtilizationWeight: cpuWeight,
					MemUtilizationWeight: memWeight,
				},
			},
		},
	}
	log.Printf("GenerateSmartPlacementIntent: intent = %+v\n", spIntent)
	migParam.SmartPlacementIntent = spIntent
	return &migParam, nil
}

func CallPlacementController(ctx context.Context, migParam WorkflowParams) (*WorkflowParams, error) {
	log.Printf("CallPlacementController: activity start\n")
	var resp model.Cluster

	plcCtrlUrl := migParam.GetParamByKey("plcControllerUrl")
	data := migParam.SmartPlacementIntent

	responseBody, err := postHttpRespBody(plcCtrlUrl, data)
	if err != nil {
		log.Printf("[ERROR] Placement Controller returned error: %v. Relocation for APP[%v] failed.", err.Error(), migParam.GetParamByKey("targetAppName"))
		return &migParam, err
	}

	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		log.Printf("error occured while unmarshaling: %v. Resp body: %v", err, string(responseBody))
		return &migParam, err
	}

	migParam.OptimalCluster = resp

	log.Printf("CallPlacementController: Optimal cluster = provider{%v}, cluster={%v}.\n", resp.Provider, resp.Cluster)

	return &migParam, nil
}

func GenerateRelocateWfIntent(ctx context.Context, migParam WorkflowParams) (*WorkflowParams, error) {
	log.Printf("GenerateRelocateWfIntent: activity start\n")

	appName := migParam.GetParamByKey("targetAppName")
	clientName := migParam.GetParamByKey("rClientName")
	clientPort, _ := strconv.Atoi(migParam.GetParamByKey("rClientPort"))
	wfClientName := migParam.GetParamByKey("rWfClientName")

	erWfIntent := ti.WorkflowIntent{
		Metadata: ti.Metadata{
			Name:        appName + "er-intent-" + RandString(),
			Description: "Edge Relocation WfIntent",
		},
		Spec: ti.WorkflowIntentSpec{
			WfClientSpec: ti.WfClientSpec{
				WfClientEndpointName: clientName,
				WfClientEndpointPort: clientPort,
			},
			WfTemporalSpec: eta.WfTemporalSpec{
				WfClientName: wfClientName,
				WfStartOpts: client.StartWorkflowOptions{
					ID:                                       appName + RandString(),
					TaskQueue:                                "RELOCATE_TASK_Q",
					WorkflowExecutionTimeout:                 0,
					WorkflowRunTimeout:                       0,
					WorkflowTaskTimeout:                      0,
					WorkflowIDReusePolicy:                    0,
					WorkflowExecutionErrorWhenAlreadyStarted: false,
					RetryPolicy:                              &temporal.RetryPolicy{MaximumAttempts: 1},
					CronSchedule:                             "",
					Memo:                                     nil,
					SearchAttributes:                         nil,
				},
				WfParams: eta.WorkflowParams{
					ActivityOpts: map[string]wf.ActivityOptions{"all-activities": {StartToCloseTimeout: 60000000000,
						HeartbeatTimeout: 50000000000, RetryPolicy: &temporal.RetryPolicy{InitialInterval: 10, MaximumAttempts: 1}}},
					ActivityParams: map[string]map[string]string{"all-activities": {
						"emcoOrchEndpoint":       migParam.GetParamByKey("emcoOrchEndpoint"),
						"emcoOrchStatusEndpoint": migParam.GetParamByKey("emcoOrchStatusEndpoint"),
						"emcoClmEndpoint":        migParam.GetParamByKey("emcoClmEndpoint"),
						"project":                migParam.GetParamByKey("project"),
						"compositeApp":           migParam.GetParamByKey("compositeApp"),
						"compositeAppVersion":    migParam.GetParamByKey("compositeAppVersion"),
						"deploymentIntentGroup":  migParam.GetParamByKey("deploymentIntentGroup"),
						"targetClusterProvider":  migParam.OptimalCluster.Provider,
						"targetClusterName":      migParam.OptimalCluster.Cluster,
						"targetAppName":          migParam.GetParamByKey("targetAppName"),
					}},
				},
			},
		},
	}

	migParam.ErWfIntent = erWfIntent
	log.Printf("GenerateRelocateWfIntent: generated WfIntent = %+v\n", migParam.ErWfIntent)

	return &migParam, nil
}

func CallTemporalWfController(ctx context.Context, migParam WorkflowParams) (*WorkflowParams, error) {
	log.Printf("CallTemporalWfController: activity start\n")

	createWfUrl := migParam.buildWfMgrURL()

	_, err := postHttp(createWfUrl, migParam.ErWfIntent)
	if err != nil {
		return nil, err
	}
	log.Printf("CallTemporalWfController: generated WfIntent = %+v\n", migParam.ErWfIntent)
	fmt.Printf("CallTemporalWfController: created ER Wf Intent.\n")

	startWfUrl := createWfUrl + "/" + migParam.ErWfIntent.Metadata.Name + "/start"

	startBody, err := postHttp(startWfUrl, "")
	if err != nil {
		return nil, err
	}

	fmt.Printf("CallTemporalWfController: started ER Wf Intent, response body: %v.\n", startBody)

	return &migParam, nil
}
