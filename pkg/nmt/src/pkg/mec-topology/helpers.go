package mec_topology

import (
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/config"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/model"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
)

func (g *Graph) CheckGraphContainsVertex(mecHost model.MecHost) bool {

	for _, v := range g.MecHosts {
		if mecHost.Identity.Cluster == v.Identity.Cluster &&
			mecHost.Identity.Provider == v.Identity.Cluster {
			return true
		}
	}
	return false
}

// this func checks in bidirectional way
func (g *Graph) CheckAlreadExistLink(k model.Edge) bool {

	for _, v := range g.Edges {
		if (k.SourceVertexName == v.SourceVertexName &&
			k.SourceVertexProviderName == v.SourceVertexProviderName &&
			k.TargetVertexName == v.TargetVertexName &&
			k.TargetVertexProviderName == v.TargetVertexProviderName) ||
			(k.SourceVertexName == v.TargetVertexName &&
				k.SourceVertexProviderName == v.TargetVertexProviderName &&
				k.TargetVertexName == v.SourceVertexName &&
				k.TargetVertexProviderName == v.SourceVertexProviderName) {
			return true
		}

	}
	return false

}

func FindCanidateMec(app model.MECApp, cell *model.Cell, mhs []model.MecHost, graph *Graph) (model.MecHost, error) {
	candidates := []model.MecHost{}

	fmt.Printf("Looking for MEC Candidates for:    ")
	app.PrintApplication()

	for index, mh := range mhs {
		latency, err := graph.ShortestPath(cell, &mh)
		if err != nil {
			err := "Find MEC Candidate has failed due to '0' value shortest path"
			fmt.Printf(err)
			return model.MecHost{}, errors.New(err)
		}
		mhs[index].TmpLatency = latency
	}

	for _, mh := range mhs {
		if resourcesOk(app, mh) && latencyOk(app, mh) {
			candidates = append(candidates, mh)
			// We need to prioritize N+1 and N+2 level clusters
			if mh.Identity.Location.Level == 2 {
				candidates = append(candidates, mh)
				candidates = append(candidates, mh)
				candidates = append(candidates, mh)
			} else if mh.Identity.Location.Level == 1 {
				candidates = append(candidates, mh)
				candidates = append(candidates, mh)
			}
		} else {
		}
	}
	if len(candidates) <= 0 {
		return model.MecHost{}, errors.New("no candidates found")
	}
	fmt.Printf("Total number of candidates: %v\n", len(candidates))
	for i, _ := range candidates {
		fmt.Println(candidates[i].Identity.Cluster)
	}
	return candidates[rand.Intn(len(candidates))], nil
}

func updateMecResourcesInfo(mhs []model.MecHost, cmh model.MecHost, app model.MECApp) []model.MecHost {
	cmh.CpuResources.Used += app.Requirements.RequestedCPU
	cmh.MemoryResources.Used += app.Requirements.RequestedMEMORY
	cmh.CpuResources.Utilization = 100 * (cmh.CpuResources.Used / cmh.CpuResources.Capacity)
	cmh.MemoryResources.Utilization = 100 * (cmh.MemoryResources.Used / cmh.MemoryResources.Capacity)

	for index, mh := range mhs {
		if mh.Identity.Provider == cmh.Identity.Provider && mh.Identity.Cluster == cmh.Identity.Cluster {
			mhs[index] = cmh
			break
		}
	}

	return mhs
}

// latencyOk checks if latency constraints specified in intent (i) are met
func latencyOk(app model.MECApp, mec model.MecHost) bool {
	latency := mec.TmpLatency
	latencyMax := app.Requirements.RequestedLatency

	if mec.TmpLatency < 0 {
		return false
	} else if latencyMax > latency {
		return true
	} else {
		return false
	}
}

// resourcesOk checks if resource constraints specified in intent (i) are met
func resourcesOk(app model.MECApp, mec model.MecHost) bool {
	var cpuUtilization, memUtilization float64

	tau, _ := strconv.ParseFloat(config.GetConfiguration().Tau, 64)

	cpuUtilization = (mec.GetCpuUsed() + app.Requirements.RequestedCPU) / mec.GetCpuCapacity()
	memUtilization = (mec.GetMemoryUsed() + app.Requirements.RequestedMEMORY) / mec.GetMemoryCapacity()

	cpuMecAvaliable := mec.GetCpuCapacity() - mec.GetCpuUsed()
	memMecAvaliable := mec.GetMemoryCapacity() - mec.GetMemoryUsed()

	if cpuUtilization < 0 || memUtilization < 0 {
		return false
	} else if cpuMecAvaliable < app.Requirements.RequestedCPU {
		return false
	} else if memMecAvaliable < app.Requirements.RequestedMEMORY {
		return false
	} else if cpuUtilization <= tau && memUtilization <= tau {
		//log.Warnf("[RES-CHECK][DEBUG] Resources OK!")
		return true
	} else {
		//log.Warnf("[RES-CHECK][DEBUG] Resources not OK :/")
		return false
	}
}

func GenerateRandomCellsForUsers(usersNumber int, graph Graph) map[int]int {
	var cells map[int]int = map[int]int{}
	max, _ := strconv.Atoi(config.GetConfiguration().MaxCellNumber)
	for i := 1; i <= usersNumber; i++ {
		cells[i] = rand.Intn(max) + 1
	}
	return cells
}

type ShortestPathResult struct {
	latencyResults float64
	path           []string
}
