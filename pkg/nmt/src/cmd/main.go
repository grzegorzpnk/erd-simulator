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

	//add MECs
	mec := createMecHost("mec1", "edge-provider")
	graph.AddMecHost(mec)
	mec = createMecHost("mec5", "edge-provider")
	graph.AddMecHost(mec)
	mec = createMecHost("mec7", "edge-provider")
	graph.AddMecHost(mec)
	mec = createMecHost("mec15", "edge-provider")
	graph.AddMecHost(mec)

	//addEdges
	link := createLink("mec1", "edge-provider", "mec5", "edge-provider")
	graph.AddEdge(link)
	link = createLink("mec5", "edge-provider", "mec1", "edge-provider")
	graph.AddEdge(link)
	link = createLink("mec5", "edge-provider", "mec7", "edge-provider")
	graph.AddEdge(link)
	link = createLink("mec7", "edge-provider", "mec1", "edge-provider")
	graph.AddEdge(link)
	link = createLink("mec1", "edge-provider", "mec15", "edge-provider")
	graph.AddEdge(link)

	//graph.PrintGraph()

}

func createMecHost(clusterName, clusterProvider string) mec_topology.MecHost {

	var mec mec_topology.MecHost
	mec.Identity.ClusterName = clusterName
	mec.Identity.Provider = clusterProvider

	return mec
}

func createLink(startMecHost, startMecProvider, destMecHost, destMecProvider string) mec_topology.Edge {

	var link mec_topology.Edge
	link.SourceVertexName = startMecHost
	link.SourceVertexProviderName = startMecProvider
	link.TargetVertexName = destMecHost
	link.TargetVertexProviderName = destMecProvider

	return link
}
