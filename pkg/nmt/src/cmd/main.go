package main

import (
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/api"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/config"
	log "10.254.188.33/matyspi5/erd/pkg/nmt/src/logger"
	mec_topology "10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/mec-topology"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
)

var graph *mec_topology.Graph

func main() {
	log.Infof("[SERVER] Starting NMT server. Port: %v", config.GetConfiguration().ServicePort)

	graph = &mec_topology.Graph{}
	readMECNetworkTopologyConfig()
	// setup clusters resources
	graph.AssigneCapacityToClusters()

	// TODO: setup latency on links - to be included if we already created links !
	graph.NetworkMetricsUpdate()

	//for each MEC define shortest path ( lower latency) to reach each of cell
	//this has been introduced for ML mask immidiate calcualation
	graph.FindShortestPathsConfigurationForMECs()

	//start NMT server
	startNMTserver()

}

func readMECNetworkTopologyConfig() {

	graph.ReadTopologyConfigFile("mecTopologyforTwoCZ.json")
	graph.ReadMECConnectionFile("mecLinksforTwoCZ.json")
	graph.ReadNetworkTopologyConfigFile("networkTopologyforTwoCZ.json")

}

func startNMTserver() {

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
