// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package lcm

import (
	eri "10.254.188.33/matyspi5/erd/pkg/lcm-workflow/src/module"
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

var waitNotification = make(chan types.CellId)

func SubCellChangedNotification(ctx context.Context, migParam MigParam) (*MigParam, error) {

	log.Printf("SubCellChangedNotification got params: %#v\n", migParam)
	log.Printf("SubCellChangedNotification: activity start\n")

	migParam.NotifyUrl = GenerateListenerEndpoint("workflow-listener", "cell-changed")

	innotUrl := migParam.GetInnotUrl()

	log.Printf("\nSubCellChangedNotification: InnotUrl = %s\n", innotUrl)

	host := fmt.Sprintf("http://10.254.185.48:%v", os.Getenv("NOTIFICATION_NODE_PORT"))

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

func GetCellChangedNotification(ctx context.Context, migParam MigParam) (*MigParam, error) {
	log.Printf("GetCellChangedNotification: activity start\n")

	router := mux.NewRouter()
	router.HandleFunc(migParam.NotifyUrl, serveNotification).Methods("POST")

	httpServer := &http.Server{
		Handler: router,
		Addr:    ":8585",
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

func GenerateERIntent(ctx context.Context, migParam MigParam) (*MigParam, error) {
	log.Printf("GenerateERIntent: activity start\n")

	targetAppName := migParam.GetParamByKey("targetAppName")
	priorityLevel := migParam.GetParamByKey("appPriorityLevel")

	latMax, _ := strconv.ParseFloat(migParam.GetParamByKey("latencyMax"), 64)
	cpuUtilMax, _ := strconv.ParseFloat(migParam.GetParamByKey("cpuUtilizationMax"), 64)
	memUtilMax, _ := strconv.ParseFloat(migParam.GetParamByKey("memUtilizationMax"), 64)

	latWeight, _ := strconv.ParseFloat(migParam.GetParamByKey("latencyWeight"), 64)
	cpuWeight, _ := strconv.ParseFloat(migParam.GetParamByKey("cpuUtilizationWeight"), 64)
	memWeight, _ := strconv.ParseFloat(migParam.GetParamByKey("memUtilizationWeight"), 64)

	erIntent := eri.ErIntent{
		MetaData: eri.MetaData{
			Name:        targetAppName + "-er-intent",
			Description: fmt.Sprintf("Edge Relocation Intent for app: %s", targetAppName),
		},
		Spec: eri.SpecData{
			AppName: targetAppName,
			Intent: eri.IntentStruct{
				PriorityLevel: getPriorityLevel(priorityLevel),
				ConstraintsList: eri.Constraints{
					LatencyMax:        latMax,
					CpuUtilizationMax: cpuUtilMax,
					MemUtilizationMax: memUtilMax,
				},
				ParametersWeights: eri.Weights{
					LatencyWeight:        latWeight,
					CpuUtilizationWeight: cpuWeight,
					MemUtilizationWeight: memWeight,
				},
			},
		},
	}
	log.Printf("GenerateERIntent: intent = %+v\n", erIntent)
	migParam.ErIntent = erIntent
	return &migParam, nil
}

func CallPlacementController(ctx context.Context, migParam MigParam) (*MigParam, error) {
	log.Printf("CallPlacementController: activity start\n")
	var resp Cluster

	plcCtrlUrl := migParam.GetParamByKey("plcControllerUrl")
	data := migParam.ErIntent

	responseBody, err := postHttp(plcCtrlUrl, data)
	if err != nil {
		log.Fatal(err)
	}

	_ = json.Unmarshal([]byte(responseBody), &resp)

	migParam.OptimalCluster = resp

	log.Printf("CallPlacementController: Optimal cluster = provider{%v}, cluster={%v}.\n", resp.ProviderName, resp.ClusterName)

	return &migParam, nil
}

func CreateTemporalErWfIntent(ctx context.Context, migParam MigParam) (*MigParam, error) {
	log.Printf("CreateTemporalErWfIntent: activity start\n")

	appName := migParam.GetParamByKey("appName")
	clientName := migParam.GetParamByKey("rClientName")
	clientPort, _ := strconv.Atoi(migParam.GetParamByKey("rClientPort"))
	wfClientName := migParam.GetParamByKey("rWfClientName")

	erWfIntent := ti.WorkflowIntent{
		Metadata: ti.Metadata{
			Name:        appName + "er-intent",
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
						HeartbeatTimeout: 50000000000, RetryPolicy: &temporal.RetryPolicy{InitialInterval: 10}}},
					ActivityParams: map[string]map[string]string{"all-activities": {
						"emcoOrchEndpoint":       migParam.GetParamByKey("emcoOrchEndpoint"),
						"emcoOrchStatusEndpoint": migParam.GetParamByKey("emcoOrchStatusEndpoint"),
						"emcoClmEndpoint":        migParam.GetParamByKey("emcoClmEndpoint"),
						"project":                migParam.GetParamByKey("project"),
						"compositeApp":           migParam.GetParamByKey("compositeApp"),
						"compositeAppVersion":    migParam.GetParamByKey("compositeAppVersion"),
						"deploymentIntentGroup":  migParam.GetParamByKey("deploymentIntentGroup"),
						"targetClusterProvider":  migParam.OptimalCluster.ProviderName,
						"targetClusterName":      migParam.OptimalCluster.ClusterName,
						"targetAppName":          migParam.GetParamByKey("targetAppName"),
					}},
				},
			},
		},
	}

	migParam.ErWfIntent = erWfIntent
	log.Printf("CreateTemporalErWfIntent: generated WfIntent = %+v\n", migParam.ErWfIntent)

	createWfUrl := migParam.buildWfMgrURL()

	_, err := postHttp(createWfUrl, erWfIntent)
	if err != nil {
		return nil, err
	}
	fmt.Printf("CreateTemporalErWfIntent: created ER Wf Intent.\n")

	startWfUrl := createWfUrl + "/" + erWfIntent.Metadata.Name + "/start"

	startBody, err := postHttp(startWfUrl, "")
	if err != nil {
		return nil, err
	}

	fmt.Printf("CreateTemporalErWfIntent: started ER Wf Intent, resp body: %v.\n", startBody)

	return &migParam, nil
}
