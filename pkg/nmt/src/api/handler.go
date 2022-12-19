package api

import (
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/djikstra"
	log "10.254.188.33/matyspi5/erd/pkg/nmt/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/mec-topology"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

//prereqquesties types and function
type apiHandler struct {
	graphClient mec_topology.Graph
}

func (h *apiHandler) SetClients(graphClient mec_topology.Graph) {
	h.graphClient = graphClient
}

//main functions

func (h *apiHandler) createMecHostHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var mecHost model.MecHost
	_ = json.NewDecoder(r.Body).Decode(&mecHost)
	log.Infof("Client tries to add new mecHost ID: %v, provider: %v\n", mecHost.Identity.Cluster, mecHost.Identity.Provider)
	if h.graphClient.CheckGraphContainsVertex(mecHost) {
		err := fmt.Errorf("Mec Host %v, %v not added beacuse it already exists", mecHost.Identity.Cluster, mecHost.Identity.Provider)
		log.Errorf(err.Error())
		w.WriteHeader(http.StatusConflict)
	} else {
		//this methods is only for create Mec Host, this is blocker for not to create list of edges at that time
		if containsAnyEdge(mecHost) {
			mecHost.Neighbours = nil
		}
		h.graphClient.AddMecHost(mecHost)
		w.WriteHeader(http.StatusOK)
	}
}

func (h *apiHandler) getMecHostHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for i, v := range h.graphClient.MecHosts {
		if v.Identity.Cluster == params["cluster"] &&
			v.Identity.Provider == params["provider"] {
			json.NewEncoder(w).Encode(h.graphClient.MecHosts[i].Identity)
			break
		}
	}
}

func (h *apiHandler) getCellAssociatedMecHostsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	cellId, _ := params["cell-id"]
	var zone, region string
	response := make([]model.MecIdentity, 0)

	//send N level
	for i, v := range h.graphClient.MecHosts {
		if v.CheckMECsupportsCell(cellId) {
			response = append(response, h.graphClient.MecHosts[i].Identity)
			//json.NewEncoder(w).Encode(h.graphClient.MecHosts[i].Identity)
			zone = v.Identity.Location.Zone
			region = v.Identity.Location.Region
		}
	}

	//send N+1 level
	for i, v := range h.graphClient.MecHosts {
		if v.Identity.Location.Level == 1 {
			if v.Identity.Location.Zone == zone {
				//json.NewEncoder(w).Encode(h.graphClient.MecHosts[i].Identity)
				response = append(response, h.graphClient.MecHosts[i].Identity)
			}
		}
	}

	//send N+2 level
	for i, v := range h.graphClient.MecHosts {
		if v.Identity.Location.Level == 2 {
			if v.Identity.Location.Region == region {
				//json.NewEncoder(w).Encode(h.graphClient.MecHosts[i].Identity)
				response = append(response, h.graphClient.MecHosts[i].Identity)
			}
		}
	}
	json.NewEncoder(w).Encode(response)
}

func (h *apiHandler) shortestPathHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	cell, _ := params["cell-id"]
	cProvider, _ := params["provider"]
	cName, _ := params["cluster"]

	destCluster := h.graphClient.GetMecHost(cName, cProvider)
	startCell := h.graphClient.GetCell(cell)
	if destCluster == nil {
		log.Fatalln("destination MEC host not recognized!")

	}
	//check if they are direct neighbours, if so the latency is just between start and stop node
	if destCluster.CheckMECsupportsCell(startCell.Id) {
		json.NewEncoder(w).Encode(destCluster.GetCell(cell).Latency)
		fmt.Printf("direct nodes, latency between cell: %v and mec: [%v+%v], is: %v", startCell.Id, cProvider, cName, destCluster.GetCell(cell).Latency)

	} else {
		// if not, we have to calculate path between all MEC clusters that are in the same local zone as cell, to the target cluster, the final latency is a sum of the calculated one + between started mec and cell
		var startClusters []model.MecHost

		for _, v := range h.graphClient.MecHosts {
			if v.Identity.Location.LocalZone == startCell.LocalZone {
				startClusters = append(startClusters, *v)
			}
		}

		var inputGraph djikstra.InputGraph
		inputGraph.Graph = make([]djikstra.InputData, 200)

		//add all mec hosts to temp graph todo: should be only subset of graph nodes
		for i, v := range h.graphClient.Edges {
			inputGraph.Graph[i].Source = v.SourceVertexName
			inputGraph.Graph[i].Destination = v.TargetVertexName
			inputGraph.Graph[i].Weight = v.EdgeMetrics.Latency
		}
		itemGraph := djikstra.CreateGraph(inputGraph)

		//calculate shortest path between all []startClusters and stopNode, where startClusters is a list of cluster directly associated with cell
		results := make([]ShortestPathResult, 0)

		for _, v := range startClusters {

			startNd := djikstra.Node{v.Identity.Cluster}
			stopNd := djikstra.Node{destCluster.Identity.Cluster}

			var resultTmp ShortestPathResult
			resultTmp.path, resultTmp.latencyResults = djikstra.GetShortestPath(&startNd, &stopNd, itemGraph)

			//add latency between cell and start MEC host
			resultTmp.latencyResults += h.graphClient.GetMecHost(v.Identity.Cluster, v.Identity.Provider).GetCell(cell).Latency

			results = append(results, resultTmp)

		}

		//find minimal value
		min := results[0].latencyResults
		for _, v := range results {
			if v.latencyResults < min {
				min = v.latencyResults
			}
		}

		for i, v := range results {

			if v.latencyResults == min {
				fmt.Printf("final path is: %v\n", results[i].path)
			}
		}

		json.NewEncoder(w).Encode(min)
		log.Infof("indirect nodes, latency between cell: %v and mec: [%v+%v], is: %v", startCell.Id, cProvider, cName, min)

	}

}

