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

	// SAMPLE URL: http://localhost:8282/v1/obs/ksm/provider/orange/cluster/meh02/memory-requests
	router.HandleFunc(ksmUrl+"/provider/{provider}/cluster/{cluster}/cpu-requests", handler.getCpuReqHandler).Methods("GET")
	router.HandleFunc(ksmUrl+"/provider/{provider}/cluster/{cluster}/cpu-limits", handler.getCpuLimHandler).Methods("GET")
	router.HandleFunc(ksmUrl+"/provider/{provider}/cluster/{cluster}/memory-requests", handler.getMemReqHandler).Methods("GET")
	router.HandleFunc(ksmUrl+"/provider/{provider}/cluster/{cluster}/memory-limits", handler.getMemLimHandler).Methods("GET")

	// SAMPLE URL: http://localhost:8282/v1/obs/ltc/cell/1/meh/edge-provider+meh01/latency-ms
	router.HandleFunc(ltcUrl+"/cell/{cell-id}/meh/{meh-id}/latency-ms", handler.getLatencyHandler).Methods("GET")

	return router
}
