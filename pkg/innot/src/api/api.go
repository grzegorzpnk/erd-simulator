package api

import (
	"github.com/gorilla/mux"
)

// NewRouter returns rew MUX ROUTER with registered handlers
func NewRouter() *mux.Router {

	var handler apiHandler
	handler.SetClients()

	router := mux.NewRouter().PathPrefix("/v1").Subrouter()

	baseCollectionUrl := "/intermediate-notifier"
	subUrl := baseCollectionUrl + "/subscribe"
	unsubUrl := baseCollectionUrl + "/unsubscribe"
	unsubByIdUrl := unsubUrl + "/{subscriptionId}"
	subscriptionsUrl := baseCollectionUrl + "/subscriptions"
	handleSubUrl := subscriptionsUrl + "/{subscription-id}/handle"

	router.HandleFunc(subUrl, handler.subscribeHandler).Methods("POST")
	router.HandleFunc(unsubUrl, handler.unsubscribeByEndpointHandler).Methods("POST")
	router.HandleFunc(unsubByIdUrl, handler.unsubscribeByIdHandler).Methods("POST")

	router.HandleFunc(handleSubUrl, handler.handleSubscriptionHandler).Methods("POST")
	router.HandleFunc(subscriptionsUrl, handler.getAllSubscriptionsHandler).Methods("GET")

	// API which will be used to collect information about experiments
	resultsCollectionUrl := "/results"
	failedUrl := resultsCollectionUrl + "/relocation-failed"
	successfulUrl := resultsCollectionUrl + "/relocation-successful"
	redundantUrl := resultsCollectionUrl + "/relocation-redundant"

	router.HandleFunc(failedUrl, handler.relocationFailedHandler).Methods("POST")
	router.HandleFunc(successfulUrl, handler.relocationSuccessfulHandler).Methods("POST")
	router.HandleFunc(redundantUrl, handler.relocationRedundantHandler).Methods("POST")
	router.HandleFunc(resultsCollectionUrl, handler.getResultsHandler).Methods("GET")

	return router
}
