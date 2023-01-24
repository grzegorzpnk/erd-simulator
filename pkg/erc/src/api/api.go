package api

import (
	"fmt"
	"reflect"

	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/module"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/results"
	"github.com/gorilla/mux"
)

// NewRouter creates a router that registers the various routes.
// If the mockClient parameter is not nil, the router is configured with a mock handler.
func NewRouter(mockClient interface{}, resultClient *results.Client) *mux.Router {
	const baseURL string = "/erc/smart-placement-intents"
	const spiURL string = baseURL + "/optimal-mec"
	r := mux.NewRouter().PathPrefix("/v2").Subrouter()
	c := module.NewClient()
	h := intentHandler{
		client:       setClient(c.SmartPlacementIntent, mockClient).(module.SmartPlacementIntentManager),
		resultClient: resultClient,
	}

	r.HandleFunc(spiURL+"/optimal", h.handleSmartPlacementIntentOptimal).Methods("POST")
	r.HandleFunc(spiURL+"/heuristic", h.handleSmartPlacementIntentHeuristic).Methods("POST")
	r.HandleFunc(spiURL+"/ear-heuristic", h.handleSmartPlacementIntentEarHeuristic).Methods("POST")

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
