package api

import (
	log "10.254.188.33/matyspi5/erd/pkg/obs/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/obs/src/pkg/latency"
	"10.254.188.33/matyspi5/erd/pkg/obs/src/pkg/model"
	"10.254.188.33/matyspi5/erd/pkg/obs/src/pkg/observability"
	"10.254.188.33/matyspi5/erd/pkg/obs/src/pkg/promql"
	"errors"
	"reflect"
	"strconv"
	"strings"

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

func (h *apiHandler) getMemInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	cluster := vars["cluster"]

	alloc, err := h.obsClient.GetClusterAlloc(promql.MEMORY, provider, cluster)
	if err != nil {
		log.Errorf("[API] Error: %v", err)
	}

	req, err := h.obsClient.GetClusterReq(promql.MEMORY, provider, cluster)
	if err != nil {
		log.Errorf("[API] Error: %v", err)
	}

	util := 100 * (req / alloc)

	val := model.MecResInfo{
		Used:        req,
		Allocatable: alloc,
		Utilization: util,
	}

	if alloc != -1 && req != -1 {
		sendResponse(w, val, http.StatusOK)
	} else {
		sendResponse(w, err.Error(), http.StatusNoContent)
	}
}

func (h *apiHandler) getCpuInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	cluster := vars["cluster"]

	alloc, err := h.obsClient.GetClusterAlloc(promql.CPU, provider, cluster)
	if err != nil {
		log.Errorf("[API] Error: %v", err)
	}

	req, err := h.obsClient.GetClusterReq(promql.CPU, provider, cluster)
	if err != nil {
		log.Errorf("[API] Error: %v", err)
	}

	util := 100 * (req / alloc)

	val := model.MecResInfo{
		Used:        req,
		Allocatable: alloc,
		Utilization: util,
	}

	body := val

	if alloc != -1 && req != -1 {
		sendResponse(w, body, http.StatusOK)
	} else {
		sendResponse(w, err.Error(), http.StatusNoContent)
	}
}

func (h *apiHandler) getCpuRequestsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	cluster := vars["cluster"]

	val, err := h.obsClient.GetClusterReq(promql.CPU, provider, cluster)
	if err != nil {
		log.Errorf("[API] Error: %v", err)
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
		log.Errorf("[API] Error: %v", err)
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
		log.Errorf("[API] Error: %v", err)
	}

	if val != -1 && val != 0 {
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
		log.Errorf("[API] Error: %v", err)
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
		log.Errorf("[API] Error: %v", err)
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
		log.Errorf("[API] Error: %v", err)
	}

	if val != -1 && val != 0 {
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
		log.Errorf("[API] Error: %v", err)
	}

	alloc, err := h.obsClient.GetClusterAlloc(promql.CPU, provider, cluster)
	if err != nil {
		log.Errorf("[API] Error: %v", err)
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
		log.Errorf("[API] Error: %v", err)
	}

	alloc, err := h.obsClient.GetClusterAlloc(promql.MEMORY, provider, cluster)
	if err != nil {
		log.Errorf("[API] Error: %v", err)
	}

	val := 100 * (req / alloc)

	if req != -1 && alloc != -1 {
		sendResponse(w, val, http.StatusOK)
	} else {
		sendResponse(w, err.Error(), http.StatusNoContent)
	}
}

func (h *apiHandler) getLatencyHandler(w http.ResponseWriter, r *http.Request) {
	var target, source interface{}
	var err error

	vars := mux.Vars(r)
	sourceId := vars["source-node"]
	targetId := vars["target-node"]

	source, err = h.checkType(sourceId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	target, err = h.checkType(targetId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if reflect.TypeOf(source) == reflect.TypeOf(target) && reflect.TypeOf(source) == reflect.TypeOf(model.Cell{}) {
		err := errors.New(fmt.Sprintf("source and target node are both of %T type, which is not permited", target))
		log.Warnf("couldn't get latency: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	value, err := h.ltcClient.GetMockedLatencyMs(source, target)
	if err != nil {
		log.Errorf("[API] Error: %v", err)
	}

	if value != -1 {
		sendResponse(w, value, http.StatusOK)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *apiHandler) checkType(param string) (object interface{}, err error) {
	// CellID suppose to be an Integer value
	_, err = strconv.Atoi(param)
	if err != nil {
		// MecHost suppose to be in format provider+cluster
		ok := strings.Contains(param, "+")
		if !ok {
			err = errors.New(fmt.Sprintf("couldn't find cell or mecHost based on sourceId[%v]", param))
			return
		} else {
			object, err = h.ltcClient.GetMehByFqdn(param)
			if err != nil {
				return
			}
		}
	} else {
		object, err = h.ltcClient.GetCellById(param)
		if err != nil {
			return
		}
	}
	return
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
