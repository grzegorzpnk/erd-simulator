package api

import (
	"fmt"
	"reflect"

	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/module"
	"github.com/gorilla/mux"
)

// NewRouter creates a router that registers the various routes.
// If the mockClient parameter is not nil, the router is configured with a mock handler.
func NewRouter(mockClient interface{}) *mux.Router {
	const baseURL string = "/erc/smart-placement-intents"
	const spiURL string = baseURL + "/optimal-mec"
	r := mux.NewRouter().PathPrefix("/v2").Subrouter()
	c := module.NewClient()
	h := intentHandler{
		client: setClient(c.SmartPlacementIntent, mockClient).(module.SmartPlacementIntentManager),
	}

	r.HandleFunc(spiURL+"/optimal", h.handleSmartPlacementIntentOptimal).Methods("POST")
	r.HandleFunc(spiURL+"/heuristic", h.handleSmartPlacementIntentHeuristic).Methods("POST")
	r.HandleFunc(spiURL+"/ear-heuristic", h.handleSmartPlacementIntentEarHeuristic).Methods("POST")

	//RESULTS
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

	r.HandleFunc(failedUrl, h.relocationFailedHandler).Methods("POST")
	r.HandleFunc(successfulUrl, h.relocationSuccessfulHandler).Methods("POST")
	r.HandleFunc(redundantUrl, h.relocationRedundantHandler).Methods("POST")
	r.HandleFunc(skippedUrl, h.relocationSkippedHandler).Methods("POST")

	r.HandleFunc(failedTimeUrl, h.failedTimeHandler).Methods("POST")
	r.HandleFunc(successfulTimeUrl, h.successfulTimeHandler).Methods("POST")
	r.HandleFunc(redundantTimeUrl, h.redundantTimeHandler).Methods("POST")
	r.HandleFunc(skippedTimeUrl, h.skippedTimeHandler).Methods("POST")

	r.HandleFunc(resetUrl, h.resetHandler).Methods("POST")
	r.HandleFunc(resultsCollectionUrl, h.getResultsHandler).Methods("GET")
	r.HandleFunc(resultsCollectionUrl+"/csv", h.getResultsCSVHandler).Methods("GET")

	return r
}

// setClient set the client and its corresponding manager interface.
// If the mockClient parameter is not nil and implements the manager interface
// corresponding to the client return the mockClient. Otherwise, return the client.
func setClient(client, mockClient interface{}) interface{} {
	switch cl := client.(type) {
	case *module.SmartPlacementIntentClient:
		if mockClient != nil && reflect.TypeOf(mockClient).Implements(reflect.TypeOf((*module.SmartPlacementIntentManager)(nil)).Elem()) {
			c, ok := mockClient.(module.SmartPlacementIntentManager)
			if ok {
				return c
			}
		}
	default:
		fmt.Printf("unknown type %T\n", cl)
	}
	return client
}
