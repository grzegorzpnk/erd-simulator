package api

import (
	"fmt"
	"simu/src/pkg/model"
	"strconv"
	"testing"
)

func TestGenerateUserToMove(t *testing.T) {
	// Create a mock struct for the SimuClient
	mockSimuClient := model.SimuClient{Apps: []model.MECApp{model.MECApp{Id: "1"}, model.MECApp{Id: "2"}, model.MECApp{Id: "3"}}}

	// Create a mock struct for the apiHandler
	mockApiHandler := &apiHandler{SimuClient: mockSimuClient}

	// Call the generateUserToMove function
	result := mockApiHandler.generateUserToMove()

	// Convert the result to an int
	resultInt, _ := strconv.Atoi(result)

	// Check if the result is within the range of the number of apps in the SimuClient
	if resultInt < 1 || resultInt > len(mockSimuClient.Apps) {
		t.Errorf("generateUserToMove returned an invalid result: %d", resultInt)
	} else {
		fmt.Printf(result)
	}
}
