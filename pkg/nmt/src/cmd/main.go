package main

import (
	"log"
	"net/http"
	"nmt/src/api"
	"nmt/src/config"
	"nmt/src/package/topology"
)

var graph *topology.Graph

func main() {

	graph = &topology.Graph{}
	initializingGraph()

	//go topology.TopologyMetricsUpdate()

	httpRouter := api.NewRouter(graph)

	httpServer := &http.Server{
		Handler: httpRouter,
		Addr:    ":" + config.GetConfiguration().ServicePort,
	}

	log.Fatal(httpServer.ListenAndServe())

}

func initializingGraph() {

	for i := 0; i < 6; i++ {
		if i%2 == 0 {
			graph.AddVertex(topology.Vertex{Id: i, Type: "MEC", VertexMetrics: topology.ClusterMetrics{20, 50, 80}})
		} else {
			graph.AddVertex(topology.Vertex{Id: i, Type: "CELL"})
		}
	}

	graph.AddEdge(topology.Edge{1, 4, topology.NetworkMetrics{1.3, 10}})
	graph.AddEdge(topology.Edge{2, 5, topology.NetworkMetrics{1.3, 10}})
	graph.AddEdge(topology.Edge{3, 2, topology.NetworkMetrics{1.3, 10}})
	graph.AddEdge(topology.Edge{1, 0, topology.NetworkMetrics{1.3, 10}})
	graph.AddEdge(topology.Edge{4, 5, topology.NetworkMetrics{1.3, 10}})
	graph.PrintGraph()
}
