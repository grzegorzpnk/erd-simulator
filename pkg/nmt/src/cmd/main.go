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

	graph = &mec_topology.Graph{}
	readGraphConfig()

	//need to be refactored for simulator
	graph.NetworkMetricsUpdate(false)

	//start NMT server

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

func readGraphConfig() {

	graph.ReadTopologyConfigFile("mecTopology.json")
	graph.ReadMECConnectionFile("mecLinks.json")
	graph.ReadNetworkTopologyConfigFile("networkTopology.json")

}
