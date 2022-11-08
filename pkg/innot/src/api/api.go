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
	router.HandleFunc(subscriptionsUrl+"/paths", handler.handleGetUsersPaths).Methods("GET")
	router.HandleFunc(unsubUrl, handler.unsubscribeByEndpointHandler).Methods("POST")
	router.HandleFunc(unsubByIdUrl, handler.unsubscribeByIdHandler).Methods("POST")

	router.HandleFunc(handleSubUrl, handler.handleSubscriptionHandler).Methods("POST")
	router.HandleFunc(subscriptionsUrl, handler.getAllSubscriptionsHandler).Methods("GET")

	// API which will be used to collect information about experiments
	resultsCollectionUrl := "/results"
	failedUrl := resultsCollectionUrl + "/relocation-failed/inc/{type}"
	successfulUrl := resultsCollectionUrl + "/relocation-successful/inc/{type}"
	redundantUrl := resultsCollectionUrl + "/relocation-redundant/inc/{type}"
	skippedUrl := resultsCollectionUrl + "/relocation-skipped/inc/{type}"
	resetUrl := resultsCollectionUrl + "/reset"

	failedTimeUrl := resultsCollectionUrl + "/relocation-failed/time"
	successfulTimeUrl := resultsCollectionUrl + "/relocation-successful/time"
	redundantTimeUrl := resultsCollectionUrl + "/relocation-redundant/time"
	skippedTimeUrl := resultsCollectionUrl + "/relocation-skipped/time"

	router.HandleFunc(failedUrl, handler.relocationFailedHandler).Methods("POST")
	router.HandleFunc(successfulUrl, handler.relocationSuccessfulHandler).Methods("POST")
	router.HandleFunc(redundantUrl, handler.relocationRedundantHandler).Methods("POST")
	router.HandleFunc(skippedUrl, handler.relocationSkippedHandler).Methods("POST")

	router.HandleFunc(failedTimeUrl, handler.failedTimeHandler).Methods("POST")
	router.HandleFunc(successfulTimeUrl, handler.successfulTimeHandler).Methods("POST")
	router.HandleFunc(redundantTimeUrl, handler.redundantTimeHandler).Methods("POST")
	router.HandleFunc(skippedTimeUrl, handler.skippedTimeHandler).Methods("POST")

	router.HandleFunc(resetUrl, handler.resetHandler).Methods("POST")
	router.HandleFunc(resultsCollectionUrl, handler.getResultsHandler).Methods("GET")
	router.HandleFunc(resultsCollectionUrl+"/csv", handler.getResultsCSVHandler).Methods("GET")

	return router
}
