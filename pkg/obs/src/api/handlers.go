package api

import (
	log "10.254.188.33/matyspi5/erd/pkg/obs/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/obs/src/pkg/latency"
	"10.254.188.33/matyspi5/erd/pkg/obs/src/pkg/observability"
	"10.254.188.33/matyspi5/erd/pkg/obs/src/pkg/promql"
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

func (h *apiHandler) getCpuRequestsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	cluster := vars["cluster"]

	val, err := h.obsClient.GetClusterReq(promql.CPU, provider, cluster)
	if err != nil {
		fmt.Errorf("[API] Error: %v", err)
	}

	if val != -1 {
		sendResponse(w, val, http.StatusOK)
	} else {
		sendResponse(w, err.Error(), http.StatusNoContent)
	}
}

func (h *apiHandler) getCpuLimitsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	cluster := vars["cluster"]

	val, err := h.obsClient.GetClusterLim(promql.CPU, provider, cluster)
	if err != nil {
		fmt.Errorf("[API] Error: %v", err)
	}

	if val != -1 {
		sendResponse(w, val, http.StatusOK)
	} else {
		sendResponse(w, err.Error(), http.StatusNoContent)
	}
}

func (h *apiHandler) getCpuAllocHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	cluster := vars["cluster"]

	val, err := h.obsClient.GetClusterAlloc(promql.CPU, provider, cluster)
	if err != nil {
		fmt.Errorf("[API] Error: %v", err)
	}

	if val != -1 {
		sendResponse(w, val, http.StatusOK)
	} else {
		sendResponse(w, err.Error(), http.StatusNoContent)
	}
}

func (h *apiHandler) getMemRequestsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	cluster := vars["cluster"]

	val, err := h.obsClient.GetClusterReq(promql.MEMORY, provider, cluster)
	if err != nil {
		fmt.Errorf("[API] Error: %v", err)
	}

	if val != -1 {
		sendResponse(w, val, http.StatusOK)
	} else {
		sendResponse(w, err.Error(), http.StatusNoContent)
	}
}

func (h *apiHandler) getMemLimitsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	cluster := vars["cluster"]

	val, err := h.obsClient.GetClusterLim(promql.MEMORY, provider, cluster)
	if err != nil {
		fmt.Errorf("[API] Error: %v", err)
	}

	if val != -1 {
		sendResponse(w, val, http.StatusOK)
	} else {
		sendResponse(w, err.Error(), http.StatusNoContent)
	}
}

func (h *apiHandler) getMemAllocHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	cluster := vars["cluster"]

	val, err := h.obsClient.GetClusterAlloc(promql.MEMORY, provider, cluster)
	if err != nil {
		fmt.Errorf("[API] Error: %v", err)
	}

	if val != -1 {
		sendResponse(w, val, http.StatusOK)
	} else {
		sendResponse(w, err.Error(), http.StatusNoContent)
	}
}

func (h *apiHandler) getCpuReqUtilizationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	cluster := vars["cluster"]

	req, err := h.obsClient.GetClusterReq(promql.CPU, provider, cluster)
	if err != nil {
		fmt.Errorf("[API] Error: %v", err)
	}

	alloc, err := h.obsClient.GetClusterAlloc(promql.CPU, provider, cluster)
	if err != nil {
		fmt.Errorf("[API] Error: %v", err)
	}

	val := 100 * (req / alloc)

	if req != -1 && alloc != -1 {
		sendResponse(w, val, http.StatusOK)
	} else {
		sendResponse(w, err.Error(), http.StatusNoContent)
	}
}

func (h *apiHandler) getMemReqUtilizationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	cluster := vars["cluster"]

	req, err := h.obsClient.GetClusterReq(promql.MEMORY, provider, cluster)
	if err != nil {
		fmt.Errorf("[API] Error: %v", err)
	}

	alloc, err := h.obsClient.GetClusterAlloc(promql.MEMORY, provider, cluster)
	if err != nil {
		fmt.Errorf("[API] Error: %v", err)
	}

	val := 100 * (req / alloc)

	if req != -1 && alloc != -1 {
		sendResponse(w, val, http.StatusOK)
	} else {
		sendResponse(w, err.Error(), http.StatusNoContent)
	}
}

func (h *apiHandler) getLatencyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cell := vars["cell-id"]
	meh := vars["mec-id"]

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
