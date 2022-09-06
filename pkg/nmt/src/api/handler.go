package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"nmt/src/package/topology"
)

type apiHandler struct {
	graphClient topology.Graph
}

func (h *apiHandler) SetClients(graphClient topology.Graph) {
	h.graphClient = graphClient
}

func (h *apiHandler) createEdgeHandler(w http.ResponseWriter, r *http.Request) {
	//todo: validate body of REST POST
	w.Header().Set("Content-Type", "application/json")

	var edge topology.Edge
	_ = json.NewDecoder(r.Body).Decode(&edge)
	fmt.Printf("Client tries to add new Edge: %v --- %v \n", edge.Source, edge.Target)
	h.graphClient.AddEdge(edge)

}

func (h *apiHandler) updateClusterMetrics(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, _ := params["Id"]
	var clusterMetrics topology.ClusterMetrics
	_ = json.NewDecoder(r.Body).Decode(&clusterMetrics)

	if topology.ContainsVertex(h.graphClient.Vertices, topology.Vertex{Id: id, Type: "MEC"}) {
		h.graphClient.GetVertex(id).VertexMetrics.UpdateClusterMetrics(clusterMetrics)
		w.WriteHeader(http.StatusOK)
		fmt.Printf("Client updates cluster metrics for vertex ID: %v\n", params["Id"])
	} else {
		err := fmt.Errorf("Vertex %v not updated beacuse it's not exist", id)
		fmt.Println(err.Error())
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
	/*if exist(graph.Vertices, id) {
		graph.getVertex(id).VertexMetrics.updateClusterMetrics(clusterMetrics)
		w.WriteHeader(http.StatusOK)
		fmt.Printf("Client updates cluster metrics for vertex ID: %v\n", params["Id"])
	} else {
		err := fmt.Errorf("Vertex %v not updated beacuse it's not exist", id)
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusConflict)
	}*/
}

func (h *apiHandler) getVertexHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for i, v := range h.graphClient.Vertices {
		if v.Id == params["Id"] {
			json.NewEncoder(w).Encode(h.graphClient.Vertices[i])
		}
	}

}

func (h *apiHandler) createVertex(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var vertex topology.Vertex
	_ = json.NewDecoder(r.Body).Decode(&vertex)
	fmt.Printf("Client tries to add new vertex ID: %v\n", vertex.Id)
	if topology.ContainsVertex(h.graphClient.Vertices, vertex) {
		err := fmt.Errorf("Vertex %v not added beacuse it is an existing key", vertex.Id)
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusConflict)
	} else {
		if containsAnyEdge(vertex) {
			vertex.Neighbours = nil
		}
		h.graphClient.AddVertex(vertex)
		w.WriteHeader(http.StatusOK)
	}
}

func containsAnyEdge(vertex topology.Vertex) bool {

	if vertex.Neighbours != nil {
		return true
	} else {
		return false
	}

}

func (h *apiHandler) getEdgesHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	for i, _ := range h.graphClient.Edges {
		json.NewEncoder(w).Encode(h.graphClient.Edges[i])
	}

}

func (h *apiHandler) getAllVertexesHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.graphClient.Vertices)

}
