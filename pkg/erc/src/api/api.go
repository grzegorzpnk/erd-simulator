// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

// Package api defines all the routes and their associated handler functions.
// This example implements two HTTP methods.
// It registers two routes to create and retrieve the intents associated with a deployment group.
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
	//const baseURL string = "/projects/{project}/composite-apps/{compositeApp}/{compositeAppVersion}/deployment-intent-groups/{deploymentIntentGroup}/smartPlacementIntents"
	const baseURL string = "/erc/smart-placement-intents"
	r := mux.NewRouter().PathPrefix("/v2").Subrouter()
	c := module.NewClient()
	h := intentHandler{
		client: setClient(c.SmartPlacementIntent, mockClient).(module.SmartPlacementIntentManager),
	}

	//r.HandleFunc(baseURL, h.handleSmartPlacementIntentCreate).Methods("POST")
	//r.HandleFunc(baseURL+"/{smartPlacementIntent}", h.handleSmartPlacementIntentGet).Methods("GET")

	r.HandleFunc(baseURL+"/optimal-mec", h.handleSmartPlacementIntentOutsideEMCO).Methods("GET")

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
