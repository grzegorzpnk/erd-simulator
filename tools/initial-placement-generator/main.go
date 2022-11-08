package main

import (
	"10.254.188.33/matyspi5/erd/initial-placement-generator/model"
	"10.254.188.33/matyspi5/erd/initial-placement-generator/topology"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"time"
)

const (
	TopologyShortestpathEndpoint = "http://10.254.185.44:32139/v1/topology/cells/:cell-id/mecHosts/provider/:provider/cluster/:cluster/shortest-path"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var inputFilePath string
	flag.StringVar(&inputFilePath, "input", "", "Input file to make transformation")
	flag.Parse()

	// Fetch yaml file
	var inputFileYaml model.InputYaml = model.InputYaml{}
	inputFileYaml.GetYamlFile(inputFilePath)

	tc := topology.NewTopologyClient()

	// Modify yaml file
	inputFileYaml = generateInitialClusters(tc, inputFileYaml)

	// Save yaml to the file
	//fmt.Printllsn(inputFileYaml)
	fmt.Println("----- SAVING -----")
	inputFileYaml.SaveYamlFile("values.yaml")

}

func generateInitialClusters(tc *topology.Client, in model.InputYaml) model.InputYaml {
	var mhs []model.MecHost

	for _, provider := range in.Providers {
		for _, cluster := range provider.Clusters {
			mh, err := tc.GetMecHost(provider.Name, cluster.Name)
			if err != nil {
				log.Errorf("Could not get MEC[%v+%v]", provider.Name, cluster.Name)
			}
			mh, err = tc.CollectResourcesInfo(mh)
			if err != nil {
				log.Errorf("Could not collect resources info MEC[%v+%v]", provider.Name, cluster.Name)
			}
			mhs = append(mhs, mh)
		}
	}

	var mhsT = make([]model.MecHost, len(mhs))
	var cells = map[int]int{}
	search := true
	for search {
		fmt.Println("[DEBUG] Starting search.")
		copy(mhsT, mhs)
		cells = generateRandomCells()

		search = false
		for index, edgeApp := range in.Deployments {
			cmh, err := findCanidateMec(tc, edgeApp.Apps[0].Workflows[0].Params, cells[index+1], mhsT)
			if err != nil {
				search = true
				fmt.Printf("Could not find candidate mec for App[%v]. Search failed.\n", edgeApp.Apps[0].Name)
				break
			} else {
				fmt.Printf("Found candidate mec: %v\n", cmh)
			}
			mhsT = updateMecResourcesInfo(mhsT, cmh, edgeApp.Apps[0].Workflows[0].Params)
			in.Deployments[index].Apps[0].PlacementClusters[0].Clusters = []string{cmh.Identity.Cluster}
		}
	}
	printCellsInfo(cells)

	return in
}

func updateMecResourcesInfo(mhs []model.MecHost, cmh model.MecHost, eap model.EdgeAppParams) []model.MecHost {
	cmh.Resources.Cpu.Used += eap.AppCPUReq
	cmh.Resources.Memory.Used += eap.AppMemReq
	cmh.Resources.Cpu.Utilization = 100 * (cmh.Resources.Cpu.Used / cmh.Resources.Cpu.Allocatable)
	cmh.Resources.Memory.Utilization = 100 * (cmh.Resources.Memory.Used / cmh.Resources.Memory.Allocatable)

	for index, mh := range mhs {
		if mh.Identity.Provider == cmh.Identity.Provider && mh.Identity.Cluster == cmh.Identity.Cluster {
			mhs[index] = cmh
			break
		}
	}

	return mhs
}

func findCanidateMec(tc *topology.Client, eap model.EdgeAppParams, cell int, mhs []model.MecHost) (model.MecHost, error) {
	candidates := []model.MecHost{}

	for index, mh := range mhs {
		latency, _ := tc.GetShortestPath(model.CellId(strconv.Itoa(cell)), mh)

		mhs[index].Resources.Latency = latency
	}

	for _, mh := range mhs {
		if resourcesOk(eap, mh) && latencyOk(eap, mh) {
			candidates = append(candidates, mh)
			// We need to prioritize N+1 and N+2 level clusters
			if mh.Identity.Cluster == "mec1" {
				candidates = append(candidates, mh)
				candidates = append(candidates, mh)
				candidates = append(candidates, mh)
			} else if mh.Identity.Cluster == "mec3" || mh.Identity.Cluster == "mec4" || mh.Identity.Cluster == "mec5" ||
				mh.Identity.Cluster == "mec6" || mh.Identity.Cluster == "mec7" {
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

func generateRandomCells() map[int]int {
	var cells map[int]int = map[int]int{}
	for i := 1; i <= 50; i++ {
		cells[i] = rand.Intn(42) + 1
	}
	return cells
}

// latencyOk checks if latency constraints specified in intent (i) are met
func latencyOk(eap model.EdgeAppParams, mec model.MecHost) bool {
	latency := mec.GetLatency()
	latencyMax := eap.LatencyMax

	if mec.GetLatency() < 0 {
		return false
	} else if latencyMax > latency {
		return true
	} else {
		return false
	}
}

// resourcesOk checks if resource constraints specified in intent (i) are met
func resourcesOk(eap model.EdgeAppParams, mec model.MecHost) bool {
	var cpuUtilization, memUtilization float64

	cpuUtilization = 100 * (mec.GetCpuUsed() + eap.AppCPUReq) / mec.GetCpuAllocatable()
	memUtilization = 100 * (mec.GetMemUsed() + eap.AppMemReq) / mec.GetMemAllocatable()

	cpuMax := eap.CPUUtilMax // 80
	memMax := eap.MemUtilMax // 80
	cpuMecAvaliable := mec.GetCpuAllocatable() - mec.GetCpuUsed()
	memMecAvaliable := mec.GetMemAllocatable() - mec.GetMemUsed()

	if cpuUtilization < 0 || memUtilization < 0 {
		return false
	} else if cpuMecAvaliable < eap.AppCPUReq {
		return false
	} else if memMecAvaliable < eap.AppMemReq {
		return false
	} else if cpuMax >= cpuUtilization && memMax >= memUtilization {
		//log.Warnf("[RES-CHECK][DEBUG] Resources OK!")
		return true
	} else {
		//log.Warnf("[RES-CHECK][DEBUG] Resources not OK :/")
		return false
	}
}

func printCellsInfo(val interface{}) {
	jsonCells, err := json.Marshal(val)
	if err != nil {
		log.Errorf("Marshal err: %v", err)
	}

	fmt.Println("----- CELLS -----")
	fmt.Println(string(jsonCells))
}
