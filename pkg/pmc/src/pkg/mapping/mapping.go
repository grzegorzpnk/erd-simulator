package mapping

import (
	log "10.254.188.33/matyspi5/pmc/src/logger"
	"10.254.188.33/matyspi5/pmc/src/pkg/promql"
	"strconv"
	"time"
)

const PROVIDER = "multus-erd-orange"

var NODE_CLUSTER_MAP = map[string]string{
	PROVIDER + "+" + "ran01":  "10.21.1.116:9100",
	PROVIDER + "+" + "core01": "10.31.1.129:9100",
	PROVIDER + "+" + "meh01":  "10.41.1.168:9100",
	PROVIDER + "+" + "meh02":  "10.61.1.121:9100",
}

type Node struct {
	InternalName   string
	ExternalName   string
	CpuUtilisation float64
	RamUtilisation float64
}

func newNode(iName, eName string, cpu, ram float64) Node {
	return Node{
		InternalName:   iName,
		ExternalName:   eName,
		CpuUtilisation: cpu,
		RamUtilisation: ram,
	}
}

func (n *Node) SetCpuUtil(val float64) {
	n.CpuUtilisation = val
}

func (n *Node) GetCpuUtil() float64 {
	return n.CpuUtilisation
}

func (n *Node) SetRamUtil(val float64) {
	n.RamUtilisation = val
}

func (n *Node) GetRamUtil() float64 {
	return n.RamUtilisation
}

type NodesInfo struct {
	client promql.PromQL
	nodes  []Node
}

func (ni *NodesInfo) String() {

}

func (ni *NodesInfo) InitializeNodesInfo(client promql.PromQL) {
	log.Info("[MAPPING] Initializing NodesInfo...")
	ni.client = client

	for key, value := range NODE_CLUSTER_MAP {
		cpu, err := ni.client.GetCpuUtilisation(value)
		if err != nil {
			log.Errorf("Error: %v", err)
		}

		ram, err := ni.client.GetRamUtilisation(value)
		if err != nil {
			log.Errorf("Error: %v", err)
		}

		cpuFloat, _ := strconv.ParseFloat(cpu, 64)
		ramFloat, _ := strconv.ParseFloat(ram, 64)
		node := newNode(value, key, cpuFloat, ramFloat)
		ni.nodes = append(ni.nodes, node)
	}
	go ni.updateNodesInfo()
}

func (ni *NodesInfo) updateNodesInfo() {
	log.Info("[MAPPING] Updating NodesInfo...")
	ni.client.Time = time.Now()

	for id, node := range ni.nodes {
		cpu, err := ni.client.GetCpuUtilisation(node.InternalName)
		if err != nil {
			log.Errorf("[MAPPING] Error: %v", err)
		}

		ram, err := ni.client.GetRamUtilisation(node.InternalName)
		if err != nil {
			log.Errorf("[MAPPING] Error: %v", err)
		}

		cpuFloat, _ := strconv.ParseFloat(cpu, 64)
		ramFloat, _ := strconv.ParseFloat(ram, 64)

		ni.nodes[id].SetCpuUtil(cpuFloat)
		ni.nodes[id].SetRamUtil(ramFloat)
	}
	log.Infof("[MAPPING] Current state: %v", ni.nodes)
	time.Sleep(5 * time.Second)

	ni.updateNodesInfo()
}
