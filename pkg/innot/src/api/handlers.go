package api

import (
	log "10.254.188.33/matyspi5/erd/pkg/innot/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/innot/src/pkg/db"
	"10.254.188.33/matyspi5/erd/pkg/innot/src/pkg/results"
	"10.254.188.33/matyspi5/erd/pkg/innot/src/pkg/subscription"
	"10.254.188.33/matyspi5/erd/pkg/innot/src/pkg/types"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type apiHandler struct {
	resClient *results.Client
}

func (h *apiHandler) SetClients() {
	h.resClient = results.NewClient()
}

// subscribeHandler is a function which is called to handle new subscription
func (h *apiHandler) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	var body types.AmfEventSubscription
	var eventType types.AmfEventType
	var endpoint string
	//var subId db.SubscriptionId

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("GOT BODY: %v", body)

	if body == (types.AmfEventSubscription{}) {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "Subscription body not found!", http.StatusInternalServerError)
		return
	}

	// For now allow subscription for single EventType per request
	eventTypes := *body.EventList
	eventType = eventTypes[0].Type

	endpoint = body.EventNotifyUri

	sub := db.Subscriber{
		Endpoint:     types.ClientListenerUri(endpoint),
		AmfEventType: eventType,
		//BodyRequest:  body,
	}

	_, err = db.DummyDB.PutItem(sub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//go subscription.ServeSubscription(subId)

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

// OLD IMPLEMENTATION
//func (h apiHandler) subscribeHandler(w http.ResponseWriter, r *http.Request) {
//	var body types.AmfEventSubscription
//	var eventType types.AmfEventType
//	var endpoint string
//	var subId db.SubscriptionId
//
//	err := json.NewDecoder(r.Body).Decode(&body)
//	if err != nil {
//		w.Header().Set("Content-Type", "application/json")
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	log.Infof("GOT BODY: %v", body)
//
//	if body == (types.AmfEventSubscription{}) {
//		w.Header().Set("Content-Type", "application/json")
//		http.Error(w, "Subscription body not found!", http.StatusInternalServerError)
//		return
//	}
//
//	// For now allow subscription for single EventType per request
//	eventTypes := *body.EventList
//	eventType = eventTypes[0].Type
//
//	endpoint = body.EventNotifyUri
//
//	sub := db.Subscriber{
//		Endpoint:     types.ClientListenerUri(endpoint),
//		AmfEventType: eventType,
//		BodyRequest:  body,
//	}
//
//	subId, err = db.DummyDB.PutItem(sub)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	go subscription.ServeSubscription(subId)
//
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusCreated)
//	err = json.NewEncoder(w).Encode(sub)
//	if err != nil {
//		log.Error("[API] Error encoding.")
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	log.Infof("[API] Created new subscription! Endpoint: %v", sub.Endpoint)
//}

// TODO not implemented
func (h *apiHandler) unsubscribeByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
}

// TODO not implemented
func (h *apiHandler) unsubscribeByEndpointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *apiHandler) handleSubscriptionHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	subIdStr := params["subscription-id"]

	subId, err := strconv.Atoi(subIdStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't parse subscription-id[%v]", subIdStr), http.StatusInternalServerError)
		return
	}

	//sub, err := db.DummyDB.GetItemByKey(db.SubscriptionId(subId))
	//if err != nil {
	//	http.Error(w, fmt.Sprintf("Couldn't find subciption for subscription-id[%v]", subIdStr), http.StatusNoContent)
	//}

	subscription.ServeSubscription(db.SubscriptionId(subId))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err != nil {
		log.Error("[API] Error encoding.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *apiHandler) getAllSubscriptionsHandler(w http.ResponseWriter, r *http.Request) {

	subs := db.DummyDB.GetItems()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(subs)
	if err != nil {
		log.Error("[API] Error encoding.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Results collection handlers

func (h *apiHandler) relocationFailedHandler(w http.ResponseWriter, r *http.Request) {

	h.resClient.Results.IncFailed()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func (h *apiHandler) relocationSuccessfulHandler(w http.ResponseWriter, r *http.Request) {

	h.resClient.Results.IncSuccessful()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func (h *apiHandler) relocationRedundantHandler(w http.ResponseWriter, r *http.Request) {

	h.resClient.Results.IncRedundant()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func (h *apiHandler) getResultsHandler(w http.ResponseWriter, r *http.Request) {

	subs := h.resClient.Results

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(subs)
	if err != nil {
		log.Error("[API] Error encoding.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
