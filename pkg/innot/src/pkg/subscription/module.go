package subscription

import (
	log "10.254.188.33/matyspi5/erd/pkg/innot/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/innot/src/pkg/db"
	"10.254.188.33/matyspi5/erd/pkg/innot/src/pkg/types"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"math"
	"math/rand"
	"net/http"
	"strconv"
)

type counter float64

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
		sub, newCell := generateTargetCellId(id, sub)
		db.DummyDB.UpdateItem(id, sub)
		sendCellIdCellNotification(sub.Endpoint, types.CellId(strconv.Itoa(newCell)), id)

	} else {
		log.Errorf("[SUBSCRIPTION] Event Type %s doesn't exist!", sub.AmfEventType)
		err = errors.New(fmt.Sprintf("[SUBSCRIPTION] Exiting! Event Type %s doesn't exist!", sub.AmfEventType))
		db.INDEX--
		return
	}

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

// TODO: Try to make a path unique. Don't allow users to visit the same cell twice.
func generateTargetCellId(id db.SubscriptionId, sub db.Subscriber) (db.Subscriber, int) {
	if sub.CurrentCell == "" {
		initialState := initialPlacementStateMachine[id]
		sub.CurrentCell = types.CellId(strconv.Itoa(initialState))
		sub.UserPath = append(sub.UserPath, sub.CurrentCell)
	}
	var nextState int
	possibleStates := cellStateMachine[sub.CurrentCell]
	for {
		nextState = possibleStates[rand.Intn(len(possibleStates))]
		if len(sub.UserPath) < 2 {
			break
		}
		if types.CellId(strconv.Itoa(nextState)) == sub.UserPath[len(sub.UserPath)-2] {
			log.Warnf("[DEBUG] Chosen previous cell. Skipping...")
			continue
		}
		break
	}
	//possibleStates := cellStateMachine[sub.CurrentCell]
	//nextState := possibleStates[rand.Intn(len(possibleStates)-1)]
	log.Infof("[DEBUG] Candidate cells for CELL[%v] are [%v] chosen [%v]", sub.CurrentCell, possibleStates, nextState)
	sub.CurrentCell = types.CellId(strconv.Itoa(nextState))
	sub.UserPath = append(sub.UserPath, sub.CurrentCell)
	return sub, nextState
}

//var cellStateMachine = map[types.CellId][]int{
//	"1":  {2, 4, 5},
//	"2":  {1, 3, 4},
//	"3":  {2, 4, 7, 11, 14},
//	"4":  {1, 2, 3, 5, 6, 7},
//	"5":  {1, 4, 6, 8},
//	"6":  {4, 5, 7, 8, 9, 10},
//	"7":  {3, 4, 6, 10, 11, 14, 17},
//	"8":  {5, 6, 9},
//	"9":  {6, 8, 10},
//	"10": {6, 7, 9, 17, 20},
//	"11": {3, 12, 14},
//	"12": {11, 13, 14, 15, 16},
//	"13": {12, 16, 22},
//	"14": {3, 7, 11, 12, 15, 17},
//	"15": {12, 14, 16, 17, 18, 19},
//	"16": {12, 13, 15, 19, 22, 25},
//	"17": {7, 10, 14, 15, 18, 20},
//	"18": {15, 17, 19, 20, 21},
//	"19": {15, 16, 18, 21, 25, 29},
//	"20": {10, 17, 18},
//	"21": {18, 19, 29},
//	"22": {13, 16, 23, 25, 26},
//	"23": {22, 24, 26},
//	"24": {23, 26, 27, 32, 35},
//	"25": {16, 19, 22, 26, 28, 29},
//	"26": {22, 23, 24, 25, 27, 28},
//	"27": {24, 26, 28, 31, 35, 38},
//	"28": {25, 26, 27, 29, 30, 31},
//	"29": {19, 21, 25, 28, 30},
//	"30": {28, 29, 31},
//	"31": {27, 28, 30, 38, 41},
//	"32": {24, 33, 35},
//	"33": {32, 34, 35, 36, 37},
//	"34": {33, 37}, // consider only cell in 1st Coverage Zone
//	"35": {24, 27, 32, 33, 36, 38},
//	"36": {33, 35, 37, 38, 39, 40},
//	"37": {33, 34, 36, 40}, // consider only cell in 1st Coverage Zone
//	"38": {27, 31, 35, 36, 39, 41},
//	"39": {36, 38, 40, 41, 42},
//	"40": {36, 37, 39, 42}, // consider only cell in 1st Coverage Zone
//	"41": {31, 38, 39},
//	"42": {39, 40}, // consider only cell in 1st Coverage Zone
//}

