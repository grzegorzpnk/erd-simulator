package subscription

import (
	"10.254.188.33/matyspi5/erd/pkg/innot/src/config"
	log "10.254.188.33/matyspi5/erd/pkg/innot/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/innot/src/pkg/db"
	"10.254.188.33/matyspi5/erd/pkg/innot/src/pkg/types"
	"math"

	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
)

type counter float64

// ServeSubscription method is invoked (as goroutine) when new db.Subscriber is created
// It will request notification from AMF (once per second) and if conditions are met, the notification to the MEO
// client is sent. For now there is no mechanism implemented, which will allow to cancel this subscription.
func ServeSubscription(id db.SubscriptionId) {
	var respBody types.AmfCreatedEventSubscription
	var err error
	var sub db.Subscriber

	var lcounter counter = 0

	currentCell := ""
	sub, err = db.DummyDB.GetItemByKey(id)
	if err != nil {
		return
	}
	amfEndpoint := config.GetConfiguration().AMFEndpoint

	err = verifyNotificationType(sub.AmfEventType)
	if err != nil {
		err = cancelSubscription(id, err)
		return
	}

	for {
		respBody, err = getNotification(sub, amfEndpoint, id)

		if err != nil {
			err = cancelSubscription(id, err)
			return
		}
		if sub.AmfEventType == types.LOCATIONREPORT {
			lcounter.logWaitingForNotification(id, sub.AmfEventType)

			if respBody.ReportList == nil {
				err = errors.New("ReportList is empty")
				log.Errorf("[SUBSCRIPTION][ID=%v] Error: %v. UE probably is not registered!", id, err)
				err = cancelSubscription(id, err)
				return
			}

			reportList := *respBody.ReportList

			if currentCell == "" {
				currentCell = reportList[0].Location.NrLocation.Ncgi.NrCellId
			}
			newCell := reportList[0].Location.NrLocation.Ncgi.NrCellId
			if currentCell == newCell {
				time.Sleep(1 * time.Second)
				continue
			} else {
				sendCellIdCellNotification(sub.Endpoint, types.CellId(newCell), id)
				currentCell = newCell
			}
		} else if sub.AmfEventType == types.ACCESSTYPEREPORT {
			log.Infof("[SUBSCRIPTION][ID=%v] Event Type %s is not implemented. Aborting...", id, sub.AmfEventType)
			return
		} else if sub.AmfEventType == types.COMMUNICATIONFAILUREREPORT {
			log.Infof("[SUBSCRIPTION][ID=%v] Event Type %s is not implemented. Aborting...", id, sub.AmfEventType)
			return
		} else if sub.AmfEventType == types.CONNECTIVITYSTATEREPORT {
			log.Infof("[SUBSCRIPTION][ID=%v] Event Type %s is not implemented. Aborting...", id, sub.AmfEventType)
			return
		} else if sub.AmfEventType == types.PRESENCEINAOIREPORT {
			log.Infof("[SUBSCRIPTION][ID=%v] Event Type %s is not implemented. Aborting...", id, sub.AmfEventType)
			return
		} else if sub.AmfEventType == types.REACHABILITYREPORT {
			log.Infof("[SUBSCRIPTION][ID=%v] Event Type %s is not implemented. Aborting...", id, sub.AmfEventType)
			return
		} else if sub.AmfEventType == types.REGISTRATIONSTATEREPORT {
			log.Infof("[SUBSCRIPTION][ID=%v] Event Type %s is not implemented. Aborting...", id, sub.AmfEventType)
			return
		} else if sub.AmfEventType == types.SUBSCRIBEDDATAREPORT {
			log.Infof("[SUBSCRIPTION][ID=%v] Event Type %s is not implemented. Aborting...", id, sub.AmfEventType)
			return
		} else if sub.AmfEventType == types.SUBSCRIPTIONIDADDITION {
			log.Infof("[SUBSCRIPTION][ID=%v] Event Type %s is not implemented. Aborting...", id, sub.AmfEventType)
			return
		} else if sub.AmfEventType == types.SUBSCRIPTIONIDCHANGE {
			log.Infof("[SUBSCRIPTION][ID=%v] Event Type %s is not implemented. Aborting...", id, sub.AmfEventType)
			return
		} else if sub.AmfEventType == types.TIMEZONEREPORT {
			log.Infof("[SUBSCRIPTION][ID=%v] Event Type %s is not implemented. Aborting...", id, sub.AmfEventType)
			return
		} else if sub.AmfEventType == types.UESINAREAREPORT {
			log.Infof("[SUBSCRIPTION][ID=%v] Event Type %s is not implemented. Aborting...", id, sub.AmfEventType)
			return
		} else {
			log.Errorf("[SUBSCRIPTION] Event Type %s doesn't exist!", sub.AmfEventType)
			err = errors.New(fmt.Sprintf("[SUBSCRIPTION] Exiting! Event Type %s doesn't exist!", sub.AmfEventType))
			db.INDEX--
			return
		}
	}
	return
}