func (h *apiHandler) getMECCpu(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for i, v := range h.graphClient.MecHosts {
		if v.Identity.Cluster == params["cluster"] &&
			v.Identity.Provider == params["provider"] {
			json.NewEncoder(w).Encode(h.graphClient.MecHosts[i].CpuResources)
		}
	}
}

func (h *apiHandler) getMECMemory(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for i, v := range h.graphClient.MecHosts {
		if v.Identity.Cluster == params["cluster"] &&
			v.Identity.Provider == params["provider"] {
			json.NewEncoder(w).Encode(h.graphClient.MecHosts[i].MemoryResources)
		}
	}
}

func (h *apiHandler) getMECNeighbours(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for i, v := range h.graphClient.MecHosts {
		if v.Identity.Cluster == params["cluster"] &&
			v.Identity.Provider == params["provider"] {
			json.NewEncoder(w).Encode(h.graphClient.MecHosts[i].Neighbours)
		}
	}
}

func (h *apiHandler) createLinkHandler(w http.ResponseWriter, r *http.Request) {
	//todo: validate body of REST POST
	w.Header().Set("Content-Type", "application/json")

	var edge model.Edge
	_ = json.NewDecoder(r.Body).Decode(&edge)
	log.Infof("Client tries to add new Edge: %v --- %v \n", edge.SourceVertexName, edge.TargetVertexName)
	h.graphClient.AddLink(edge)
}

func (h *apiHandler) getEdgesHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	for i := range h.graphClient.Edges {
		json.NewEncoder(w).Encode(h.graphClient.Edges[i])
	}

}

func (h *apiHandler) getAllMecHostsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	response := make([]model.MecIdentity, 0)

	for i := range h.graphClient.MecHosts {
		response = append(response, h.graphClient.MecHosts[i].Identity)
	}
	json.NewEncoder(w).Encode(response)
}

// Handler for EXPERIMENTS purposes
func (h *apiHandler) getAllMecHostsWithMetricsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	response := make([]string, 0)
	realMecs := []string{"mec1", "mec3", "mec4", "mec5", "mec6", "mec7", "mec11", "mec12", "mec13", "mec14", "mec15", "mec16", "mec17", "mec18",
		"mec19", "mec20", "mec21", "mec22", "mec23", "mec24", "mec25", "mec26"}
	for i := range h.graphClient.MecHosts {
		mec := *h.graphClient.MecHosts[i]

		mecExists := false
		for _, mh := range realMecs {
			if mec.Identity.Cluster == mh {
				mecExists = true
				break
			}
		}
		if !mecExists {
			continue
		}

		mecDescription := "null,null,"
		mecDescription += mec.Identity.Cluster + ","
		mecDescription += strconv.FormatFloat(mec.CpuResources.Used, 'f', 1, 64) + ","
		mecDescription += strconv.FormatFloat(mec.CpuResources.Capacity, 'f', 1, 64) + ","
		mecDescription += strconv.FormatFloat(mec.CpuResources.Utilization, 'f', 1, 64) + ","
		mecDescription += strconv.FormatFloat(mec.MemoryResources.Used, 'f', 1, 64) + ","
		mecDescription += strconv.FormatFloat(mec.MemoryResources.Capacity, 'f', 1, 64) + ","
		mecDescription += strconv.FormatFloat(mec.MemoryResources.Utilization, 'f', 1, 64)

		response = append(response, mecDescription)
	}
	json.NewEncoder(w).Encode(response)
}
