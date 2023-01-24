package api

import (
	log "10.254.188.33/matyspi5/erd/pkg/erc/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/errs"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/model"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/module"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/results"
	"encoding/json"
	"fmt"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/apierror"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/validation"
	"io"
	"net/http"
	"strconv"
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

	//measure time of searching for best cluster
	startTime := time.Now()
	mec, err := h.client.ServeSmartPlacementIntentHeuristic(false, i)
	elapsedTime := time.Since(startTime)

	if err != nil {
		// EXPERIMENTS: remove later
		if err.Error() == errs.ERR_CLUSTER_OK.Error() {
			h.resultClient.Results.IncSkipped(strconv.FormatFloat(i.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax, 'f', -1, 64))
			h.resultClient.Results.AddSkippedTime(int(elapsedTime.Milliseconds()))
			sendResponse(w, err.Error(), http.StatusNoContent)
			return
		}
		// EXPERIMENTS: remove later
		h.resultClient.Results.IncFailed(strconv.FormatFloat(i.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax, 'f', -1, 64))
		h.resultClient.Results.AddFailedTime(int(elapsedTime.Milliseconds()))
		sendResponse(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		if mec.Identity.Provider == i.CurrentPlacement.Provider && mec.Identity.Cluster == i.CurrentPlacement.Cluster {

			h.resultClient.Results.IncRedundant(strconv.FormatFloat(i.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax, 'f', -1, 64))
			h.resultClient.Results.AddRedundantTime(int(elapsedTime.Milliseconds()))

			sendResponse(w, "redundant cluster.. skipped..", http.StatusNoContent)
		} else {

			h.resultClient.Results.IncSuccessful(strconv.FormatFloat(i.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax, 'f', -1, 64))
			h.resultClient.Results.AddSuccessfulTime(int(elapsedTime.Milliseconds()))

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

	//measure time of searching for best cluster
	startTime := time.Now()
	mec, err := h.client.ServeSmartPlacementIntentHeuristic(true, i)
	elapsedTime := time.Since(startTime)

	if err != nil {
		// EXPERIMENTS: remove later
		if err.Error() == errs.ERR_CLUSTER_OK.Error() {
			h.resultClient.Results.IncSkipped(strconv.FormatFloat(i.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax, 'f', -1, 64))
			h.resultClient.Results.AddSkippedTime(int(elapsedTime.Milliseconds()))
			sendResponse(w, err.Error(), http.StatusNoContent)
			return
		}
		// EXPERIMENTS: remove later
		h.resultClient.Results.IncFailed(strconv.FormatFloat(i.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax, 'f', -1, 64))
		h.resultClient.Results.AddFailedTime(int(elapsedTime.Milliseconds()))
		sendResponse(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		if mec.Identity.Provider == i.CurrentPlacement.Provider && mec.Identity.Cluster == i.CurrentPlacement.Cluster {

			h.resultClient.Results.IncRedundant(strconv.FormatFloat(i.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax, 'f', -1, 64))
			h.resultClient.Results.AddRedundantTime(int(elapsedTime.Milliseconds()))

			sendResponse(w, "redundant cluster.. skipped..", http.StatusNoContent)
		} else {

			h.resultClient.Results.IncSuccessful(strconv.FormatFloat(i.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax, 'f', -1, 64))
			h.resultClient.Results.AddSuccessfulTime(int(elapsedTime.Milliseconds()))

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

	//measure time of searching for best cluster
	startTime := time.Now()
	mec, err := h.client.ServeSmartPlacementIntentOptimal(i)
	elapsedTime := time.Since(startTime)

	//let's find out what type of searching it was: succesfull, failed, redundant or if-skipped
	if err != nil {
		// EXPERIMENTS: remove later
		if err.Error() == errs.ERR_CLUSTER_OK.Error() {
			h.resultClient.Results.IncSkipped(strconv.FormatFloat(i.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax, 'f', -1, 64))
			h.resultClient.Results.AddSkippedTime(int(elapsedTime.Milliseconds()))
			sendResponse(w, err.Error(), http.StatusNoContent)
			return
		}
		// EXPERIMENTS: remove later
		h.resultClient.Results.IncFailed(strconv.FormatFloat(i.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax, 'f', -1, 64))
		h.resultClient.Results.AddFailedTime(int(elapsedTime.Milliseconds()))
		sendResponse(w, err.Error(), http.StatusInternalServerError)
		return
	} else {

		body := ResponseBody{
			Provider: mec.Identity.Provider,
			Cluster:  mec.Identity.Cluster,
		}

		if mec.Identity.Provider == i.CurrentPlacement.Provider && mec.Identity.Cluster == i.CurrentPlacement.Cluster {

			h.resultClient.Results.IncRedundant(strconv.FormatFloat(i.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax, 'f', -1, 64))
			h.resultClient.Results.AddRedundantTime(int(elapsedTime.Milliseconds()))
			sendResponse(w, "redundant cluster.. skipped..", http.StatusOK)

		} else {

			h.resultClient.Results.IncSuccessful(strconv.FormatFloat(i.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax, 'f', -1, 64))
			h.resultClient.Results.AddSuccessfulTime(int(elapsedTime.Milliseconds()))
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

func (h *intentHandler) resetHandler(w http.ResponseWriter, r *http.Request) {

	h.resultClient.Results.Reset()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func (h *intentHandler) getResultsHandler(w http.ResponseWriter, r *http.Request) {

	subs := h.resultClient.Results

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(subs)
	if err != nil {
		log.Error("[API] Error encoding.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *intentHandler) getResultsCSVHandler(w http.ResponseWriter, r *http.Request) {
	var body []string
	subs := h.resultClient.Results

	body = append(body, "strategy,objective,experiment,iteration,edge_relocation_cg,search_successful_cg,search_failed_cg,search_skipped_cg,edge_relocation_v2x,search_successful_v2x,search_failed_v2x,search_skipped_v2x,edge_relocation_uav,search_successful_uav,search_failed_uav,search_skipped_uav")

	body = append(body, fmt.Sprintf("null,null,null,null,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v",
		subs.Successful["10"], subs.Successful["10"]+subs.Redundant["10"], subs.Failed["10"], subs.Skipped["10"],
		subs.Successful["15"], subs.Successful["15"]+subs.Redundant["15"], subs.Failed["15"], subs.Skipped["15"],
		subs.Successful["30"], subs.Successful["30"]+subs.Redundant["30"], subs.Failed["30"], subs.Skipped["30"]))

	body = append(body, ",,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,")

	var times string
	for _, t := range subs.EvalTimes.Failed {
		times += strconv.Itoa(t) + ","
	}
	for _, t := range subs.EvalTimes.Successful {
		times += strconv.Itoa(t) + ","
	}
	for _, t := range subs.EvalTimes.Redundant {
		times += strconv.Itoa(t) + ","
	}
	for _, t := range subs.EvalTimes.Skipped {
		times += strconv.Itoa(t) + ","
	}

	body = append(body, "strategy,objective,experiment,iteration,times[ms]")
	body = append(body, fmt.Sprintf("null,null,null,null,%v", times))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Error("[API] Error encoding.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