//getNotification requests notification from AMF and returns the response body
func getNotification(sub db.Subscriber, amfEndpoint string, id db.SubscriptionId) (types.AmfCreatedEventSubscription, error) {
	var respBody types.AmfCreatedEventSubscription

	subBody := types.AmfCreateEventSubscription{
		Subscription:      sub.BodyRequest,
		SupportedFeatures: nil,
	}
	reqBody, err := json.Marshal(subBody)

	resp, err := http.Post(amfEndpoint, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Could not get notification for: EventType: %s, AMF endpoint: %s",
			sub.AmfEventType, sub.Endpoint))
		return types.AmfCreatedEventSubscription{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &respBody)
	if err != nil {
		err = errors.Wrap(err, "Failed to unmarshal body")
		log.Errorf("[SUBSCRIPTION][ID=%v] Error: %v", id, err)
		return types.AmfCreatedEventSubscription{}, err
	}
	return respBody, nil
}

// sendCellIdCellNotification will send notification to the types.ClientListenerUri that types.CellId changed.
// This is sent to the clients which subscribed to that kind of subscription.
func sendCellIdCellNotification(subscriberEndpoint types.ClientListenerUri, cellId types.CellId, id db.SubscriptionId) {
	log.Infof("[SUBSCRIPTION][ID=%v] Sending CELL_ID_CHANGED to the client.", id)
	info := types.CellChangedInfo{
		Reason: types.CELLCHANGED,
		Cell:   cellId,
	}

	body, err := json.Marshal(info)

	resp, err := http.Post(string(subscriberEndpoint), "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Errorf("[SUBSCRIPTION][ID=%v] Error: %v", id, err)
		return
	}

	log.Infof("[SUBSCRIPTION][ID=%v] Sent notification of type: %s. Status code: %v", id, types.CELLCHANGED, resp.StatusCode)
}

// verifyNotificationType checks if provided types.AmfEventType is a real Event Type
func verifyNotificationType(eventType types.AmfEventType) error {
	for _, val := range types.EVENT_TYPES {
		if eventType == val {
			return nil
		}
	}
	err := errors.New(fmt.Sprintf("Event Type %v not supported.", eventType))
	return err
}

// cancelSubscription will delete subscription from the db.DummyDB. It can be used to handle errors, but
// it will not cancel subscription if it wasn't interrupted
// TODO implement mechanism, which will allow real cancellation: channels/contexts
func cancelSubscription(id db.SubscriptionId, err error) error {
	err2 := db.DummyDB.DeleteItemByKey(id)

	if err2 != nil {
		log.Errorf("[SUBSCRIPTION][ID=%v] Could not cancel subscription: %v", id, err)
		return err
	}

	log.Errorf("[SUBSCRIPTION][ID=%v] Subscription cancelled. Reason: %v", id, err)
	return err
}

// logWaitingForNotification is used to log information about waiting for the event.
// It's used to log only once per 60 requests for single Subscription.
func (c *counter) logWaitingForNotification(id db.SubscriptionId, et types.AmfEventType) {
	if math.Mod(float64(*c), 60) == 0 {
		log.Infof("[SUBSCRIPTION][ID=%v] Handling %v. Waiting for event", id, et)
		*c++
	} else {
		*c++
	}
}
