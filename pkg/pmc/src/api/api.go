package api

import (
	"10.254.188.33/matyspi5/pmc/src/pkg/latency"
	"10.254.188.33/matyspi5/pmc/src/pkg/observability"

	"github.com/gorilla/mux"
)

// NewRouter returns rew MUX ROUTER with registered handlers
func NewRouter(ksmClient *observability.ClustersInfo, ltcClient *latency.MockClient) *mux.Router {

	var handler apiHandler
	handler.SetClients(*ksmClient, *ltcClient)

	router := mux.NewRouter().PathPrefix("/v1").Subrouter()

	baseUrl := "/pmc"
	ksmUrl := baseUrl + "/ksm"
	ltcUrl := baseUrl + "/ltc"

	// SAMPLE URL: http://localhost:8282/v1/pmc/ksm/provider/orange/cluster/meh02/get-mem-req
	router.HandleFunc(ksmUrl+"/provider/{provider}/cluster/{cluster}/get-cpu-req", handler.getCpuReqHandler).Methods("GET")
	router.HandleFunc(ksmUrl+"/provider/{provider}/cluster/{cluster}/get-cpu-lim", handler.getCpuLimHandler).Methods("GET")
	router.HandleFunc(ksmUrl+"/provider/{provider}/cluster/{cluster}/get-mem-req", handler.getMemReqHandler).Methods("GET")
	router.HandleFunc(ksmUrl+"/provider/{provider}/cluster/{cluster}/get-mem-lim", handler.getMemLimHandler).Methods("GET")

	// SAMPLE URL: http://localhost:8282/v1/pmc/ltc/cell/1/meh/edge-provider+meh01/get-latency-ms
	router.HandleFunc(ltcUrl+"/cell/{cell-id}/meh/{meh-id}/get-latency-ms", handler.getLatencyHandler).Methods("GET")

	return router
}
