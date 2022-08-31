package db

import (
	log "10.254.188.33/matyspi5/erd/pkg/innot/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/innot/src/pkg/types"
	"errors"
	"fmt"
)

type SubscriptionId int

// INDEX is used to generate Subscription IDs
var INDEX SubscriptionId = 0

// DummyDB is an instance of internal database.
// It's represented by map[subscriptionId]Subscriber
var DummyDB *subscriptionDb

// Subscriber struct represents a single subscriber which is stored in the db.DummyDB
// Endpoint is of type: types.ClientListenerUri (string) and represents Subscriber URI
// on which client will be listening for notifications.
type Subscriber struct {
	Endpoint     types.ClientListenerUri
	AmfEventType types.AmfEventType
	BodyRequest  types.AmfEventSubscription
}

// subscriptionDb is a struct which represents subscription database
// As key we use SubscriptionId (Integer) and as value Subscriber struct
type subscriptionDb struct {
	db map[SubscriptionId]Subscriber
}

// newDummyDb creates new Database and returns instance of that DB
func newDummyDb() *subscriptionDb {
	return &subscriptionDb{
		db: map[SubscriptionId]Subscriber{},
	}
}

// InitDummyDb initializes database, where subscription data will be stored
func InitDummyDb() {
	log.Info("[DB] Initializing dummy database.")
	DummyDB = newDummyDb()
	log.Info("[DB] Dummy database initialized!")
}

// GetItemByKey gets subscription data from the db.DummyDB database, based on subscription ID
func (ddb *subscriptionDb) GetItemByKey(key SubscriptionId) (Subscriber, error) {
	var err error

	sub, ok := ddb.db[key]
	if ok {
		log.Infof("[DB] GETTING Endpoint: %v, EventType: %v", sub.Endpoint, sub.AmfEventType)
		return sub, nil
	} else {
		err = errors.New(fmt.Sprintf("KEY %s doesn't exist", key))
		log.Errorf("[DB] Could not GET subscription: %v", err)
		return Subscriber{}, err
	}
}

// PutItem saves subscription data to the db.DummyDB database
func (ddb *subscriptionDb) PutItem(sub Subscriber) (SubscriptionId, error) {
	var err error

	for _, val := range ddb.db {
		if val.Endpoint == sub.Endpoint {
			if val.AmfEventType == sub.AmfEventType {
				err = errors.New(fmt.Sprintf("Endpoint %s already subscribed", sub.Endpoint))
				log.Errorf("[DB] Could not add subscriber to the database: %v", err)
				return -1, err
			}
		}
	}
	INDEX++
	ddb.db[INDEX] = sub

	log.Infof("[DB] ADDED subscriber to the database. Subscription ID:: %v", INDEX)
	return INDEX, nil
}

// DeleteItemByKey requires some mechanism which will allow to cancel
// goroutine which is serving specific subscription -> for example channel/context
// For now it's used to delete Subscription when some error occurred
func (ddb *subscriptionDb) DeleteItemByKey(id SubscriptionId) error {
	var err error

	_, err = ddb.GetItemByKey(id)
	if err != nil {
		return err
	}

	delete(ddb.db, id)
	return nil

}

// DeleteItemByEndpointAndEventType requires some mechanism which will allow to cancel
// goroutine which is serving specific subscription -> for example channel/context
// TODO
func (ddb *subscriptionDb) DeleteItemByEndpointAndEventType(e types.ClientListenerUri, et types.AmfEventType) error {
	return nil
}
