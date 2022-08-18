package api

import (
	"10.254.188.33/matyspi5/pmc/src/pkg/observability"

	"github.com/gorilla/mux"
)

// NewRouter returns rew MUX ROUTER with registered handlers
func NewRouter(client *observability.ClustersInfo) *mux.Router {

	var handler apiHandler
	handler.SetKsmClient(*client)

	router := mux.NewRouter().PathPrefix("/v1").Subrouter()

	baseUrl := "/pmc"
	ksmUrl := baseUrl + "/ksm"
	router.HandleFunc(baseUrl, handler.baseHandler).Methods("GET")

	// SAMPLE URL: http://localhost:8282/v1/pmc/ksm/provider/orange/cluster/meh02/get-mem-lim
	router.HandleFunc(ksmUrl+"/provider/{provider}/cluster/{cluster}/get-cpu-req", handler.getCpuReqHandler).Methods("GET")
	router.HandleFunc(ksmUrl+"/provider/{provider}/cluster/{cluster}/get-cpu-lim", handler.getCpuLimHandler).Methods("GET")
	router.HandleFunc(ksmUrl+"/provider/{provider}/cluster/{cluster}/get-mem-req", handler.getMemReqHandler).Methods("GET")
	router.HandleFunc(ksmUrl+"/provider/{provider}/cluster/{cluster}/get-mem-lim", handler.getMemLimHandler).Methods("GET")

	return router
}
