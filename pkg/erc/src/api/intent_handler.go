// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package api

import (
	"encoding/json"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/apierror"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/validation"
	"io"
	"net/http"

	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/model"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/module"
)

type intentHandler struct {
	client module.SmartPlacementIntentManager
}

var ErJSONFile string = "./json-schemas/intent.json"

func (h intentHandler) handleSmartPlacementIntent(w http.ResponseWriter, r *http.Request) {
	var i model.SmartPlacementIntent

	isValid := validateRequestBody(w, r, &i, ErJSONFile)
	if !isValid {
		return
	}

	mecFqdn, err := h.client.ServeSmartPlacementIntentOutsideEMCO(i)
	if err != nil {
		//handleError(w, map[string]string{}, err, i)
		sendResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendResponse(w, mecFqdn, http.StatusOK)
}

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
