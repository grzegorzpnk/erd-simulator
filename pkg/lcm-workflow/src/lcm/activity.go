// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package lcm

import (
	"10.254.188.33/matyspi5/erd/pkg/lcm-workflow/src/types"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var waitNotification = make(chan types.CellId)

func SubCellChangedNotification(ctx context.Context, migParam MigParam) (*MigParam, error) {

	log.Printf("SubCellChangedNotification got params: %#v\n", migParam)

	migParam.NotifyUrl = GenerateListenerEndpoint("workflow-listener", "cell-changed")
	log.Printf("\nSubCellChangedNotification: listener endpoint = %s\n", migParam.NotifyUrl)

	// TODO: un-hardcore innot url
	innotUrl := migParam.GetInnotUrl()

	log.Printf("\nSubCellChangedNotification: innotUrl = %s\n", innotUrl)

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

	httpServer := &http.Server{
		Handler: router,
		Addr:    ":8585",
	}

	go func() {
		migParam.NewCellId = <-waitNotification

		fmt.Printf("SubCellChangedNotification: Got notification. New CELL_ID: %v\n.", migParam.NewCellId)
		_ = httpServer.Shutdown(context.Background())
	}()

	fmt.Println("SubCellChangedNotification: Listening for notification..")
	_ = httpServer.ListenAndServe()

	return &migParam, nil
}

func serveNotification(w http.ResponseWriter, r *http.Request) {
	var info types.CellChangedInfo
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		fmt.Errorf("error while decoding: %v\n", err)
	}

	log.Printf("Notification reason: %v, cell id: %v", info.Reason, info.Cell)

	waitNotification <- info.Cell

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
		EventNotifyUri:                "http://localhost/workflow-listener/cell-changed/ABCDEFGHIJ/notify",
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
