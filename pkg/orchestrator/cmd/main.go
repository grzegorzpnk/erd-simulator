package main

import (
	"orchestrator/api"
	"orchestrator/config"
	log "orchestrator/logger"

	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	log.Infof("[SERVER] Starting Orchestrator service. Port: %v", config.GetConfiguration().ServicePort)
	log.Infof("[SERVER] NMT endpoint: %v", config.GetConfiguration().TopologyEndpoint)

	httpRouter := api.NewRouter()

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
