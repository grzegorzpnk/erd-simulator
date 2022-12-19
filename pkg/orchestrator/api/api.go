package api

import (
	"github.com/gorilla/mux"
)

var r *mux.Router

func NewRouter() *mux.Router {

	var handler apiHandler
	handler.SetClients()

	r := mux.NewRouter().PathPrefix("/v1").Subrouter()

	baseUrl := ""

	//Setters:
	r.HandleFunc(baseUrl+"/topology/mecHost", handler.relocateApp).Methods("POST")

	return r

}
