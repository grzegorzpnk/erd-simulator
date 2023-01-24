package api

import (
	log "10.254.188.33/matyspi5/erd/pkg/erc/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/errs"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/model"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/module"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/results"
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/apierror"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/validation"
	"io"
	"net/http"
	"time"
)

type ResponseBody struct {
	Provider string `json:"provider"`
	Cluster  string `json:"cluster"`
}

type intentHandler struct {
	client       module.SmartPlacementIntentManager
	resultClient *results.Client
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
		body := ResponseBody{
			Provider: mec.Identity.Provider,
			Cluster:  mec.Identity.Cluster,
		}

		sendResponse(w, body, http.StatusOK)
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
		body := ResponseBody{
			Provider: mec.Identity.Provider,
			Cluster:  mec.Identity.Cluster,
		}

		sendResponse(w, body, http.StatusOK)
	}
}

func (h intentHandler) handleSmartPlacementIntentOptimal(w http.ResponseWriter, r *http.Request) {
	var i model.SmartPlacementIntent

	isValid := validateRequestBody(w, r, &i, ErJSONFile)
	if !isValid {
		return
	}

	//measure time of searching for best cluster
	startTime := time.Now()
	mec, err := h.client.ServeSmartPlacementIntentOptimal(i)
	elapsedTime := time.Since(startTime)
	respBody, err2 := json.Marshal(elapsedTime)
	if err2 != nil {
		log.Warnf("Could not unmarshal: %v. error: %v", elapsedTime, err2)
	}

	if err != nil {
		// EXPERIMENTS: remove later
		if err.Error() == errs.ERR_CLUSTER_OK.Error() {
			http.Post(fmt.Sprintf("http://10.254.185.44:32137/v1/results/relocation-skipped/inc/%v",
				i.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax), "application/json", bytes.NewBuffer([]byte{}))
			http.Post("http://10.254.185.44:32137/v1/results/relocation-skipped/time", "application/json", bytes.NewBuffer(respBody))
			sendResponse(w, err.Error(), http.StatusNoContent)
			return
		}
		// EXPERIMENTS: remove later
		http.Post(fmt.Sprintf("http://10.254.185.44:32137/v1/results/relocation-failed/inc/%v",
			i.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax), "application/json", bytes.NewBuffer([]byte{}))
		http.Post("http://10.254.185.44:32137/v1/results/relocation-failed/time", "application/json", bytes.NewBuffer(respBody))
		sendResponse(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		if mec.Identity.Provider == i.CurrentPlacement.Provider && mec.Identity.Cluster == i.CurrentPlacement.Cluster {
			http.Post(fmt.Sprintf("http://10.254.185.44:32137/v1/results/relocation-redundant/inc/%v",
				i.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax), "application/json", bytes.NewBuffer([]byte{}))
			http.Post("http://10.254.185.44:32137/v1/results/relocation-redundant/time", "application/json", bytes.NewBuffer(respBody))
			sendResponse(w, "Relocation redundant. Skipping...", http.StatusNoContent)
			return
		} else {
			http.Post(fmt.Sprintf("http://10.254.185.44:32137/v1/results/relocation-successful/inc/%v",
				i.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax), "application/json", bytes.NewBuffer([]byte{}))
			http.Post("http://10.254.185.44:32137/v1/results/relocation-successful/time", "application/json", bytes.NewBuffer(respBody))
			body := ResponseBody{
				Provider: mec.Identity.Provider,
				Cluster:  mec.Identity.Cluster,
			}
			sendResponse(w, body, http.StatusOK)
		}
	}

	//////////////////////////////////////////////////////////////////////////////////////////////////

	//let's find out what type of searching it was: succesfull, failed, redundant or if-skipped
	if err != nil {
		msg := fmt.Sprintf("Search has returned failed. Reason: %v", err.Error())
		sendResponse(w, msg, http.StatusInternalServerError)
		return
	} else {
		if mec.Identity.Provider == i.CurrentPlacement.Provider && mec.Identity.Cluster == i.CurrentPlacement.Cluster {

			h.resultClient.Results.IncRedundant(i.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax)

			/*http.Post(fmt.Sprintf("http://10.254.185.44:32137/v1/results/relocation-redundant/inc/%v",
					i.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax), "application/json", bytes.NewBuffer([]byte{}))
				http.Post("http://10.254.185.44:32137/v1/results/relocation-redundant/time", "application/json", bytes.NewBuffer(respBody))
				sendResponse(w, "Relocation redundant. Skipping...", http.StatusNoContent)
				return
			}*/
		} else {

			http.Post(fmt.Sprintf("http://10.254.185.44:32137/v1/results/relocation-successful/inc/%v",
				i.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax), "application/json", bytes.NewBuffer([]byte{}))
			http.Post("http://10.254.185.44:32137/v1/results/relocation-successful/time", "application/json", bytes.NewBuffer(respBody))
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
