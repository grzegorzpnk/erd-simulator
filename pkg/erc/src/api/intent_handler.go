package api

import (
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/errs"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/model"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/module"
	"encoding/json"
	"fmt"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/apierror"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/validation"
	"io"
	"net/http"
)

type ResponseBody struct {
	Provider string `json:"provider"`
	Cluster  string `json:"cluster"`
}

type intentHandler struct {
	client module.SmartPlacementIntentManager
}

//replace below if you are running locally
var ErJSONFile = "./json-schemas/intent.json"

//var ErJSONFile = "./../../json-schemas/intent.json"

func (h intentHandler) handleSmartPlacementIntentHeuristic(w http.ResponseWriter, r *http.Request) {
	var i model.SmartPlacementIntent

	isValid := validateRequestBody(w, r, &i, ErJSONFile)
	if !isValid {
		return
	}

	mec, err := h.client.ServeSmartPlacementIntentHeuristic(false, i)

	if err != nil {
		msg := fmt.Sprintf("Search failed. Reason: %v", err.Error())
		sendResponse(w, msg, http.StatusInternalServerError)
		return
	} else {
		if mec.Identity.Provider == i.CurrentPlacement.Provider && mec.Identity.Cluster == i.CurrentPlacement.Cluster {
			sendResponse(w, "Relocation redundant. Skipping...", http.StatusNotModified)
			return
		} else {
			body := ResponseBody{
				Provider: mec.Identity.Provider,
				Cluster:  mec.Identity.Cluster,
			}

			sendResponse(w, body, http.StatusOK)
		}
	}
}

func (h intentHandler) handleSmartPlacementIntentEarHeuristic(w http.ResponseWriter, r *http.Request) {
	var i model.SmartPlacementIntent

	isValid := validateRequestBody(w, r, &i, ErJSONFile)
	if !isValid {
		return
	}

	mec, err := h.client.ServeSmartPlacementIntentHeuristic(true, i)

	if err != nil {
		if err.Error() == errs.ERR_CLUSTER_OK.Error() {
			msg := fmt.Sprintf("Skipping search. Reason: current cluster meets the requirements[%+v].", i.CurrentPlacement)
			sendResponse(w, msg, http.StatusNotModified)
			return
		}

		msg := fmt.Sprintf("Search failed. Reason: %v", err.Error())
		sendResponse(w, msg, http.StatusInternalServerError)
		return
	} else {
		if mec.Identity.Provider == i.CurrentPlacement.Provider && mec.Identity.Cluster == i.CurrentPlacement.Cluster {
			sendResponse(w, "Relocation redundant. Skipping...", http.StatusNotModified)
			return
		} else {
			body := ResponseBody{
				Provider: mec.Identity.Provider,
				Cluster:  mec.Identity.Cluster,
			}

			sendResponse(w, body, http.StatusOK)
		}
	}
}

func (h intentHandler) handleSmartPlacementIntentOptimal(w http.ResponseWriter, r *http.Request) {
	var i model.SmartPlacementIntent

	isValid := validateRequestBody(w, r, &i, ErJSONFile)
	if !isValid {
		return
	}

	mec, err := h.client.ServeSmartPlacementIntentOptimal(i)

	if err != nil {
		msg := fmt.Sprintf("Search failed. Reason: %v", err.Error())
		sendResponse(w, msg, http.StatusInternalServerError)
		return
	} else {
		//todo: replace with asking nmt instead of taking data from intent
		if mec.Identity.Provider == i.CurrentPlacement.Provider && mec.Identity.Cluster == i.CurrentPlacement.Cluster {
			sendResponse(w, "Relocation redundant. Skipping...", http.StatusOK)
			return
		} else {
			body := ResponseBody{
				Provider: mec.Identity.Provider,
				Cluster:  mec.Identity.Cluster,
			}

			sendResponse(w, body, http.StatusOK)
		}
	}
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
