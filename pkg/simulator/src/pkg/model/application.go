package model

import (
	"fmt"
	"simu/src/config"
	log "simu/src/logger"
)

type SimuClient struct {
	apps []MECApp `json:"apps"`
}

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

func (app *MECApp) PrintApplication() {

	fmt.Printf("Application ID: %v, app clusterID: %v, UE Location at cell no:%v, requirements: CPU: %v, Memory: %v, Latency: %v \n", app.Id, app.ClusterId, app.UserLocation, app.Requirements.RequestedCPU, app.Requirements.RequestedMEMORY, app.Requirements.RequestedLatency)

}

func (simuCLient *SimuClient) FetchAppsFromNMT() error {

	url := buildNMTendpoint()

	apps, err := GetAppsFromNMT(url)
	if err != nil {
		log.Errorf(err)
		return error(err)
	}

	simuCLient.setApps(apps)
	return nil
}

func GetAppsFromNMT(string url) {

}

func (simuCLient *SimuClient) setApps(apps []MECApp) {
	simuCLient.apps = apps

}

func buildNMTendpoint() string {
	url := config.GetConfiguration().NMTEndpoint
	url += "topology/application"

	return url
}


