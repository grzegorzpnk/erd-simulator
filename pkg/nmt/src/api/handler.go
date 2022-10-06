package api

import (
	log "10.254.188.33/matyspi5/erd/pkg/nmt/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/mec-topology"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/metrics"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
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
			json.NewEncoder(w).Encode(h.graphClient.MecHosts[i])
		}
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

func (h *apiHandler) createEdgeHandler(w http.ResponseWriter, r *http.Request) {
	//todo: validate body of REST POST
	w.Header().Set("Content-Type", "application/json")

	var edge model.Edge
	_ = json.NewDecoder(r.Body).Decode(&edge)
	log.Infof("Client tries to add new Edge: %v --- %v \n", edge.SourceVertexName, edge.TargetVertexName)
	h.graphClient.AddEdge(edge)
}

func (h *apiHandler) getEdgesHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	for i := range h.graphClient.Edges {
		json.NewEncoder(w).Encode(h.graphClient.Edges[i])
	}

}

func (h *apiHandler) getAllMecHostsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	for i, _ := range h.graphClient.MecHosts {
		json.NewEncoder(w).Encode(h.graphClient.MecHosts[i].Identity)
	}

}

func (h *apiHandler) updateClusterCPUResources(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	cluster, _ := params["cluster"]
	provider, _ := params["provider"]

	mecHost := model.MecHost{}
	mecHost.Identity.Cluster = cluster
	mecHost.Identity.Provider = provider

	var clusterMetrics metrics.ClusterResources
	_ = json.NewDecoder(r.Body).Decode(&clusterMetrics)

	if h.graphClient.CheckGraphContainsVertex(mecHost) {
		h.graphClient.GetMecHost(cluster, provider).CpuResources.UpdateClusterMetrics(clusterMetrics)
		w.WriteHeader(http.StatusOK)
		log.Infof("Client updates cluster metrics for vertex ID: %v\n", cluster)
	} else {
		err := fmt.Errorf("Vertex %v not updated beacuse it's not exist", cluster)
		log.Errorf(err.Error())
		w.WriteHeader(http.StatusConflict)
	}
}

func (h *apiHandler) getClusterCPUResources(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	cluster, _ := params["cluster"]
	provider, _ := params["provider"]
	mecHost := model.MecHost{}
	mecHost.Identity.Cluster = cluster
	mecHost.Identity.Provider = provider

	if h.graphClient.CheckGraphContainsVertex(mecHost) {
		json.NewEncoder(w).Encode(h.graphClient.GetMecHost(cluster, provider).CpuResources)
		w.WriteHeader(http.StatusOK)
	} else {
		err := fmt.Errorf("Vertex %v not not exist", cluster)
		log.Errorf(err.Error())
		w.WriteHeader(http.StatusConflict)
	}
}

func updateEdgeMetrics(w http.ResponseWriter, r *http.Request) {

	/*	w.Header().Set("Content-Type", "application/json")

		params := mux.Vars(r)
		idSource, _ := strconv.Atoi(params["IdSource"])
		idTarget, _ := strconv.Atoi(params["IdTarget"])
		var edgeMetrics NetworkMetrics
		_ = json.NewDecoder(r.Body).Decode(&edgeMetrics)
	*/

	//sprawdz czy istnieje dany link i go pobierz
	//update danych na łaczu
	/*if exist(graph.MecHosts, id) {
		graph.getVertex(id).VertexMetrics.updateClusterResourcesMetrics(clusterMetrics)
		w.WriteHeader(http.StatusOK)
		log.Infof("Client updates cluster metrics for vertex ID: %v\n", params["cluster"])
	} else {
		err := fmt.Errorf("Vertex %v not updated beacuse it's not exist", id)
		log.Errorf(err.Error())
		w.WriteHeader(http.StatusConflict)
	}*/
}
