package main

import (
	"10.254.188.33/matyspi5/pmc/api"
	"10.254.188.33/matyspi5/pmc/config"
	log "10.254.188.33/matyspi5/pmc/src/logger"
	"10.254.188.33/matyspi5/pmc/src/pkg/observability"
	"10.254.188.33/matyspi5/pmc/src/pkg/promql"
	"fmt"
	"github.com/gorilla/handlers"
	"golang.org/x/net/context"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	log.Info("[SERVER] Starting PMC server")

	cApi := promql.PromQL{
		Host:    config.GetConfiguration().Host,
		Timeout: config.GetConfiguration().Timeout * time.Second,
		Time:    time.Now(),
		Client:  nil,
	}

	client, err := promql.NewClient(cApi.Host)
	if err != nil {
		log.Errorf("[SERVER] Could not create new prometheus client: %v", err)
		return
	}
	cApi.Client = client

	var nodes observability.NodesInfo
	nodes.InitializeNodesInfo(cApi)

	var clusters observability.ClustersInfo
	clusters.InitializeClustersInfo(cApi)

	//// Testing...
	//fmt.Println("Testing...")
	//cpu1, cpu2 := cApi.GetCpuUtilisationNatively("meh3-worker")
	//fmt.Printf("%.2f : %.2f\n", cpu1, cpu2)

	httpRouter := api.NewRouter(&clusters)
	loggedRouter := handlers.LoggingHandler(os.Stdout, httpRouter)

	httpServer := &http.Server{
		Handler: loggedRouter,
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

	err = httpServer.ListenAndServe()
	log.Fatalln(fmt.Sprintf("[SERVER] HTTP server returned error: %s", err))
}
