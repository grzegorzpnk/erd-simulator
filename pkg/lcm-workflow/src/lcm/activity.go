// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package lcm

import (
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

func GenerateSmartPlacementIntent(ctx context.Context, migParam WorkflowParams) (*WorkflowParams, error) {
	log.Printf("GenerateSmartPlacementIntent: activity start\n")

	targetAppName := migParam.GetParamByKey("targetAppName")
	priorityLevel := migParam.GetParamByKey("appPriorityLevel")

	latMax, _ := strconv.ParseFloat(migParam.GetParamByKey("latencyMax"), 64)
	cpuUtilMax, _ := strconv.ParseFloat(migParam.GetParamByKey("cpuUtilizationMax"), 64)
	memUtilMax, _ := strconv.ParseFloat(migParam.GetParamByKey("memUtilizationMax"), 64)

	latWeight, _ := strconv.ParseFloat(migParam.GetParamByKey("latencyWeight"), 64)
	cpuWeight, _ := strconv.ParseFloat(migParam.GetParamByKey("cpuUtilizationWeight"), 64)
	memWeight, _ := strconv.ParseFloat(migParam.GetParamByKey("memUtilizationWeight"), 64)

	spIntent := spi.SmartPlacementIntent{
		Metadata: spi.Metadata{
			Name:        targetAppName + "-er-intent",
			Description: fmt.Sprintf("Edge Relocation Intent for app: %s", targetAppName),
		},
		Spec: spi.SmartPlacementIntentSpec{
			AppName: targetAppName,
			SmartPlacementIntentData: spi.SmartPlacementIntentStruct{
				TargetCell:    migParam.NewCellId,
				PriorityLevel: getPriorityLevel(priorityLevel),
				ConstraintsList: spi.Constraints{
					LatencyMax:        latMax,
					CpuUtilizationMax: cpuUtilMax,
					MemUtilizationMax: memUtilMax,
				},
				ParametersWeights: spi.Weights{
					LatencyWeight:        latWeight,
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
	var resp Cluster

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
	log.Printf("CallPlacementController: Optimal cluster = %+v.\n", resp)

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
					RetryPolicy:                              &temporal.RetryPolicy{MaximumAttempts: 10},
					CronSchedule:                             "",
					Memo:                                     nil,
					SearchAttributes:                         nil,
				},
				WfParams: eta.WorkflowParams{
					ActivityOpts: map[string]wf.ActivityOptions{"all-activities": {StartToCloseTimeout: 60000000000,
						HeartbeatTimeout: 50000000000, RetryPolicy: &temporal.RetryPolicy{InitialInterval: 10, MaximumAttempts: 3}}},
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
