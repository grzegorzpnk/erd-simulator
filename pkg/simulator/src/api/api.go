package api

import (
	"github.com/gorilla/mux"
	"simu/src/pkg/model"
)

var r *mux.Router

func NewRouter(simuClient *model.SimuClient) *mux.Router {

	var handler apiHandler
	handler.SetClients(*simuClient)

	r := mux.NewRouter().PathPrefix("/v1").Subrouter()

	baseUrl := ""

	r.HandleFunc(baseUrl+"/simu/experiment/", handler.conductExperiment).Methods("POST")

	//test API
	r.HandleFunc(baseUrl+"/simu/users", handler.getUsers).Methods("GET")
	//temporary function to check if the list of declared apps fetch from NMT == list updated in simu
	//r.HandleFunc(baseUrl+"/simu/compareLists", handler.CompareLists).Methods("GET")

	return r

}
