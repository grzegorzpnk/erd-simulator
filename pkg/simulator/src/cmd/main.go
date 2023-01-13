package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"simu/src/api"
	"simu/src/config"
	log "simu/src/logger"
	"simu/src/pkg/model"
)

func main() {
	log.Infof("[SERVER] Starting SIMU server. Port: %v", config.GetConfiguration().ServicePort)

	var simuClient *model.SimuClient

	simuClient.fetchAppsFromNMT()
	startSIMUserver(simuClient)

}

func startSIMUserver(simuClient *model.SimuClient) {

	httpRouter := api.NewRouter(simuClient)
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
