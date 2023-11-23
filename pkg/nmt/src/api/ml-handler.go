package api

import (
	log "10.254.188.33/matyspi5/erd/pkg/nmt/src/logger"
	mec_topology "10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/mec-topology"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
	"strconv"
	"strings"
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

		if v.Identity.Location.Level == 0 {
			mec_node.CPUUtilization = 1552 * 100 / mec_node.CPUCapacity
			mec_node.MemoryUtilization = 1112 * 100 / mec_node.MemoryCapacity
		}
		if v.Identity.Location.Level == 1 {
			mec_node.CPUUtilization = 1200 * 100 / mec_node.CPUCapacity
			mec_node.MemoryUtilization = 1080 * 100 / mec_node.MemoryCapacity
		}
		if v.Identity.Location.Level == 2 {
			mec_node.CPUUtilization = 1548 * 100 / mec_node.CPUCapacity
			mec_node.MemoryUtilization = 1080 * 100 / mec_node.MemoryCapacity
		}

		switch v.Identity.Location.Level {
		case 0:
			mec_node.PlacementCost = 1.0
		case 1:
			mec_node.PlacementCost = 0.66667
		case 2:
			mec_node.PlacementCost = 0.33334
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

func (h *apiHandler) GetCurrentState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	intent := SmartPlacementIntent{}
	err := json.NewDecoder(r.Body).Decode(&intent)
	if err != nil {
		log.Errorf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := [][]int{}

	if intent.Spec.SmartPlacementIntentData.Masked != "masked" {
		fmt.Printf("Received request from simu to preapre current state of MECs for input to ML NONMASKED client")
		cellINT, _ := strconv.Atoi(string(intent.Spec.SmartPlacementIntentData.TargetCell))

		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		//let's work on a sorted copy
		var mecHostsCopy []model.MecHost

		for _, mec := range h.graphClient.MecHosts {
			mecHostsCopy = append(mecHostsCopy, *mec)
		}

		sort.Slice(mecHostsCopy, func(i, j int) bool {
			id1 := mecHostsCopy[i].Identity.Cluster
			id2 := mecHostsCopy[j].Identity.Cluster

			// Extract the number from the id attribute
			num1, _ := strconv.Atoi(strings.TrimPrefix(id1, "mec"))
			num2, _ := strconv.Atoi(strings.TrimPrefix(id2, "mec"))

			// Compare the numbers extracted from the id attribute
			return num1 < num2
		})

		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		response = make([][]int, len(h.graphClient.MecHosts))
		for i := 0; i < len(response); i++ {
			response[i] = make([]int, 13)
		}

		// MEC(for MEC each)    : 1) CPU Capacity 2) CPU Utilization [%] 3) Memory Capacity 4) Memory Utilization [%] 5) Unit Cost
		for i, v := range mecHostsCopy {
			response[i][0] = int(v.GetCpuCapacity() - v.GetCpuUsed())
			response[i][1] = int(v.GetCpuCapacity())
			response[i][2] = int(intent.Spec.SmartPlacementIntentData.AppCpuReq)
			response[i][4] = int(v.GetMemoryCapacity() - v.GetMemoryUsed())
			response[i][5] = int(v.GetMemoryCapacity())
			response[i][6] = int(intent.Spec.SmartPlacementIntentData.AppMemReq)
			response[i][7] = (int(v.LatencyVector[cellINT-1]) + 1) * 100
			response[i][8] = int(intent.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax * 100)
			response[i][9] = (int(v.LatencyVector[cellINT-1]) + 1) * 100
			response[i][10] = int(intent.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax * 100)
			response[i][11] = (int(v.LatencyVector[cellINT-1]) + 1) * 100
			response[i][12] = int(intent.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax * 100)
			if v.Identity.Cluster == intent.CurrentPlacement.Cluster {
				response[i][3] = 1
			} else {
				response[i][3] = 0
			}
		}
	} else {
		fmt.Printf("Received request from simu to preapre current state of MECs for input to ML MASKED client")
		response = make([][]int, len(h.graphClient.MecHosts))
		for i := 0; i < len(response); i++ {
			response[i] = make([]int, 5)
		}

		// MEC(for MEC each)    : 1) CPU Capacity 2) CPU Utilization [%] 3) Memory Capacity 4) Memory Utilization [%] 5) Unit Cost
		for i, v := range h.graphClient.MecHosts {
			response[i][0] = determineStateOfCapacity(int(v.GetCpuCapacity()))
			response[i][1] = int(v.GetCpuUtilization() * 100)
			response[i][2] = determineStateOfCapacity(int(v.GetMemoryCapacity()))
			response[i][3] = int(v.GetMemoryUtilization() * 100)
			response[i][4] = determineStateofCost(int(v.Identity.Location.Level))
		}
	}

	fmt.Printf("State of MECs before returning to Simu:\n ")
	printSlice(response)

	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)

	//added, cause due to training of non-masked it turned out that changes of obserability space was needed
	//if cell != "masked" {
	//	cellINT, _ := strconv.Atoi(string(cell))
	//	fmt.Printf("Received request from simu to preapre current state of MECs for input to ML client")
	//	response = make([][]int, len(h.graphClient.MecHosts))
	//	for i := 0; i < len(response); i++ {
	//		response[i] = make([]int, 6)
	//	}
	//
	//	// MEC(for MEC each)    : 1) CPU Capacity 2) CPU Utilization [%] 3) Memory Capacity 4) Memory Utilization [%] 5) Unit Cost
	//	for i, v := range h.graphClient.MecHosts {
	//		response[i][0] = determineStateOfCapacity(int(v.GetCpuCapacity()))
	//		response[i][1] = int(v.GetCpuUtilization() * 100)
	//		response[i][2] = determineStateOfCapacity(int(v.GetMemoryCapacity()))
	//		response[i][3] = int(v.GetMemoryUtilization() * 100)
	//		response[i][4] = determineStateofCost(int(v.Identity.Location.Level))
	//		response[i][5] = int(v.LatencyVector[cellINT-1] + 1)
	//	}
	//}

}

func (h *apiHandler) GetCurrentMask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	app := model.MECApp{}
	//responseMask := make([]int, len(h.graphClient.MecHosts))
	var responseMask []int

	err := json.NewDecoder(r.Body).Decode(&app)
	if err != nil || app.Id == "" {
		log.Errorf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//let's work on a copy
	var mecHostsCopy []model.MecHost

	for _, mec := range h.graphClient.MecHosts {
		mecHostsCopy = append(mecHostsCopy, *mec)
	}

	//now lets sort this slice of MEC hosts in order by ID
	sort.Slice(mecHostsCopy, func(i, j int) bool {
		id1 := mecHostsCopy[i].Identity.Cluster
		id2 := mecHostsCopy[j].Identity.Cluster

		// Extract the number from the id attribute
		num1, _ := strconv.Atoi(strings.TrimPrefix(id1, "mec"))
		num2, _ := strconv.Atoi(strings.TrimPrefix(id2, "mec"))

		// Compare the numbers extracted from the id attribute
		return num1 < num2
	})

	for _, mec := range mecHostsCopy {
		if mec.CheckEnoughResources(app) && mec.CheckLatency(app) {
			responseMask = append(responseMask, 1)
		} else {
			responseMask = append(responseMask, 0)
		}
	}

	json.NewEncoder(w).Encode(responseMask)
	w.WriteHeader(http.StatusOK)

}

func printSlice(slice [][]int) {
	for i := 0; i < len(slice); i++ {
		for j := 0; j < len(slice[i]); j++ {
			fmt.Print(slice[i][j], " ")
		}
		fmt.Println()
	}
}

func determineStateOfCapacity(capacityValue int) int {
	capacityMap := map[int]int{4000: 1, 8000: 2, 12000: 3}
	return capacityMap[capacityValue]
}

//jeśli level 0 (city-level), to koszt placementu jest najwyzszy i jest reprezentowany przez wartość 3 itd..
func determineStateofCost(placementCost int) int {
	costMap := map[int]int{0: 3, 1: 2, 2: 1}
	return costMap[placementCost]
}

type CellId string

// CurrentPlacement represents where the application is currently instantiated
type CurrentPlacement struct {
	Provider        string `json:"provider"`
	Cluster         string `json:"cluster"`
	ClusterCapacity int    `json:"omitempty,clusterCapacity"`
}

// SmartPlacementIntent defines the Intent to perform Smart Placement for giver application
type SmartPlacementIntent struct {
	Metadata         Metadata                 `json:"metadata"`
	CurrentPlacement CurrentPlacement         `json:"currentPlacement"`
	Spec             SmartPlacementIntentSpec `json:"spec"`
}

type Metadata struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description,omitempty" yaml:"-"`
	UserData1   string `json:"userData1,omitempty" yaml:"-"`
	UserData2   string `json:"userData2,omitempty" yaml:"-"`
}

