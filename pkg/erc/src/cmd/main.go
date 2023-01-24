package main

import (
	"10.254.188.33/matyspi5/erd/pkg/erc/src/api"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/config"
	log "10.254.188.33/matyspi5/erd/pkg/erc/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/results"
	"context"
	"github.com/gorilla/handlers"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	rand.Seed(time.Now().UnixNano())

	if err := serve(); err != nil {
		return err
	}

	return nil
}

func serve() error {
	p := config.GetConfiguration().ServicePort
	log.Infof("Starting Smart Placement Controller. Port: %v.", p)

	var result *results.Client
	result = results.NewClient()

	r := api.NewRouter(nil, result)
	h := handlers.LoggingHandler(os.Stdout, r)
	server := &http.Server{
		Handler: h,
		Addr:    ":" + p,
	}

	connection := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		server.Shutdown(context.Background())
		close(connection)
	}()

	return server.ListenAndServe()
}
