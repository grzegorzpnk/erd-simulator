package model

import (
	"fmt"
	"math/rand"
)

type MECApp struct {
	Id           string             `json:"id"`
	Requirements RequestedResources `json:"requirements"`
	ClusterId    string             `json:"clusterId"`
}

type RequestedResources struct {
	RequestedCPU     float64 `json:"requestedCPU"`
	RequestedMEMORY  float64 `json:"requestedMEMORY"`
	RequestedLatency float64 `json:"requestedLatency"`
}

func (app *MECApp) GeneratreResourceRequirements() {

	//rand.Seed(time.Now().UnixNano())

	// Generate a random number between 100 and 200

	app.Requirements.RequestedCPU = 100.00 + rand.Float64()*101.00
	app.Requirements.RequestedMEMORY = 100 + rand.Float64()*101.00

}

func (app *MECApp) PrintApplication() {

	fmt.Printf("Application ID: %v, app clusterID: %v, requirements: CPU: %v, Memory: %v, Latency: %v \n", app.Id, app.ClusterId, app.Requirements.RequestedCPU, app.Requirements.RequestedMEMORY, app.Requirements.RequestedLatency)

}
