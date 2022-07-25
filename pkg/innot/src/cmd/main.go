package main

import (
	"10.254.188.33/matyspi5/innot/config"
	"10.254.188.33/matyspi5/innot/src/api"
	log "10.254.188.33/matyspi5/innot/src/logger"
	ddb "10.254.188.33/matyspi5/innot/src/pkg/db"
	"fmt"

	"context"
	"github.com/gorilla/handlers"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	ddb.InitDummyDb()

	httpRouter := api.NewRouter()
	loggedRouter := handlers.LoggingHandler(os.Stdout, httpRouter)
	log.Info("[SERVER] Starting IN server")

	httpServer := &http.Server{
		Handler: loggedRouter,
		Addr:    ":" + config.GetConfiguration().ServicePort,
	}
	log.Infof("[SERVER] Intermediate Notifier HTTP server will listen at endpoint: %s", httpServer.Addr)
	log.Infof("[SERVER] AMF Endpoint: %s", config.GetConfiguration().AMFEndpoint)

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