//// TODO: Try to limit possible directions of movement
//var cellStateMachine = map[types.CellId][]int{
//	"1":  {2, 4},
//	"2":  {3, 4},
//	"3":  {7, 11, 14},
//	"4":  {2, 3, 6, 7},
//	"5":  {4, 6},
//	"6":  {7, 10},
//	"7":  {3, 10, 11, 14, 17},
//	"8":  {6, 9},
//	"9":  {10},
//	"10": {7, 17, 20},
//	"11": {3, 12, 14},
//	"12": {11, 13, 14, 16},
//	"13": {12, 16, 22},
//	"14": {3, 7, 11, 12, 15, 17},
//	"15": {14, 16, 17, 19},
//	"16": {12, 13, 15, 19, 22, 25},
//	"17": {7, 10, 14, 15, 18, 20},
//	"18": {15, 17, 19, 20, 21},
//	"19": {15, 16, 18, 21, 25, 29},
//	"20": {10, 17, 18},
//	"21": {18, 19, 29},
//	"22": {13, 16, 23, 25, 26},
//	"23": {22, 24, 26},
//	"24": {23, 26, 27, 32, 35},
//	"25": {16, 19, 22, 26, 28, 29},
//	"26": {22, 24, 25, 27},
//	"27": {24, 26, 28, 31, 35, 38},
//	"28": {25, 27, 29, 31},
//	"29": {19, 21, 25, 28, 30},
//	"30": {28, 29, 31},
//	"31": {27, 28, 30, 38, 41},
//	"32": {24, 33, 35},
//	"33": {32, 35, 36},
//	"34": {33}, // consider only cell in 1st Coverage Zone
//	"35": {24, 27, 32, 33, 36, 38},
//	"36": {35, 38},
//	"37": {33, 36}, // consider only cell in 1st Coverage Zone
//	"38": {27, 31, 35, 36, 39, 41},
//	"39": {36, 38, 41},
//	"40": {36, 39}, // consider only cell in 1st Coverage Zone
//	"41": {31, 38, 39},
//	"42": {39}, // consider only cell in 1st Coverage Zone
//}

// TODO: Try to further improve mobility model v1
var cellStateMachine = map[types.CellId][]int{
	"1":  {3},
	"5":  {2, 7}, // middle
	"8":  {1},
	"2":  {11, 14},
	"4":  {11, 14},
	"6":  {1, 2},
	"9":  {4, 7},
	"3":  {12},
	"7":  {12, 15}, // middle
	"10": {8},
	"11": {13, 16},
	"14": {13, 16},
	"17": {6, 9},
	"20": {6, 9},
	"12": {22},
	"15": {22, 25}, // middle
	"18": {10},
	"13": {23, 26},
	"16": {23, 26},
	"19": {17, 20},
	"21": {17, 20},
	"22": {24},
	"25": {15, 18, 24, 27}, // middle
	"29": {18},
	"23": {32, 35},
	"26": {32, 35},
	"28": {19, 21},
	"30": {19, 21},
	"24": {33},
	"27": {25, 29}, // middle
	"31": {29},
	"32": {34},
	"35": {37},
	"38": {28, 30},
	"41": {28, 30},
	"33": {39},
	"36": {27, 31}, // middle
	"39": {31},
	"34": {36, 40},
	"37": {39, 42},
	"40": {38, 41},
	"42": {38, 41},
}

// Initial placement v3.0.9
var initialPlacementStateMachine = map[db.SubscriptionId]int{
	1:  15,
	10: 42,
	11: 13,
	12: 20,
	13: 36,
	14: 12,
	15: 41,
	16: 20,
	17: 20,
	18: 4,
	19: 37,
	2:  28,
	20: 11,
	21: 35,
	22: 34,
	23: 39,
	24: 4,
	25: 35,
	26: 1,
	27: 16,
	28: 17,
	29: 1,
	3:  38,
	30: 36,
	31: 32,
	32: 31,
	33: 13,
	34: 6,
	35: 20,
	36: 25,
	37: 25,
	38: 5,
	39: 42,
	4:  11,
	40: 6,
	41: 39,
	42: 10,
	43: 11,
	44: 26,
	45: 29,
	46: 17,
	47: 1,
	48: 21,
	49: 27,
	5:  25,
	50: 30,
	6:  39,
	7:  42,
	8:  11,
	9:  25,
}

