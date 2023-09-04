package model

import (
	"fmt"
	"math/rand"
	"time"
)

type MECApp struct {
	Id           string             `json:"id"`
	ClusterId    string             `json:"clusterId"`
	UserLocation string             `json:"userLocation"`
	Requirements RequestedResources `json:"requirements"`
}

type RequestedResources struct {
	RequestedCPU     float64 `json:"requestedCPU"`
	RequestedMEMORY  float64 `json:"requestedMEMORY"`
	RequestedLatency float64 `json:"requestedLatency"`
}

func (app *MECApp) GeneratreResourceRequirements() {

	//// Create a slice containing the possible values
	//values := []int{500, 600, 700, 800, 900, 1000}
	//
	//// Select a random value from the slice using rand.Intn()
	//randomCPUIndex := rand.Intn(len(values))
	//randomMEMIndex := rand.Intn(len(values))
	//randomCPU := values[randomCPUIndex]
	//randomMemory := values[randomMEMIndex]

	rand.Seed(time.Now().UnixNano()) // Seed the random number generator with the current time
	randomCPU := rand.Intn(521) + 530
	randomMemory := rand.Intn(521) + 530

	app.Requirements.RequestedCPU = float64(randomCPU)
	app.Requirements.RequestedMEMORY = float64(randomMemory)

}

func (app *MECApp) PrintApplication() {

	fmt.Printf("Application ID: %v, app clusterID: %v, UE Location at cell no:%v, requirements: CPU: %v, Memory: %v, Latency: %v \n", app.Id, app.ClusterId, app.UserLocation, app.Requirements.RequestedCPU, app.Requirements.RequestedMEMORY, app.Requirements.RequestedLatency)

}
