package main

import (
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/api"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/config"
	log "10.254.188.33/matyspi5/erd/pkg/nmt/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/mec-topology"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
)

var graph *mec_topology.Graph

func main() {
	log.Infof("[SERVER] Starting NMT server. Port: %v", config.GetConfiguration().ServicePort)
	log.Infof("[SERVER] OBS endpoint: %v", config.GetConfiguration().ClusterControllerEndpoint)

	graph = &mec_topology.Graph{}
	//	initializingGraph()

	graph.ReadTopologyConfigFile("mecTopology.json")
	graph.ReadMECConnectionFile("mecLinks.json")
	graph.ReadNetworkTopologyConfigFile("networkTopology.json")

	/*	var link model.Edge
		link.SourceVertexName = "mec12"
		link.TargetVertexName = "mec13"
		link.TargetVertexProviderName = "edge-provider"
		link.SourceVertexProviderName = "edge-provider"
		link.EdgeMetrics.Latency = 12.32

		graph.AddLink(link)
	*/
	//gorutines to update cluster resources and network metrics
	//go graph.NetworkMetricsUpdate()
	graph.NetworkMetricsUpdate(false)
	go graph.ClustersResourcesUpdate()

	httpRouter := api.NewRouter(graph)

	httpServer := &http.Server{
		Handler: httpRouter,
		Addr:    ":" + config.GetConfiguration().ServicePort,
	}

	connectionsClose := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		httpServer.Shutdown(context.Background())
		close(connectionsClose)
	}()

	err := httpServer.ListenAndServe()
	log.Fatalln(fmt.Sprintf("[SERVER] HTTP server returned error: %s", err))

}
