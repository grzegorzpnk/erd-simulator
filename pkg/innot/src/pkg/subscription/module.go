package subscription

import (
	log "10.254.188.33/matyspi5/erd/pkg/innot/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/innot/src/pkg/db"
	"10.254.188.33/matyspi5/erd/pkg/innot/src/pkg/types"
	"math"
	"strconv"
	"strings"

	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"math/rand"
	"net/http"
)

type counter float64

func normalizeCellId(cellStr string) (int, error) {
	// e.g. get["000000020"] -> return[2]
	if strings.HasSuffix(cellStr, "0") {
		cellStr = cellStr[:len(cellStr)-1]
	}
	cell, err := strconv.Atoi(cellStr)
	if err != nil {
		log.Errorf("Could not parse CELL_ID[%v] to integer value", cellStr)
		return -1, err
	}
	return cell, nil
}

// ServeSubscription method is invoked (as goroutine) when new db.Subscriber is created
// It will request notification from AMF (once per second) and if conditions are met, the notification to the MEO
// client is sent. For now there is no mechanism implemented, which will allow to cancel this subscription.
func ServeSubscription(id db.SubscriptionId) {
	var err error
	var sub db.Subscriber

	sub, err = db.DummyDB.GetItemByKey(id)
	if err != nil {
		return
	}

	err = verifyNotificationType(sub.AmfEventType)
	if err != nil {
		err = cancelSubscription(id, err)
		return
	}

	if sub.AmfEventType == types.LOCATIONREPORT {
		sub, newCell := generateTargetCellId(sub)
		db.DummyDB.UpdateItem(id, sub)
		sendCellIdCellNotification(sub.Endpoint, types.CellId(strconv.Itoa(newCell)), id)

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
	log.Infof("[SUBSCRIPTION][ID=%v] Sending CELL_ID_CHANGED[%v] to the client.", id, cellId)
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

func generateTargetCellId(sub db.Subscriber) (db.Subscriber, int) {
	if sub.CurrentCell == "" {
		initialState := rand.Intn(42-1) + 1
		sub.CurrentCell = types.CellId(strconv.Itoa(initialState))
		return sub, initialState
	} else {
		possibleStates := cellStateMachine[sub.CurrentCell]
		nextState := possibleStates[rand.Intn(len(possibleStates)-1)]
		log.Infof("[DEBUG] Candidate cells for CELL[%v] are [%v] chosen [%v]", sub.CurrentCell, possibleStates, nextState)
		sub.CurrentCell = types.CellId(strconv.Itoa(nextState))
		return sub, nextState
	}
}

var cellStateMachine = map[types.CellId][]int{
	"1":  {2, 4, 5},
	"2":  {1, 3, 4},
	"3":  {2, 4, 7, 11, 14},
	"4":  {1, 2, 3, 5, 6, 7},
	"5":  {1, 4, 6, 8},
	"6":  {4, 5, 7, 8, 9, 10},
	"7":  {3, 4, 6, 10, 11, 14, 17},
	"8":  {5, 6, 9},
	"9":  {6, 8, 10},
	"10": {6, 7, 9, 14, 17, 20},
	"11": {3, 12, 14},
	"12": {11, 13, 14, 15, 16},
	"13": {12, 16, 22},
	"14": {3, 7, 11, 12, 15, 17},
	"15": {12, 14, 16, 17, 18, 19},
	"16": {12, 13, 15, 19, 22, 25},
	"17": {7, 10, 14, 15, 18, 20},
	"18": {15, 17, 19, 20, 21},
	"19": {15, 16, 18, 21, 25, 29},
	"20": {10, 17, 18},
	"21": {18, 19, 29},
	"22": {13, 16, 23, 25, 26},
	"23": {22, 24, 26},
	"24": {23, 26, 27, 32, 35},
	"25": {16, 19, 22, 26, 28, 29},
	"26": {22, 23, 24, 25, 27, 28},
	"27": {24, 26, 28, 31, 35, 38},
	"28": {25, 26, 27, 29, 30, 31},
	"29": {19, 21, 25, 28, 30},
	"30": {28, 29, 31},
	"31": {27, 28, 30, 38, 41},
	"32": {24, 33, 35},
	"33": {32, 34, 35, 36, 37},
	"34": {33, 37}, // consider only cell in 1st Coverage Zone
	"35": {24, 27, 32, 33, 36, 38},
	"36": {33, 35, 37, 38, 38, 40},
	"37": {33, 34, 36, 40}, // consider only cell in 1st Coverage Zone
	"38": {27, 31, 35, 36, 39, 41},
	"39": {36, 38, 40, 41, 42},
	"40": {36, 37, 39, 42}, // consider only cell in 1st Coverage Zone
	"41": {31, 38, 39},
	"42": {39, 40}, // consider only cell in 1st Coverage Zone
}
