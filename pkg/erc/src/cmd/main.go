// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package main

import (
	"context"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"10.254.188.33/matyspi5/erd/pkg/erc/src/api"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/config"
	log "10.254.188.33/matyspi5/erd/pkg/erc/src/logger"
	"github.com/gorilla/handlers"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

// run initializes the dependencies and start the controller
func run() error {
	rand.Seed(time.Now().UnixNano())

	//// Initialize database(s)
	//if err := initDataBases(); err != nil {
	//	return err
	//}
	//
	//// Initialize grpc server, if required
	//initGrpcServer()

	// Handle requests on incoming connections
	if err := serve(); err != nil {
		return err
	}

	return nil
}

//// initDataBases initializes the emco databases
//func initDataBases() error {
//	// Initialize the emco database(Mongo DB)
//	err := db.InitializeDatabaseConnection("emco")
//	if err != nil {
//		log.Errorf("Failed to initialize mongo database connection. Error: %v", err)
//		return err
//	}
//
//	// Initialize etcd
//	err = contextdb.InitializeContextDatabase()
//	if err != nil {
//		log.Errorf("Failed to initialize etcd database connection. Error: %v", err)
//		return err
//	}
//
//	return nil
//}

// serve start the controller and handle requests on incoming connections
func serve() error {
	p := config.GetConfiguration().ServicePort

	log.Infof("Starting controller. Port: %v.", p)

	r := api.NewRouter(nil)
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

	//c, err := auth.GetTLSConfig("ca.cert", "server.cert", "server.key")
	//if err != nil {
	//	log.Info("Failed to get the TLS configuration. Starting without TLS.")
	//	return server.ListenAndServe()
	//}
	//
	//server.TLSConfig = c
	//return server.ListenAndServeTLS("", "") // empty string. tlsconfig already has this information
	return server.ListenAndServe()
}

//// initGrpcServer start the gRPC server
//func initGrpcServer() {
//	go func() {
//		if err := grpc.StartGrpcServer(); err != nil {
//			log.Fatalln(err)
//		}
//	}()
//}
