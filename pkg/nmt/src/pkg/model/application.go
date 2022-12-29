package model

import (
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/config"
	"fmt"
	"math/rand"
	"strconv"
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

	minRes, _ := strconv.ParseFloat(config.GetConfiguration().ResMin, 64)
	maxRes, _ := strconv.ParseFloat(config.GetConfiguration().ResMax, 64)

	app.Requirements.RequestedCPU = minRes + rand.Float64()*(maxRes-minRes)
	app.Requirements.RequestedMEMORY = minRes + rand.Float64()*(maxRes-minRes)

}

func (app *MECApp) PrintApplication() {

	fmt.Printf("Application ID: %v, app clusterID: %v, requirements: CPU: %v, Memory: %v, Latency: %v \n", app.Id, app.ClusterId, app.Requirements.RequestedCPU, app.Requirements.RequestedMEMORY, app.Requirements.RequestedLatency)

}
