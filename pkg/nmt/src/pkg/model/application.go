package model

import (
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/config"
	"fmt"
	"math/rand"
	"strconv"
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
	//
	////todo set precision of float 64 in order to reduce the computation
	//minRes, _ := strconv.ParseFloat(config.GetConfiguration().ResMin, 64)
	//maxRes, _ := strconv.ParseFloat(config.GetConfiguration().ResMax, 64)
	//rand.Float64()
	//
	//app.Requirements.RequestedCPU = minRes + rand.Float64()*(maxRes-minRes)
	//app.Requirements.RequestedMEMORY = minRes + rand.Float64()*(maxRes-minRes)

	minRes, _ := strconv.Atoi(config.GetConfiguration().ResMin)
	maxRes, _ := strconv.Atoi(config.GetConfiguration().ResMax)

	app.Requirements.RequestedCPU = float64(minRes + rand.Intn(maxRes-minRes))
	app.Requirements.RequestedMEMORY = float64(minRes + rand.Intn(maxRes-minRes))

}

func (app *MECApp) PrintApplication() {

	fmt.Printf("Application ID: %v, app clusterID: %v, UE Location at cell no:%v, requirements: CPU: %v, Memory: %v, Latency: %v \n", app.Id, app.ClusterId, app.UserLocation, app.Requirements.RequestedCPU, app.Requirements.RequestedMEMORY, app.Requirements.RequestedLatency)

}
