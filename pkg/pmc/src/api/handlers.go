package api

import (
	log "10.254.188.33/matyspi5/pmc/src/logger"
	"10.254.188.33/matyspi5/pmc/src/pkg/observability"

	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type apiHandler struct {
	client observability.ClustersInfo
}

func (h *apiHandler) SetKsmClient(client observability.ClustersInfo) {
	h.client = client
}

func (h *apiHandler) getCpuReqHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	cluster := vars["cluster"]

	value := h.client.GetClusterCpuReq(provider, cluster)

	if value > 0 {
		sendResponse(w, value, http.StatusOK)
	} else {
		sendResponse(w, value, http.StatusNoContent)
	}

}

func (h *apiHandler) getCpuLimHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	cluster := vars["cluster"]

	value := h.client.GetClusterCpuLim(provider, cluster)

	if value > 0 {
		sendResponse(w, value, http.StatusOK)
	} else {
		sendResponse(w, value, http.StatusNoContent)
	}
}

func (h *apiHandler) getMemReqHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	cluster := vars["cluster"]

	value := h.client.GetClusterMemReq(provider, cluster)

	if value > 0 {
		sendResponse(w, value, http.StatusOK)
	} else {
		sendResponse(w, value, http.StatusNoContent)
	}
}

func (h *apiHandler) getMemLimHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	cluster := vars["cluster"]

	value := h.client.GetClusterMemLim(provider, cluster)

	if value > 0 {
		sendResponse(w, value, http.StatusOK)
	} else {
		sendResponse(w, value, http.StatusNoContent)
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
