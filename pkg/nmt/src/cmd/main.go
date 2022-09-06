package main

import (
	"log"
	"net/http"
	"nmt/src/api"
	"nmt/src/config"
	"nmt/src/package/topology"
	"strconv"
)

var graph *topology.Graph

//todo:
// 1 add edge provider
// 2 add api
// 3 documentation

func main() {

	graph = &topology.Graph{}
	initializingGraph()

	go topology.NetworkMetricsUpdate(graph)
	go topology.ClustersMetricsUpdate(graph)

	httpRouter := api.NewRouter(graph)

	httpServer := &http.Server{
		Handler: httpRouter,
		Addr:    ":" + config.GetConfiguration().ServicePort,
	}

	log.Fatal(httpServer.ListenAndServe())

}

func initializingGraph() {

	//add cells
	for i := 1; i <= 18; i++ {
		graph.AddVertex(topology.Vertex{Id: strconv.Itoa(i), Type: "CELL"})
	}
	//add MECs
	graph.AddVertex(topology.Vertex{Id: "mec1", Type: "MEC"})
	graph.AddVertex(topology.Vertex{Id: "mec5", Type: "MEC"})
	graph.AddVertex(topology.Vertex{Id: "mec7", Type: "MEC"})
	graph.AddVertex(topology.Vertex{Id: "mec15", Type: "MEC"})

	//addEdges
	graph.AddEdge(topology.Edge{Source: "mec1", Target: "2"})
	graph.AddEdge(topology.Edge{Source: "mec7", Target: "2"})
	graph.AddEdge(topology.Edge{Source: "mec1", Target: "8"})
	graph.AddEdge(topology.Edge{Source: "mec5", Target: "8"})
	graph.AddEdge(topology.Edge{Source: "mec7", Target: "8"})

	graph.GetVertex("mec1").VertexMetrics.UpdateClusterMetrics(topology.ClusterMetrics{12.231, 1.23})

	graph.PrintGraph()
	graph.PrintGraph()
}
