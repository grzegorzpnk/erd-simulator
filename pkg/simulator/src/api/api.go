package api

import (
	"github.com/gorilla/mux"
	"simu/src/pkg/model"
	"simu/src/pkg/results"
)

var r *mux.Router

func NewRouter(sClient *model.SimuClient, rClient *results.Client) *mux.Router {

	var handler apiHandler
	handler.SetClients(*sClient, *rClient)

	r := mux.NewRouter().PathPrefix("/v1").Subrouter()

	baseUrl := ""
	//r.HandleFunc(baseUrl+"/simu/ml-state", handler.stateTest).Methods("POST")
	//r.HandleFunc(baseUrl+"/simu/ml-experiment", handler.conductMLExperiment).Methods("POST")
	r.HandleFunc(baseUrl+"/simu/single-experiment", handler.conductSingleExperiment).Methods("POST")
	r.HandleFunc(baseUrl+"/simu/experiment", handler.conductExperiment).Methods("POST")
	r.HandleFunc(baseUrl+"/simu/experiment-globcom", handler.conductExperimentGlobcom).Methods("POST")

	//test API
	r.HandleFunc(baseUrl+"/simu/users", handler.getUsers).Methods("GET")
	//temporary function to check if the list of declared apps fetch from NMT == list updated in simu
	//r.HandleFunc(baseUrl+"/simu/compareLists", handler.CompareLists).Methods("GET")

	r.HandleFunc(baseUrl+"/results/all", handler.getAllResults).Methods("GET")
	r.HandleFunc(baseUrl+"/results/charts", handler.generateChartPkg).Methods("GET")

	return r

}
