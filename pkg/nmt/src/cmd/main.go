package main

import (
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/api"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/config"
	log "10.254.188.33/matyspi5/erd/pkg/nmt/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/mec-topology"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/model"
	"net/http"
)

var graph *mec_topology.Graph

//todo:
// 1 add edge provider
// 2 add api
// 3 documentation

func main() {
	log.Infof("[SERVER] Starting NMT server. Port: %v", config.GetConfiguration().ServicePort)
	log.Infof("[SERVER] OBS endpoint: %v", config.GetConfiguration().ClusterControllerEndpoint)

	graph = &mec_topology.Graph{}
	initializingGraph()

	//gorutines to update cluster resources and network metrics
	//go mec_topology.NetworkMetricsUpdate(graph)
	go graph.ClustersResourcesUpdate()

	httpRouter := api.NewRouter(graph)

	httpServer := &http.Server{
		Handler: httpRouter,
		Addr:    ":" + config.GetConfiguration().ServicePort,
	}

	log.Fatalln(httpServer.ListenAndServe())

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

func createMecHost(clusterName, clusterProvider string) model.MecHost {

	var mec model.MecHost
	mec.Identity.Cluster = clusterName
	mec.Identity.Provider = clusterProvider

	return mec
}

func createLink(startMecHost, startMecProvider, destMecHost, destMecProvider string) model.Edge {

	var link model.Edge
	link.SourceVertexName = startMecHost
	link.SourceVertexProviderName = startMecProvider
	link.TargetVertexName = destMecHost
	link.TargetVertexProviderName = destMecProvider

	return link
}
