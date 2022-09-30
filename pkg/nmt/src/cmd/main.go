package main

import (
	"log"
	"net/http"
	"nmt/src/api"
	"nmt/src/config"
	"nmt/src/package/mec-topology"
)

var graph *mec_topology.Graph

//todo:
// 1 add edge provider
// 2 add api
// 3 documentation

func main() {

	graph = &mec_topology.Graph{}
	initializingGraph()

	//gorutines to update cluster resources and network metrics
	//go mec_topology.NetworkMetricsUpdate(graph)
	go mec_topology.ClustersResourcesUpdate(graph)

	httpRouter := api.NewRouter(graph)

	httpServer := &http.Server{
		Handler: httpRouter,
		Addr:    ":" + config.GetConfiguration().ServicePort,
	}

	log.Fatal(httpServer.ListenAndServe())

}

func initializingGraph() {

	//add cells
	/*	for i := 1; i <= 18; i++ {
			graph.AddMecHost(mec_topology.MecHost{Id: strconv.Itoa(i), Type: "CELL"})
		}
		//add MECs
		graph.AddVertex(mec_topology.Vertex{Id: "mec1", Type: "MEC"})
		graph.AddVertex(mec_topology.Vertex{Id: "mec5", Type: "MEC"})
		graph.AddVertex(mec_topology.Vertex{Id: "mec7", Type: "MEC"})
		graph.AddVertex(mec_topology.Vertex{Id: "mec15", Type: "MEC"})

		//addEdges
		graph.AddEdge(mec_topology.Edge{Source: "mec1", Target: "2"})
		graph.AddEdge(mec_topology.Edge{Source: "mec7", Target: "2"})
		graph.AddEdge(mec_topology.Edge{Source: "mec1", Target: "8"})
		graph.AddEdge(mec_topology.Edge{Source: "mec5", Target: "8"})
		graph.AddEdge(mec_topology.Edge{Source: "mec7", Target: "8"})

		graph.GetVertex("mec1").VertexMetrics.UpdateClusterMetrics(mec_topology.ClusterMetrics{12.231, 1.23})

		graph.PrintGraph()
		graph.PrintGraph()*/
}
