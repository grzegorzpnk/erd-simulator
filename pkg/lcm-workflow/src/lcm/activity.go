// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package lcm

import (
	"10.254.188.33/matyspi5/erd/pkg/lcm-workflow/src/types"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var NEXT_PORT = 8585
var waitNotification chan int = make(chan int)

func SubCellChangedNotification(ctx context.Context, migParam MigParam) (*MigParam, error) {

	log.Printf("SubCellChangedNotification got params: %#v\n", migParam)

	migParam.NotifyUrl = GenerateListenerEndpoint("workflow-listener", "cell-changed")
	log.Printf("\nSubCellChangedNotification: listener endpoint = %s\n", migParam.NotifyUrl)

	// TODO: un-hardcore innot url
	innotUrl := "http://10.254.185.44:32137/v1/intermediate-notifier/subscribe"

	log.Printf("\nSubCellChangedNotification: innotUrl = %s\n", innotUrl)
	//if migParam.Listener == reflect.Zero(reflect.TypeOf(new(net.Listener))).Interface() {
	//	log.Printf("SubCellChangedNorification: 	Creating new listener...")
	//	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", strconv.Itoa(NEXT_PORT)))
	//	if err != nil {
	//		log.Fatalf("error: %v\n", err)
	//	}
	//	log.Printf("SubCellChangedNotification: Created listener on port: %v\n", listener.Addr().(*net.TCPAddr).Port)
	//	migParam.Listener = &listener
	//}

	//host := fmt.Sprintf("http://10.254.185.50:%v", listener.Addr().(*net.TCPAddr).Port)
	host := fmt.Sprintf("http://10.254.185.48:%v", os.Getenv("NOTIFICATION_NODE_PORT"))

	data := generateSubscriptionBody()
	data.EventNotifyUri = host + migParam.NotifyUrl

	log.Printf("SubCellChangedNotification: EventNorifyUri: %v\n", data.EventNotifyUri)

	err := postHttp(innotUrl, data)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("SubCellChangedNotification: Subscribed for notifications.")
	}

	return &migParam, nil
}

func GetCellChangedNotification(ctx context.Context, migParam MigParam) (*MigParam, error) {
	log.Printf("SubCellChangedNotification: activity start\n")

	router := mux.NewRouter()
	router.HandleFunc(migParam.NotifyUrl, serveNotification).Methods("POST")

	fmt.Printf("SubCellChangedNotification: Creating server\n")
	httpServer := &http.Server{
		Handler: router,
		Addr:    ":8585",
	}

	log.Printf("SubCellChangedNotification: registered handler\n")

	go func() {
		fmt.Printf("SubCellChangedNotification: Running new goroutine to handle server shutdown\n")
		<-waitNotification
		fmt.Printf("SubCellChangedNotification: Got notification on channel... HTTP server shutdown\n.")
		_ = httpServer.Shutdown(context.Background())
	}()

	fmt.Println("SubCellChangedNotification: Listening..")
	_ = httpServer.ListenAndServe()

	fmt.Printf("SubCellChangedNotification: activity ended.\n")
	return &migParam, nil
}

func serveNotification(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got notification: %v.", r.Body)
	waitNotification <- 1
}

func generateSubscriptionBody() types.AmfEventSubscription {
	boole := true

	amfEventArea := types.AmfEventArea{
		LadnInfo: &types.LadnInfo{
			Ladn:     "",
			Presence: nil,
		},
		PresenceInfo: &types.PresenceInfo{
			EcgiList:            nil,
			GlobalRanNodeIdList: nil,
			NcgiList:            nil,
			PraId:               nil,
			PresenceState:       nil,
			TrackingAreaList: &[]types.Tai{types.Tai{
				PlmnId: types.PlmnId{"208", "93"},
				Tac:    "000001",
			}},
		},
	}
	amfEvent := types.AmfEvent{
		AreaList:                 &[]types.AmfEventArea{amfEventArea},
		ImmediateFlag:            &boole,
		LocationFilterList:       nil,
		SubscribedDataFilterList: nil,
		Type:                     "LOCATION_REPORT",
	}

	body := types.AmfEventSubscription{
		AnyUE:                         &boole,
		EventList:                     &[]types.AmfEvent{amfEvent},
		EventNotifyUri:                "http://localhost/workflow-listener/cell-changed/LEoF2qn1jk/notify",
		Gpsi:                          nil,
		GroupId:                       nil,
		NfId:                          uuid.UUID{},
		NotifyCorrelationId:           "",
		Options:                       &types.AmfEventMode{Trigger: "ONE_TIME"},
		Pei:                           nil,
		SubsChangeNotifyCorrelationId: nil,
		SubsChangeNotifyUri:           nil,
		Supi:                          nil,
	}
	return body
}
