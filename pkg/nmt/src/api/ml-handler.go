package api

import (
	log "10.254.188.33/matyspi5/erd/pkg/nmt/src/logger"
	mec_topology "10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/mec-topology"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type MecNode struct {
	ID                string    `json:"id"`
	CPUCapacity       int       `json:"cpu_capacity"`
	MemoryCapacity    int       `json:"memory_capacity"`
	CPUUtilization    int       `json:"cpu_utilization"`
	MemoryUtilization int       `json:"memory_utilization"`
	LatencyVector     []float64 `json:"latency_array"`
	PlacementCost     float64   `json:"placement_cost"`
}

func (h *apiHandler) MLInitialState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	appQuantity, _ := strconv.Atoi(params["applications"])

	var ac mec_topology.AppCounter
	ac.V2x = int(appQuantity / 3)
	ac.Cg = int(appQuantity / 3)
	ac.Uav = int(appQuantity / 3)

	h.graphClient.DeleteAllDeclaredApps()
	h.graphClient.DeclareApplications(ac)

	//find candidates mec and assign
	status, mecHosts := h.graphClient.FindInitialClusters()
	if status == true {
		fmt.Printf("Apps with clusters:\n")
		for i := 0; i < len(h.graphClient.Application); i++ {
			h.graphClient.Application[i].PrintApplication()
		}
	} else {
		w.WriteHeader(http.StatusConflict)
		log.Errorf("Cannot find clusters for number declared of apps. Probably due to overloading\n")
		return
	}

	var response []MecNode
	var mec_node MecNode
	for _, v := range mecHosts {
		mec_node.ID = v.Identity.Cluster
		mec_node.CPUCapacity = int(v.CpuResources.Capacity)
		mec_node.MemoryCapacity = int(v.MemoryResources.Capacity)
		mec_node.CPUUtilization = int(v.CpuResources.Utilization)
		mec_node.MemoryUtilization = int(v.MemoryResources.Utilization)
		switch v.Identity.Location.Level {
		case 0:
			mec_node.PlacementCost = 1.0
		case 1:
			mec_node.PlacementCost = 0.66667
		case 2:
			mec_node.PlacementCost = 0.33333
		}
		//todo: mec_node.LatencyVector =
		response = append(response, mec_node)
	}

	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)

}

func (h *apiHandler) MLInitialConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var mecHosts []model.MecHost

	for _, v := range h.graphClient.MecHosts {
		mecHosts = append(mecHosts, *v)
	}

	var response []MecNode
	var mec_node MecNode
	for _, v := range mecHosts {
		mec_node.ID = v.Identity.Cluster
		mec_node.CPUCapacity = int(v.CpuResources.Capacity)
		mec_node.MemoryCapacity = int(v.MemoryResources.Capacity)
		mec_node.CPUUtilization = 0
		mec_node.MemoryUtilization = 0
		switch v.Identity.Location.Level {
		case 0:
			mec_node.PlacementCost = 1.0
		case 1:
			mec_node.PlacementCost = 0.66667
		case 2:
			mec_node.PlacementCost = 0.33333
		}

		for i := 1; i <= len(h.graphClient.NetworkCells); i++ {
			cell := h.graphClient.GetCell(strconv.Itoa(i))
			latency, _ := h.graphClient.ShortestPath(cell, &v)
			mec_node.LatencyVector = append(mec_node.LatencyVector, latency)
		}

		response = append(response, mec_node)
		mec_node.LatencyVector = nil
	}

	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)

}

func (h *apiHandler) MLGetRANs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//todo: consider max-cell-numbers
	log.Infof("number of RANS: %v", len(h.graphClient.NetworkCells))
	json.NewEncoder(w).Encode(len(h.graphClient.NetworkCells))
	w.WriteHeader(http.StatusOK)

}