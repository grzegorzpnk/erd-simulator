package api

import (
	"10.254.188.33/matyspi5/erd/pkg/obs/src/pkg/latency"
	"10.254.188.33/matyspi5/erd/pkg/obs/src/pkg/observability"

	"github.com/gorilla/mux"
)

// NewRouter returns rew MUX ROUTER with registered handlers
func NewRouter(ksmClient *observability.ClustersInfo, ltcClient *latency.MockClient) *mux.Router {

	var handler apiHandler
	handler.SetClients(*ksmClient, *ltcClient)

	router := mux.NewRouter().PathPrefix("/v1").Subrouter()

	baseUrl := "/obs"
	ksmUrl := baseUrl + "/ksm"
	ltcUrl := baseUrl + "/ltc"

	ksmClusterUrl := ksmUrl + "/provider/{provider}/cluster/{cluster}"

	// SAMPLE URL: http://localhost:8282/v1/obs/ksm/provider/orange/cluster/mec02/memory-requests

	router.HandleFunc(ksmClusterUrl+"/cpu", handler.getCpuInfoHandler).Methods("GET")
	router.HandleFunc(ksmClusterUrl+"/cpu/requests", handler.getCpuRequestsHandler).Methods("GET")
	router.HandleFunc(ksmClusterUrl+"/cpu/limits", handler.getCpuLimitsHandler).Methods("GET")
	router.HandleFunc(ksmClusterUrl+"/cpu/allocatable", handler.getCpuAllocHandler).Methods("GET")
	router.HandleFunc(ksmClusterUrl+"/cpu/utilization", handler.getCpuReqUtilizationHandler).Methods("GET")

	router.HandleFunc(ksmClusterUrl+"/memory", handler.getMemInfoHandler).Methods("GET")
	router.HandleFunc(ksmClusterUrl+"/memory/requests", handler.getMemRequestsHandler).Methods("GET")
	router.HandleFunc(ksmClusterUrl+"/memory/limits", handler.getMemLimitsHandler).Methods("GET")
	router.HandleFunc(ksmClusterUrl+"/memory/allocatable", handler.getMemAllocHandler).Methods("GET")
	router.HandleFunc(ksmClusterUrl+"/memory/utilization", handler.getMemReqUtilizationHandler).Methods("GET")

	//router.HandleFunc(ksmClusterUrl+"/provider/{provider}/cluster/{cluster}/cpu-allocatable", handler.getCpuAllocHandler).Methods("GET")
	//router.HandleFunc(ksmClusterUrl+"/provider/{provider}/cluster/{cluster}/cpu-requests", handler.getCpuReqHandler).Methods("GET")
	//router.HandleFunc(ksmClusterUrl+"/provider/{provider}/cluster/{cluster}/cpu-limits", handler.getCpuLimHandler).Methods("GET")
	//router.HandleFunc(ksmClusterUrl+"/provider/{provider}/cluster/{cluster}/memory-used", handler.getMemUsedHandler).Methods("GET")
	//router.HandleFunc(ksmClusterUrl+"/provider/{provider}/cluster/{cluster}/memory-allocatable", handler.getMemAllocHandler).Methods("GET")
	//router.HandleFunc(ksmClusterUrl+"/provider/{provider}/cluster/{cluster}/memory-requests", handler.getMemReqHandler).Methods("GET")
	//router.HandleFunc(ksmClusterUrl+"/provider/{provider}/cluster/{cluster}/memory-limits", handler.getMemLimHandler).Methods("GET")

	// SAMPLE URL: http://localhost:8282/v1/obs/ltc/cell/1/mec/edge-provider+mec1/latency-ms
	router.HandleFunc(ltcUrl+"/source/{source-node}/target/{target-node}/latency-ms", handler.getLatencyHandler).Methods("GET")

	return router
}
