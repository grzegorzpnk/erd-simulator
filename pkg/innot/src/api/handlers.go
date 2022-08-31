package api

import (
	log "10.254.188.33/matyspi5/erd/pkg/innot/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/innot/src/pkg/db"
	"10.254.188.33/matyspi5/erd/pkg/innot/src/pkg/subscription"
	"10.254.188.33/matyspi5/erd/pkg/innot/src/pkg/types"
	"encoding/json"
	"net/http"
)

type apiHandler string

// subscribeHandler is a function which is called to handle new subscription
func (h apiHandler) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	var body types.AmfEventSubscription
	var eventType types.AmfEventType
	var endpoint string
	var subId db.SubscriptionId

	err := json.NewDecoder(r.Body).Decode(&body)

	log.Infof("GOT BODY: %v", body)

	// For now allow subscription for single EventType per request
	eventTypes := *body.EventList
	eventType = eventTypes[0].Type

	endpoint = body.EventNotifyUri

	sub := db.Subscriber{
		Endpoint:     types.ClientListenerUri(endpoint),
		AmfEventType: eventType,
		BodyRequest:  body,
	}

	subId, err = db.DummyDB.PutItem(sub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go subscription.ServeSubscription(subId)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(sub)
	if err != nil {
		log.Error("[API] Error encoding.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Infof("[API] Created new subscription! Endpoint: %v", sub.Endpoint)
}

// TODO not implemented
func (h apiHandler) unsubscribeByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
}

// TODO not implemented
func (h apiHandler) unsubscribeByEndpointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
}

// TODO not implementd
func (h apiHandler) getAllSubscriptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
}
