// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package api

import (
	"encoding/json"
	"fmt"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/validation"
	"io"
	"net/http"

	"10.254.188.33/matyspi5/erd/pkg/erc/pkg/model"
	"10.254.188.33/matyspi5/erd/pkg/erc/pkg/module"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/apierror"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"
)

type intentHandler struct {
	client module.SmartPlacementIntentManager
}

// The JSON file defines the allowed fields and values for the SmartPlacementIntent.
// It ensures that the data received from the client is valid before storing it in the database.
var ErJSONFile string = "../json-schemas/intent.json"

func (h intentHandler) handleSmartPlacementIntentOutsideEMCO(w http.ResponseWriter, r *http.Request) {
	var i model.SmartPlacementIntent

	isValid := validateRequestBody(w, r, &i, ErJSONFile)
	if !isValid {
		return
	}

	provider, cluster, err := h.client.ServeSmartPlacementIntentOutsideEMCO(i)
	if err != nil {
		//handleError(w, map[string]string{}, err, i)
		sendResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Provider: %v, cluster: %v", provider, cluster)

	sendResponse(w, provider+"+"+cluster, http.StatusOK)
}

//// handleSEIntentCreate handles the route for creating a new SmartPlacementIntent
//func (h intentHandler) handleSmartPlacementIntentCreate(w http.ResponseWriter, r *http.Request) {
//	var i model.SmartPlacementIntent
//	vars := mux.Vars(r)
//	project := vars["project"]
//	app := vars["compositeApp"]
//	version := vars["compositeAppVersion"]
//	group := vars["deploymentIntentGroup"]
//
//	// Validate the request body before storing it in the database.
//	isValid := validateRequestBody(w, r, &i, ErJSONFile)
//	if !isValid {
//		return
//	}
//
//	// Insert the new SmartPlacementIntent in the database.
//	intent, err := h.client.CreateSmartPlacementIntent(i, project, app, version, group, true) // true - fail if exists
//	if err != nil {
//		handleError(w, vars, err, i)
//		return
//	}
//
//	// Send the response to the client.
//	sendResponse(w, intent, http.StatusCreated)
//}
//
//// handleSmartPlacementIntentGet handles the route for retrieving an SmartPlacementIntent from the database
//func (h intentHandler) handleSmartPlacementIntentGet(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	name := vars["smartPlacementIntent"]
//	project := vars["project"]
//	app := vars["compositeApp"]
//	version := vars["compositeAppVersion"]
//	group := vars["deploymentIntentGroup"]
//
//	// Retrieve the SmartPlacementIntent details from the database.
//	intent, err := h.client.GetSmartPlacementIntent(name, project, app, version, group)
//	if err != nil {
//		handleError(w, vars, err, nil)
//		return
//	}
//
//	// Send the response to the client.
//	sendResponse(w, intent[0], http.StatusOK)
//}

// validateRequestBody validate the request body before storing it in the database
func validateRequestBody(w http.ResponseWriter, r *http.Request, v interface{}, filename string) bool {
	err := json.NewDecoder(r.Body).Decode(&v)
	switch {
	case err == io.EOF:
		logutils.Error("Empty request body",
			logutils.Fields{
				"Error": err})
		http.Error(w, "Empty request body", http.StatusBadRequest)
		// Not a valid request body.
		return false
	case err != nil:
		logutils.Error("Error decoding the request body",
			logutils.Fields{
				"Error": err})
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		// Not a valid request body.
		return false
	}

	// Ensure that the request body matches the schema defined in the JSON file.
	err, httpError := validation.ValidateJsonSchemaData(filename, v)
	if err != nil {
		logutils.Error(err.Error(),
			logutils.Fields{})
		http.Error(w, err.Error(), httpError)
		// Not a valid request body.
		return false
	}

	return true
}

// sendResponse sends an HTTP response to the client with the provided status
func sendResponse(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		logutils.Error("Error encoding response",
			logutils.Fields{
				"Error":    err,
				"Response": v})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleError matches the error with the API errors and
// sends the appropriate message and status to the client
func handleError(w http.ResponseWriter, params map[string]string, err error, mod interface{}) {
	apiErr := apierror.HandleErrors(params, err, mod, apiErrors)
	http.Error(w, apiErr.Message, apiErr.Status)
}