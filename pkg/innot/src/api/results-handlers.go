package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"

	log "10.254.188.33/matyspi5/erd/pkg/innot/src/logger"
)

// Results collection handlers

func (h *apiHandler) relocationFailedHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	h.resClient.Results.IncFailed(params["type"])

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func (h *apiHandler) relocationSuccessfulHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	h.resClient.Results.IncSuccessful(params["type"])

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func (h *apiHandler) relocationRedundantHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	h.resClient.Results.IncRedundant(params["type"])

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func (h *apiHandler) relocationSkippedHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	h.resClient.Results.IncSkipped(params["type"])

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

// Save computation times

func (h *apiHandler) failedTimeHandler(w http.ResponseWriter, r *http.Request) {
	var body time.Duration

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Infof("Got new time[%v] for failed relocation", body)
	h.resClient.Results.AddFailedTime(int(body.Milliseconds()))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func (h *apiHandler) successfulTimeHandler(w http.ResponseWriter, r *http.Request) {
	var body time.Duration

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Infof("Got new time[%v] for successful relocation", body)
	h.resClient.Results.AddSuccessfulTime(int(body.Milliseconds()))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func (h *apiHandler) redundantTimeHandler(w http.ResponseWriter, r *http.Request) {
	var body time.Duration

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Infof("Got new time[%v] for redundant relocation", body)
	h.resClient.Results.AddRedundantTime(int(body.Milliseconds()))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func (h *apiHandler) skippedTimeHandler(w http.ResponseWriter, r *http.Request) {
	var body time.Duration

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Infof("Got new time[%v] for skipped relocation", body)
	h.resClient.Results.AddSkippedTime(int(body.Milliseconds()))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func (h *apiHandler) resetHandler(w http.ResponseWriter, r *http.Request) {

	h.resClient.Results.Reset()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func (h *apiHandler) getResultsHandler(w http.ResponseWriter, r *http.Request) {

	subs := h.resClient.Results

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(subs)
	if err != nil {
		log.Error("[API] Error encoding.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *apiHandler) getResultsCSVHandler(w http.ResponseWriter, r *http.Request) {
	var body []string
	subs := h.resClient.Results

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
