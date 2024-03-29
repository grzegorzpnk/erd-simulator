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
	"simu/src/pkg/results"
)

func main() {
	log.Infof("[SERVER] Starting SIMU server. Port: %v", config.GetConfiguration().ServicePort)

	//var simuClient *model.SimuClient
	simuClient := &model.SimuClient{}
	resultsClient := results.NewClient()

	/*	err := simuClient.FetchAppsFromNMT()
		if err != nil {
			log.Errorf(err.Error())
		} else {
			log.Infof("Initial app list fetched from NMT")
		}*/

	startSIMUserver(simuClient, resultsClient)

}

func startSIMUserver(sClient *model.SimuClient, rClient *results.Client) {

	httpRouter := api.NewRouter(sClient, rClient)
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
