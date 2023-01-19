package api

import (
	"math/rand"
	log "simu/src/logger"
	"simu/src/pkg/model"
	"strconv"
	"time"
)

func (h *apiHandler) generateUserToMove() string {

	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(len(h.SimuClient.Apps)) + 1)
}

// TODO: Try to make a path unique. Don't allow users to visit the same cell twice.
func (h *apiHandler) generateTargetCellId(app *model.MECApp) int {

	var nextState int

	user := app

	possibleStates := cellStateMachine[user.UserLocation]
	for {
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

	log.Infof("[DEBUG] Candidate cells for CELL[%v] are [%v] chosen [%v]", user.UserLocation, possibleStates, nextState)
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
	"33": {39},
	"36": {27, 31}, // middle
	"39": {31},
	"34": {36, 40},
	"37": {39, 42},
	"40": {38, 41},
	"42": {38, 41},
}