type SmartPlacementIntentSpec struct {
	AppName                  string                     `json:"app"`
	SmartPlacementIntentData SmartPlacementIntentStruct `json:"smartPlacementIntentData"`
}

type SmartPlacementIntentStruct struct {
	TargetCell        CellId      `json:"targetCell"`
	AppCpuReq         float64     `json:"appCpuReq"`
	AppMemReq         float64     `json:"appMemReq"`
	ConstraintsList   Constraints `json:"constraintsList"`
	ParametersWeights Weights     `json:"parametersWeights,omitempty"`
	Masked            string      `json:"masked,omitempty"`
}

type Constraints struct {
	LatencyMax float64 `json:"latencyMax"`
}

type Weights struct {
	LatencyWeight        float64 `json:"latencyWeight"`
	ResourcesWeight      float64 `json:"resourcesWeight"`
	CpuUtilizationWeight float64 `json:"cpuUtilizationWeight"`
	MemUtilizationWeight float64 `json:"memUtilizationWeight"`
}

type SmartPlacementIntentKey struct {
	Project               string `json:"project"`
	CompositeApp          string `json:"compositeApp"`
	CompositeAppVersion   string `json:"compositeAppVersion"`
	DeploymentIntentGroup string `json:"deploymentIntentGroup"`
	SmartPlacementIntent  string `json:"smartPlacementIntent"`
}
