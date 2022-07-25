package api

import (
	"github.com/gorilla/mux"
)

// NewRouter returns rew MUX ROUTER with registered handlers
func NewRouter() *mux.Router {

	var handler apiHandler

	router := mux.NewRouter().PathPrefix("/v1").Subrouter()

	baseUrl := "/intermediate-notifier"
	subUrl := baseUrl + "/subscribe"
	unsubUrl := baseUrl + "/unsubscribe"
	unsubByIdUrl := unsubUrl + "/{subscriptionId}"
	getAllUrl := baseUrl + "/get-all"

	router.HandleFunc(subUrl, handler.subscribeHandler).Methods("POST")
	router.HandleFunc(unsubUrl, handler.unsubscribeByEndpointHandler).Methods("POST")
	router.HandleFunc(unsubByIdUrl, handler.unsubscribeByIdHandler).Methods("POST")
	router.HandleFunc(getAllUrl, handler.getAllSubscriptions).Methods("GET")

	return router
}
