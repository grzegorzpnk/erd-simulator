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
	r.HandleFunc(baseUrl+"/simu/single-experiment", handler.conductSingleExperiment).Methods("POST")
	r.HandleFunc(baseUrl+"/simu/experiment", handler.conductExperiment).Methods("POST")

	//Heuristic
	r.HandleFunc(baseUrl+"/simu/experiment-icc", handler.conductExperimentICC).Methods("POST")
	r.HandleFunc(baseUrl+"/simu/experiment-icc-tunning", handler.conductExperimentICCTunning).Methods("POST")
	r.HandleFunc(baseUrl+"/simu/experiment-icc-tunning-iter", handler.conductExperimentICCTunningIter).Methods("POST")

	//ML
	r.HandleFunc(baseUrl+"/simu/experiment-globcom", handler.conductExperimentGlobcom).Methods("POST")
	//r.HandleFunc(baseUrl+"/simu/experiment-ml-compare", handler.conductExperimentMLCompare).Methods("POST")
	//r.HandleFunc(baseUrl+"/simu/ml-state", handler.stateTest).Methods("POST")
	//r.HandleFunc(baseUrl+"/simu/ml-experiment", handler.conductMLExperiment).Methods("POST")

	//Charts and results
	r.HandleFunc(baseUrl+"/results/all", handler.getAllResults).Methods("GET")
	r.HandleFunc(baseUrl+"/results/charts/icc-heuristic", handler.generateICCHeuristicChart).Methods("GET")
	r.HandleFunc(baseUrl+"/results/charts/icc-heuristic-tunning", handler.generateICCTunningHeuristicChart).Methods("GET")
	r.HandleFunc(baseUrl+"/results/chartsml", handler.generateMLChart).Methods("GET")
	r.HandleFunc(baseUrl+"/results/test", handler.test).Methods("GET")

	//test API
	r.HandleFunc(baseUrl+"/simu/users", handler.getUsers).Methods("GET")
	//temporary function to check if the list of declared apps fetch from NMT == list updated in simu
	//r.HandleFunc(baseUrl+"/simu/compareLists", handler.CompareLists).Methods("GET")

	return r

}
