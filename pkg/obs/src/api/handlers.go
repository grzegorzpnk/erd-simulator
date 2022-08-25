package api

import (
	log "10.254.188.33/matyspi5/erd/pkg/obs/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/obs/src/pkg/latency"
	"10.254.188.33/matyspi5/erd/pkg/obs/src/pkg/observability"
	"fmt"

	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type apiHandler struct {
	obsClient observability.ClustersInfo
	ltcClient latency.MockClient
}

func (h *apiHandler) SetClients(ksmClient observability.ClustersInfo, ltcClient latency.MockClient) {
	h.obsClient = ksmClient
	h.ltcClient = ltcClient
}

func (h *apiHandler) getCpuReqHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	cluster := vars["cluster"]

	value, err := h.obsClient.GetClusterCpuReq(provider, cluster)
	if err != nil {
		fmt.Errorf("[API] Error: %v", err)
	}

	if value != -1 {
		sendResponse(w, value, http.StatusOK)
	} else {
		sendResponse(w, err.Error(), http.StatusNoContent)
	}

}

func (h *apiHandler) getCpuLimHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	cluster := vars["cluster"]

	value, err := h.obsClient.GetClusterCpuLim(provider, cluster)
	if err != nil {
		fmt.Errorf("[API] Error: %v", err)
	}

	if value != -1 {
		sendResponse(w, value, http.StatusOK)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *apiHandler) getMemReqHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	cluster := vars["cluster"]

	value, err := h.obsClient.GetClusterMemReq(provider, cluster)
	if err != nil {
		fmt.Errorf("[API] Error: %v", err)
	}

	if value != -1 {
		sendResponse(w, value, http.StatusOK)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *apiHandler) getMemLimHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	cluster := vars["cluster"]

	value, err := h.obsClient.GetClusterMemLim(provider, cluster)
	if err != nil {
		fmt.Errorf("[API] Error: %v", err)
	}

	if value != -1 {
		sendResponse(w, value, http.StatusOK)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *apiHandler) getLatencyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cell := vars["cell-id"]
	meh := vars["meh-id"]

	value, err := h.ltcClient.GetMockedLatencyMs(cell, meh)
	if err != nil {
		fmt.Errorf("[API] Error: %v", err)
	}

	if value != -1 {
		sendResponse(w, value, http.StatusOK)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// sendResponse sends an HTTP response to the client with the provided status
func sendResponse(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Error("Error encoding response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