//// Initial placement v3.0.8
//var initialPlacementStateMachine = map[db.SubscriptionId]int{
//	1:  6,
//	10: 25,
//	11: 42,
//	12: 29,
//	13: 30,
//	14: 25,
//	15: 23,
//	16: 1,
//	17: 22,
//	18: 3,
//	19: 34,
//	2:  41,
//	20: 25,
//	21: 36,
//	22: 14,
//	23: 16,
//	24: 41,
//	25: 26,
//	26: 22,
//	27: 25,
//	28: 35,
//	29: 13,
//	3:  33,
//	30: 13,
//	31: 3,
//	32: 15,
//	33: 23,
//	34: 39,
//	35: 9,
//	36: 3,
//	37: 32,
//	38: 21,
//	39: 14,
//	4:  42,
//	40: 27,
//	41: 9,
//	42: 32,
//	43: 15,
//	44: 40,
//	45: 23,
//	46: 38,
//	47: 13,
//	48: 33,
//	49: 27,
//	5:  11,
//	50: 8,
//	6:  36,
//	7:  42,
//	8:  28,
//	9:  16,
//}

//// Initial placement v3.0.7
//var initialPlacementStateMachine = map[db.SubscriptionId]int{
//	1:  9,
//	10: 34,
//	11: 25,
//	12: 11,
//	13: 14,
//	14: 6,
//	15: 4,
//	16: 16,
//	17: 8,
//	18: 20,
//	19: 31,
//	2:  20,
//	20: 1,
//	21: 37,
//	22: 10,
//	23: 28,
//	24: 26,
//	25: 8,
//	26: 41,
//	27: 36,
//	28: 34,
//	29: 34,
//	3:  2,
//	30: 16,
//	31: 35,
//	32: 11,
//	33: 40,
//	34: 17,
//	35: 14,
//	36: 19,
//	37: 36,
//	38: 3,
//	39: 21,
//	4:  36,
//	40: 33,
//	41: 20,
//	42: 31,
//	43: 33,
//	44: 25,
//	45: 31,
//	46: 39,
//	47: 27,
//	48: 23,
//	49: 8,
//	5:  31,
//	50: 33,
//	6:  14,
//	7:  18,
//	8:  26,
//	9:  22,
//}

//// Initial placement v3.0.6
//var initialPlacementStateMachine = map[db.SubscriptionId]int {
//	1:  22,
//	10: 13,
//	11: 17,
//	12: 13,
//	13: 13,
//	14: 33,
//	15: 2,
//	16: 33,
//	17: 20,
//	18: 40,
//	19: 12,
//	2:  33,
//	20: 13,
//	21: 19,
//	22: 2,
//	23: 18,
//	24: 41,
//	25: 21,
//	26: 39,
//	27: 18,
//	28: 11,
//	29: 5,
//	3:  39,
//	30: 34,
//	31: 35,
//	32: 32,
//	33: 11,
//	34: 17,
//	35: 16,
//	36: 2,
//	37: 5,
//	38: 28,
//	39: 28,
//	4:  41,
//	40: 27,
//	41: 24,
//	42: 31,
//	43: 16,
//	44: 15,
//	45: 21,
//	46: 8,
//	47: 5,
//	48: 10,
//	49: 13,
//	5:  24,
//	50: 9,
//	6:  36,
//	7:  25,
//	8:  32,
//	9:  28,
//}

//// Initial placement v3.0.5
//var initialPlacementStateMachine = map[db.SubscriptionId]int{
//	1:  30,
//	10: 31,
//	11: 22,
//	12: 5,
//	13: 15,
//	14: 7,
//	15: 10,
//	16: 11,
//	17: 3,
//	18: 41,
//	19: 21,
//	2:  24,
//	20: 12,
//	21: 37,
//	22: 5,
//	23: 40,
//	24: 21,
//	25: 8,
//	26: 10,
//	27: 6,
//	28: 12,
//	29: 34,
//	3:  39,
//	30: 7,
//	31: 21,
//	32: 7,
//	33: 11,
//	34: 41,
//	35: 12,
//	36: 4,
//	37: 1,
//	38: 29,
//	39: 30,
//	4:  11,
//	40: 21,
//	41: 2,
//	42: 24,
//	43: 38,
//	44: 5,
//	45: 40,
//	46: 1,
//	47: 24,
//	48: 36,
//	49: 28,
//	5:  34,
//	50: 28,
//	6:  17,
//	7:  38,
//	8:  37,
//	9:  31,
//}

