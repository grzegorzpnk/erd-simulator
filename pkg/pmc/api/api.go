package api

import (
	"github.com/gorilla/mux"
)

// NewRouter returns rew MUX ROUTER with registered handlers
func NewRouter() *mux.Router {

	var handler apiHandler

	router := mux.NewRouter().PathPrefix("/v1").Subrouter()

	baseUrl := "/prometheus-metrics-controller"

	router.HandleFunc(baseUrl, handler.baseHandler).Methods("GET")

	return router
}
