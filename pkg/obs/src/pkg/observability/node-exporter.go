// Package observability - node-exporter: I don't think that we will use node-exporter. We should consider K8s
// parameters such as requests and limits due to fact, that real node CPU/memory utilization can be low (e.g. 10%), but
// requests may be about 100%. If it's true, pods won't run, and they will stay in Pending state.
package observability

import (
	log "10.254.188.33/matyspi5/obs/src/logger"
	"10.254.188.33/matyspi5/obs/src/pkg/promql"

	"strconv"
	"time"
)

const (
	PROVIDER = "orange-node-exporter"

	CLUSTER1 = "ran01"
	CLUSTER2 = "core01"
	CLUSTER3 = "meh01"
	CLUSTER4 = "meh02"
)

var NODE_EXPORTER_CLUSTER_MAP = map[string]string{
	PROVIDER + "+" + CLUSTER1: "10.21.1.116:9100",
	PROVIDER + "+" + CLUSTER2: "10.31.1.129:9100",
	PROVIDER + "+" + CLUSTER3: "10.41.1.168:9100",
	PROVIDER + "+" + CLUSTER4: "10.61.1.121:9100",
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

func (ni *NodesInfo) InitializeNodesInfo(client promql.PromQL) {
	log.Info("[NODE-EXPORTER] Initializing NodesInfo...")
	ni.client = client

	for key, host := range NODE_EXPORTER_CLUSTER_MAP {
		cpu, err := ni.client.GetCpuUtilisation(host)
		if err != nil {
			log.Errorf("[NODE-EXPORTER] Error: %v", err)
		}

		ram, err := ni.client.GetRamUtilisation(host)
		if err != nil {
			log.Errorf("[NODE-EXPORTER] Error: %v", err)
		}

		cpuFloat, _ := strconv.ParseFloat(cpu, 64)
		ramFloat, _ := strconv.ParseFloat(ram, 64)
		node := newNode(host, key, cpuFloat, ramFloat)
		ni.nodes = append(ni.nodes, node)
	}
	go ni.updateNodesInfo()
}

func (ni *NodesInfo) updateNodesInfo() {
	//log.Info("[NODE-EXPORTER] Updating NodesInfo...")
	ni.client.Time = time.Now()

	for id, node := range ni.nodes {
		cpu, err := ni.client.GetCpuUtilisation(node.InternalName)
		if err != nil {
			log.Errorf("[NODE-EXPORTER] Error: %v", err)
		}

		ram, err := ni.client.GetRamUtilisation(node.InternalName)
		if err != nil {
			log.Errorf("[NODE-EXPORTER] Error: %v", err)
		}

		cpuFloat, _ := strconv.ParseFloat(cpu, 64)
		ramFloat, _ := strconv.ParseFloat(ram, 64)

		ni.nodes[id].SetCpuUtil(cpuFloat)
		ni.nodes[id].SetRamUtil(ramFloat)
	}

	//log.Infof("[NODE-EXPORTER] Current state: %v", ni.nodes)
	time.Sleep(5 * time.Second)

	ni.updateNodesInfo()
}

func (ni *NodesInfo) GetNodeCpu(nodeName string) float64 {
	for _, node := range ni.nodes {
		if nodeName == node.ExternalName {
			return node.GetCpuUtil()
		}
	}
	return -1
}

func (ni *NodesInfo) GetNodeRam(nodeName string) float64 {
	for _, node := range ni.nodes {
		if nodeName == node.ExternalName {
			return node.GetRamUtil()
		}
	}
	return -1
}