//// Initial placement v3.0.4
//var initialPlacementStateMachine = map[db.SubscriptionId]int{
//	1:  14,
//	10: 31,
//	11: 14,
//	12: 12,
//	13: 23,
//	14: 13,
//	15: 6,
//	16: 19,
//	17: 27,
//	18: 3,
//	19: 23,
//	2:  24,
//	20: 26,
//	21: 17,
//	22: 9,
//	23: 21,
//	24: 8,
//	25: 28,
//	26: 10,
//	27: 34,
//	28: 41,
//	29: 13,
//	3:  42,
//	30: 7,
//	31: 20,
//	32: 18,
//	33: 24,
//	34: 15,
//	35: 11,
//	36: 6,
//	37: 2,
//	38: 20,
//	39: 24,
//	4:  31,
//	40: 37,
//	41: 1,
//	42: 42,
//	43: 41,
//	44: 36,
//	45: 20,
//	46: 2,
//	47: 1,
//	48: 15,
//	49: 14,
//	5:  6,
//	50: 30,
//	6:  37,
//	7:  42,
//	8:  34,
//	9:  26,
//}

//// Initial placement v3.0.3
//var initialPlacementStateMachine = map[db.SubscriptionId]int{
//	1:  23,
//	10: 19,
//	11: 27,
//	12: 8,
//	13: 9,
//	14: 25,
//	15: 3,
//	16: 10,
//	17: 22,
//	18: 17,
//	19: 19,
//	2:  29,
//	20: 2,
//	21: 6,
//	22: 12,
//	23: 17,
//	24: 39,
//	25: 39,
//	26: 2,
//	27: 12,
//	28: 14,
//	29: 18,
//	3:  15,
//	30: 24,
//	31: 37,
//	32: 17,
//	33: 23,
//	34: 40,
//	35: 4,
//	36: 40,
//	37: 35,
//	38: 37,
//	39: 23,
//	4:  18,
//	40: 5,
//	41: 38,
//	42: 2,
//	43: 30,
//	44: 28,
//	45: 39,
//	46: 14,
//	47: 18,
//	48: 41,
//	49: 13,
//	5:  21,
//	50: 36,
//	6:  12,
//	7:  32,
//	8:  28,
//	9:  25,
//}

//// Initial placement v3.0.2
//var initialPlacementStateMachine = map[db.SubscriptionId]int{
//	1:  40,
//	10: 25,
//	11: 33,
//	12: 16,
//	13: 5,
//	14: 11,
//	15: 20,
//	16: 19,
//	17: 19,
//	18: 3,
//	19: 23,
//	2:  13,
//	20: 38,
//	21: 40,
//	22: 30,
//	23: 4,
//	24: 20,
//	25: 12,
//	26: 41,
//	27: 42,
//	28: 24,
//	29: 13,
//	3:  12,
//	30: 12,
//	31: 38,
//	32: 35,
//	33: 15,
//	34: 38,
//	35: 1,
//	36: 40,
//	37: 17,
//	38: 2,
//	39: 11,
//	4:  18,
//	40: 42,
//	41: 19,
//	42: 33,
//	43: 6,
//	44: 22,
//	45: 23,
//	46: 42,
//	47: 30,
//	48: 30,
//	49: 38,
//	5:  23,
//	50: 2,
//	6:  10,
//	7:  39,
//	8:  5,
//	9:  41,
//}

// Initial placement v3.0.1
//var initialPlacementStateMachine = map[db.SubscriptionId]int{
//	1:  37,
//	2:  31,
//	3:  9,
//	4:  40,
//	5:  4,
//	6:  24,
//	7:  3,
//	8:  21,
//	9:  18,
//	10: 12,
//	11: 21,
//	12: 18,
//	13: 1,
//	14: 39,
//	15: 24,
//	16: 17,
//	17: 21,
//	18: 33,
//	19: 40,
//	20: 29,
//	21: 35,
//	22: 3,
//	23: 11,
//	24: 14,
//	25: 16,
//	26: 22,
//	27: 33,
//	28: 34,
//	29: 15,
//	30: 13,
//	31: 33,
//	32: 30,
//	33: 4,
//	34: 11,
//	35: 21,
//	36: 29,
//	37: 2,
//	38: 24,
//	39: 8,
//	40: 36,
//	41: 27,
//	42: 21,
//	43: 21,
//	44: 5,
//	45: 8,
//	46: 25,
//	47: 15,
//	48: 34,
//	49: 14,
//	50: 8,
//}
