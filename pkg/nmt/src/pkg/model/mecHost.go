package model

import (
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/config"
	log "10.254.188.33/matyspi5/erd/pkg/nmt/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/metrics"
	"errors"
	"strconv"
)

type MecInfo string
type MecType int

type MecHost struct {
	Identity        MecIdentity              `json:"identity"`
	CpuResources    metrics.ClusterResources `json:"cpu_resources,omitempty"`
	MemoryResources metrics.ClusterResources `json:"memory_resources,omitempty"`
	Neighbours      []MecIdentity
	SupportingCells []Cell   `json:"supporting_cells"`
	MECApps         []MECApp `json:"mec_apps,omitempty"`
	TmpLatency      float64
}

const (
	MecCpuUtil MecInfo = "/cpu"
	MecMemUtil MecInfo = "/memory"
)

const (
	MecLocal MecType = iota
	MecRegional
	MecCentral
)

type MecIdentity struct {
	Provider string      `json:"provider"`
	Cluster  string      `json:"cluster"`
	Location MecLocation `json:"location"`
}

type MecResInfo struct {
	Used        float64 `json:"used"`        // How many cpu/memory used (value)
	Allocatable float64 `json:"allocatable"` // How many cpu/memory available (value)
	Utilization float64 `json:"utilization"` // How much is cpu/memory utilized (percentage)
}

type MecLocation struct {
	Level  MecType `json:"type"`
	Region string  `json:"region"`
	Zone   string  `json:"zone,omitempty"`
	//city-level eg. west, "" if type different from MecLocal and MecRegional
	LocalZone string `json:"local-zone,omitempty"` // eg. wroclaw, "" if type different from MecLocal
}

type Cell struct {
	Id        string  `json:"id"`
	Latency   float64 `json:"latency"`
	LocalZone string  `json:"local-zone,omitempty"`
}

func (m *MecHost) GetCell(cellId string) *Cell {

	for i, v := range m.SupportingCells {
		if v.Id == cellId {
			return &m.SupportingCells[i]
		}
	}
	return nil
}

//check if cluster supports given Cell ID
func (mec *MecHost) CheckMECsupportsCell(cellId string) bool {

	for _, v := range mec.SupportingCells {
		if v.Id == cellId {
			return true
			break
		}
	}
	return false
}

func (mec *MecHost) CheckEnoughResources(app MECApp) bool {

	tau, _ := strconv.ParseFloat(config.GetConfiguration().Tau, 64)
	availableCPU := mec.CpuResources.Capacity - mec.CpuResources.Used
	availableMem := mec.MemoryResources.Capacity - mec.MemoryResources.Used
	if !(app.Requirements.RequestedCPU < (availableCPU * tau)) {
		return false
	}
	if !(app.Requirements.RequestedMEMORY < (availableMem * tau)) {
		return false
	}

	return true
}

func (mec *MecHost) InstantiateApp(app MECApp) error {

	tau, _ := strconv.ParseFloat(config.GetConfiguration().Tau, 64)
	if (mec.CpuResources.Used+app.Requirements.RequestedCPU)/mec.CpuResources.Capacity > tau ||
		(mec.MemoryResources.Used+app.Requirements.RequestedMEMORY)/mec.MemoryResources.Capacity > tau {
		err := "Cannot initialize app %v on cluster %v, due to overload!!"
		return errors.New(err)
	}

	//add app to list
	mec.MECApps = append(mec.MECApps, app)

	//update resources on mec host

	//update used resources
	mec.CpuResources.Used += app.Requirements.RequestedCPU
	mec.MemoryResources.Used += app.Requirements.RequestedMEMORY

	//update utilization
	mec.CpuResources.Utilization = mec.CpuResources.Used / mec.CpuResources.Capacity
	mec.MemoryResources.Utilization = mec.MemoryResources.Used / mec.MemoryResources.Capacity

	log.Infof("App: %v instantiated on cluster: %v", app.Id, mec.Identity.Cluster)
	return nil

}

//TOBE tested
func (mec *MecHost) CheckAppExists(app MECApp) bool {

	for _, v := range mec.MECApps {
		if v.Id == app.Id {
			return true
		}
	}
	return false
}

//TOBE tested
func (mec *MecHost) UninstallApp(app MECApp) {

	//delete app to list
	var initialID int

	//check id where the app is in slice
	for i, v := range mec.MECApps {
		if v.Id == app.Id {
			initialID = i
		}
	}

	//delete given app
	mec.MECApps[initialID] = mec.MECApps[len(mec.MECApps)-1]
	mec.MECApps = mec.MECApps[:len(mec.MECApps)-1]

	if len(mec.MECApps) == 0 {
		mec.MECApps = nil
	}

	//update resources on mec host
	//update used resources
	mec.CpuResources.Used -= app.Requirements.RequestedCPU
	mec.MemoryResources.Used -= app.Requirements.RequestedMEMORY

	//update utilization
	mec.CpuResources.Utilization = mec.CpuResources.Used / mec.CpuResources.Capacity
	mec.MemoryResources.Utilization = mec.MemoryResources.Used / mec.MemoryResources.Capacity

	log.Infof("App: %v UNinstalled from cluster: %v", app.Id, app.ClusterId)

}

func (mec *MecHost) GetCpuUsed() float64 {
	return mec.CpuResources.Used
}

func (mec *MecHost) GetMemoryUsed() float64 {
	return mec.MemoryResources.Used
}

func (mec *MecHost) GetCpuCapacity() float64 {
	return mec.CpuResources.Capacity
}

func (mec *MecHost) GetMemoryCapacity() float64 {
	return mec.MemoryResources.Capacity
}
