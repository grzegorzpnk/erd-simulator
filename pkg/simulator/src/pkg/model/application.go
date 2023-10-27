package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"simu/src/config"
	log "simu/src/logger"
)

type SimuClient struct {
	Apps []MECApp `json:"apps"`
}

type MECApp struct {
	Id           string             `json:"id"`
	ClusterId    string             `json:"clusterId"`
	UserLocation string             `json:"userLocation"`
	Requirements RequestedResources `json:"requirements"`
	UserPath     []string           `json:"userPath,omitempty""`
}

type RequestedResources struct {
	RequestedCPU     float64 `json:"requestedCPU"`
	RequestedMEMORY  float64 `json:"requestedMEMORY"`
	RequestedLatency float64 `json:"requestedLatency"`
}

func (app *MECApp) PrintApplication() {

	fmt.Printf("Application ID: %v, app clusterID: %v, UE Location at cell no:%v, requirements: CPU: %v, Memory: %v, Latency: %v \n", app.Id, app.ClusterId, app.UserLocation, app.Requirements.RequestedCPU, app.Requirements.RequestedMEMORY, app.Requirements.RequestedLatency)

}

func (simuCl *SimuClient) GetMecHost() error {

	url := buildNMTendpoint()

	apps, err := GetAppsFromNMT(url)
	if err != nil {
		return err
	}

	simuCl.SetApps(apps)
	simuCl.setPath()
	return nil
}

func (simuCl *SimuClient) FetchAppsFromNMT() error {

	url := buildNMTendpoint()

	apps, err := GetAppsFromNMT(url)
	if err != nil {
		return err
	}

	simuCl.SetApps(apps)
	simuCl.setPath()
	return nil
}

func (simuCl *SimuClient) RecreateInitialPlacementAtNMT() error {

	url := buildRecreateInitialPlacementNMTendpoint()

	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		log.Errorf("Could not make POST request. reason: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("response code [%s] for URL [%s]", resp.Status, url)
	}

	return nil
}

func GetAppsFromNMT(endpoint string) ([]MECApp, error) {

	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	//Convert the body to type string
	var apps []MECApp
	json.Unmarshal(body, &apps)

	return apps, nil

}

func (simuCl *SimuClient) SetApps(apps []MECApp) {
	simuCl.Apps = apps
}

func (simuCl *SimuClient) setPath() {

	for i, _ := range simuCl.Apps {
		simuCl.Apps[i].UserPath = append(simuCl.Apps[i].UserPath, simuCl.Apps[i].UserLocation)
		//v.UserPath = append(v.UserPath, v.UserLocation)
	}

}

func buildNMTendpoint() string {
	url := config.GetConfiguration().NMTEndpoint
	url += "/v1/topology/application"

	return url
}

func buildNMTGetMECendpoint() string {
	url := config.GetConfiguration().NMTEndpoint
	url += "/v1/topology/application"

	return url
}

func buildRecreateInitialPlacementNMTendpoint() string {
	url := config.GetConfiguration().NMTEndpoint
	url += "/v1/topology/recreate-initial"

	return url
}
func (simuCl *SimuClient) GetApps(id string) *MECApp {

	for i, v := range simuCl.Apps {
		if v.Id == id {
			return &simuCl.Apps[i]
		}
	}
	return nil
}
