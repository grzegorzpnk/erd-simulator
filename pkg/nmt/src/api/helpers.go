package api

import (
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/config"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/djikstra"
	mec_topology "10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/mec-topology"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/model"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"strconv"

	"fmt"
)

func containsAnyEdge(vertex model.MecHost) bool {

	if vertex.Neighbours != nil {
		return true
	} else {
		return false
	}

}

type ShortestPathResult struct {
	latencyResults float64
	path           []string
}

func ShortestPath(startCell *model.Cell, destCluster *model.MecHost, graph *mec_topology.Graph) float64 {

	var min float64

	if destCluster == nil {
		log.Fatalln("destination MEC host not recognized!")

	}
	//check if they are direct neighbours, if so the latency is just between start and stop node
	if destCluster.CheckMECsupportsCell(startCell.Id) {
		min = destCluster.GetCell(startCell.Id).Latency
		fmt.Printf("direct nodes, latency between cell: %v and mec: [%v+%v], is: %v", startCell.Id, destCluster.Identity.Provider, destCluster.Identity.Cluster, destCluster.GetCell(startCell.Id).Latency)

	} else {
		// if not, we have to calculate path between all MEC clusters that are in the same local zone as cell, to the target cluster, the final latency is a sum of the calculated one + between started mec and cell
		var startClusters []model.MecHost

		for _, v := range graph.MecHosts {
			if v.Identity.Location.LocalZone == startCell.LocalZone {
				startClusters = append(startClusters, *v)
			}
		}

		var inputGraph djikstra.InputGraph
		inputGraph.Graph = make([]djikstra.InputData, 200)

		//add all mec hosts to temp graph todo: should be only subset of graph nodes
		for i, v := range graph.Edges {
			inputGraph.Graph[i].Source = v.SourceVertexName
			inputGraph.Graph[i].Destination = v.TargetVertexName
			inputGraph.Graph[i].Weight = v.EdgeMetrics.Latency
		}
		itemGraph := djikstra.CreateGraph(inputGraph)

		//calculate shortest path between all []startClusters and stopNode, where startClusters is a list of cluster directly associated with cell
		results := make([]ShortestPathResult, 0)

		for _, v := range startClusters {

			startNd := djikstra.Node{v.Identity.Cluster}
			stopNd := djikstra.Node{destCluster.Identity.Cluster}

			var resultTmp ShortestPathResult
			resultTmp.path, resultTmp.latencyResults = djikstra.GetShortestPath(&startNd, &stopNd, itemGraph)

			//add latency between cell and start MEC host
			resultTmp.latencyResults += graph.GetMecHost(v.Identity.Cluster, v.Identity.Provider).GetCell(startCell.Id).Latency

			results = append(results, resultTmp)

		}

		//find minimal value
		min := results[0].latencyResults
		for _, v := range results {
			if v.latencyResults < min {
				min = v.latencyResults
			}
		}

		for i, v := range results {

			if v.latencyResults == min {
				fmt.Printf("final path is: %v\n", results[i].path)
			}
		}

	}
	return min
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

func findCanidateMec(app model.MECApp, cell *model.Cell, mhs []model.MecHost, graph *mec_topology.Graph) (model.MecHost, error) {
	candidates := []model.MecHost{}

	for index, mh := range mhs {
		latency := ShortestPath(cell, &mh, graph)
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
	return candidates[rand.Intn(len(candidates))], nil
}

// latencyOk checks if latency constraints specified in intent (i) are met
func latencyOk(app model.MECApp, mec model.MecHost) bool {
	latency := mec.TmpLatency
	latencyMax := app.Requirements.RequestedCPU

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

	cpuUtilization = 100 * (mec.GetCpuUsed() + app.Requirements.RequestedCPU) / mec.GetCpuCapacity()
	memUtilization = 100 * (mec.GetMemoryUsed() + app.Requirements.RequestedMEMORY) / mec.GetMemoryCapacity()

	cpuMax := mec.CpuResources.Capacity * tau    // 80
	memMax := mec.MemoryResources.Capacity * tau // 80
	cpuMecAvaliable := mec.GetCpuCapacity() - mec.GetCpuUsed()
	memMecAvaliable := mec.GetMemoryCapacity() - mec.GetMemoryUsed()

	if cpuUtilization < 0 || memUtilization < 0 {
		return false
	} else if cpuMecAvaliable < app.Requirements.RequestedCPU {
		return false
	} else if memMecAvaliable < app.Requirements.RequestedMEMORY {
		return false
	} else if cpuMax >= cpuUtilization && memMax >= memUtilization {
		//log.Warnf("[RES-CHECK][DEBUG] Resources OK!")
		return true
	} else {
		//log.Warnf("[RES-CHECK][DEBUG] Resources not OK :/")
		return false
	}
}

func generateRandomCells() map[int]int {
	var cells map[int]int = map[int]int{}
	for i := 1; i <= 50; i++ {
		cells[i] = rand.Intn(42) + 1
	}
	return cells
}

func printCellsInfo(val interface{}) {
	jsonCells, err := json.Marshal(val)
	if err != nil {
		log.Fatal("Marshal err: %v", err)
	}

	fmt.Println("----- CELLS -----")
	fmt.Println(string(jsonCells))
}
