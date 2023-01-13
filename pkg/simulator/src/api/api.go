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

	//Create infrastructure:
	r.HandleFunc(baseUrl+"/simu/mecHost", handler.createMecHostHandler).Methods("POST")

	return r

}
