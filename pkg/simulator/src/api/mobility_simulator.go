package api

import (
	"math/rand"
	log "simu/src/logger"
	"simu/src/pkg/model"
	"strconv"
)

func (h *apiHandler) generateUserToMove() string {

	//rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(len(h.SimuClient.Apps)) + 1)
}

// TODO: Try to make a path unique. Don't allow users to visit the same cell twice.
func (h *apiHandler) generateTargetCellId(app *model.MECApp) int {

	var nextState int

	user := app
	possibleStates := cellStateMachine[user.UserLocation]
	for {
		log.Infof("Possible states for app: %v : %v, len: %v", app.Id, possibleStates, len(possibleStates))
		nextState = possibleStates[rand.Intn(len(possibleStates))]
		if len(user.UserPath) < 2 {
			break
		}
		//todo: check if this is correct
		if strconv.Itoa(nextState) == user.UserPath[len(user.UserPath)-1] {
			log.Warnf("[DEBUG] Chosen previous cell. Skipping...")
			continue
		}
		break
	}
	//log.Infof("User(app) with ID: %v moved FROM cell: %v, towards cell: %v", user.Id, user.UserLocation, nextState)
	user.UserLocation = strconv.Itoa(nextState)
	user.UserPath = append(user.UserPath, user.UserLocation)
	return nextState
}

// TODO: Try to further improve mobility model v1
var cellStateMachine = map[string][]int{
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
	"33": {39, 43},
	"36": {27, 31}, // middle
	"39": {31, 50},
	"34": {36, 40},
	"37": {39, 42, 46},
	"40": {38, 48},
	"42": {48, 51},
	"43": {45, 49},
	"44": {48, 53, 34},
	"45": {54, 58, 47},
	"46": {37, 57},
	"47": {52, 39},
	"48": {44, 61},
	"49": {54, 58},
	"50": {49, 52, 36},
	"51": {61, 64},
	"52": {58, 62},
	"53": {55, 58},
	"54": {45, 64},
	"55": {65, 68},
	"56": {64, 67},
	"57": {59, 63, 46},
	"58": {59, 68},
	"59": {48, 51},
	"60": {52, 70},
	"61": {71, 73},
	"62": {48, 51},
	"63": {62, 73},
	"64": {66, 69},
	"65": {55, 74},
	"66": {75, 78},
	"67": {57, 69},
	"68": {58, 77},
	"69": {78, 81},
	"70": {60},
	"71": {61, 80},
	"72": {70, 81},
	"73": {63, 83},
	"74": {76, 79},
	"75": {69, 81},
	"76": {82, 80},
	"77": {76, 79},
	"78": {69, 72},
	"79": {84, 80},
	"80": {71, 73},
	"81": {72, 79},
	"82": {80, 83},
	"83": {73, 70},
	"84": {83, 79},
}
