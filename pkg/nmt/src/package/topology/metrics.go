package topology

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"nmt/src/config"
	"strconv"
	"time"
)

type ClusterMetrics struct {
	CpuUsage    int `json:"cpuUsage"`
	MemoryUsage int `json:"memoryUsage"`
	RamUsage    int `json:"ramUsage"`
}

type NetworkMetrics struct {

	//ms
	Latency    float64 `json:"latency"`
	PacketDrop int     `json:"packetDrop"`
}

func (cm *ClusterMetrics) UpdateClusterMetrics(clusterMetrics ClusterMetrics) {

	cm.CpuUsage = clusterMetrics.CpuUsage
	cm.RamUsage = clusterMetrics.RamUsage
	cm.MemoryUsage = clusterMetrics.MemoryUsage
}

func (nm *NetworkMetrics) UpdateNetworkMetrics(networkMetrics NetworkMetrics) {

	nm.Latency = networkMetrics.Latency
}

//gorutine function
func TopologyMetricsUpdate(g *Graph) {

	endpoint := config.GetConfiguration().ClusterControllerEndpoint

	for {
		// update metrics for MEC Clusters
		for i, v := range g.Vertices {
			if v.Type == "MEC" {
				cm, err := getClusterMetricsNotification(strconv.Itoa(v.Id), endpoint)
				if err != nil {
					fmt.Errorf(err.Error())
				}
				g.Vertices[i].VertexMetrics.UpdateClusterMetrics(cm)
			}
		}
		/*
			//Update metrics for each EDGE
			for j, k := range g.Edges {
				nm := getNetworkMetricsNotification(endpoint, k, g)
				g.Edges[j].EdgeMetrics.UpdateNetworkMetrics(nm)
			}*/
	}

	time.Sleep(1 * time.Second)
}

func TopologyMetricsUpdateTest(g *Graph) {

	endpoint := config.GetConfiguration().ClusterControllerEndpoint

	for {
		// update metrics for MEC Clusters

		_, err := getClusterMetricsNotification(strconv.Itoa(1), endpoint)
		if err != nil {
			fmt.Errorf(err.Error())
		}

		time.Sleep(10 * time.Second)

	}
}

//this function takes clusterID and requests to receive latest info about Cluster Metrics at the end it returns ClusterMetrics object
func getClusterMetricsNotification(clusterId, endpoint string) (ClusterMetrics, error) {

	var cm ClusterMetrics

	clusterCPUURL := buildCpuUrl(clusterId, endpoint)
	clusterMemoryURL := buildMemoryUrl(clusterId, endpoint)

	//get current CPU
	CPUresp, err := http.Get(clusterCPUURL)
	if err != nil {
		fmt.Errorf("Cannot get CPU: %v", err)
	}
	defer CPUresp.Body.Close()

	//We Read the response body on the line below.
	CPUbody, err := ioutil.ReadAll(CPUresp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(CPUbody)
	log.Printf(sb)
	fmt.Printf("CPU Response status:%v , CPU: %v\n", CPUresp.Status, sb)

	//get current Memory
	MemoryResp, err := http.Get(clusterMemoryURL)
	if err != nil {
		panic(err)
	}
	defer MemoryResp.Body.Close()

	//We Read the response body on the line below.
	Memorybody, err := ioutil.ReadAll(MemoryResp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb2 := string(Memorybody)
	log.Printf(sb2)
	fmt.Printf("Memory Response status:%v , Memory: %v\n", MemoryResp.Status, sb2)

	return cm, nil
}

func getNetworkMetricsNotification(endpoint string, edge *Edge, g *Graph) NetworkMetrics {
	var nm NetworkMetrics
	var cellID, mecID string

	if g.GetVertex(edge.Source).Type == "CELL" {
		cellID = strconv.Itoa(edge.Source)
		mecID = strconv.Itoa(edge.Target)
	} else {
		cellID = strconv.Itoa(edge.Target)
		mecID = strconv.Itoa(edge.Source)
	}
	latencyURL := buildLatencyURL(endpoint, cellID, mecID)

	//get current latency
	latencyResp, err := http.Get(latencyURL)
	if err != nil {
		panic(err)
	}
	defer latencyResp.Body.Close()

	fmt.Println("Latency Response status:", latencyResp.Status)

	CPUscanner := bufio.NewScanner(latencyResp.Body)
	nm.Latency, _ = strconv.ParseFloat(CPUscanner.Text(), 20)

	return nm
}

//http://10.254.185.50:32138/v1/obs/ltc/cell/1/meh/edge-provider+meh01/latency-ms
func buildLatencyURL(endpoint, cellID, MECID string) string {

	var latencyURL string
	/*baseURL := endpoint + "v1/obs/ksm/provider/edge-provider/"
	providerURL := baseURL + "orange/"
	baseClusterURL := providerURL + "cluster/meh"
	clusterURL := baseClusterURL + id + "/"
	clusterCPUURL := clusterURL + "cpu-requests"
	fmt.Println("cpu url: ", clusterCPUURL)*/

	return latencyURL
}

//http://10.254.185.50:32138/v1/obs/ksm/provider/edge-provider/cluster/meh01/cpu-requests
func buildCpuUrl(id, endpoint string) string {

	baseURL := endpoint + "v1/obs/ksm/provider/"
	providerURL := baseURL + config.GetConfiguration().EdgeProvider + "/"
	baseClusterURL := providerURL + "cluster/meh0"
	clusterURL := baseClusterURL + id + "/"
	clusterCPUURL := clusterURL + "cpu-requests"
	fmt.Println("cpu url: ", clusterCPUURL)

	return clusterCPUURL
}

//np. http://10.254.185.50:32138/v1/obs/ksm/provider/edge-provider/cluster/meh02/memory-requests
func buildMemoryUrl(id, endpoint string) string {

	baseURL := endpoint + "v1/obs/ksm/provider/"
	providerURL := baseURL + config.GetConfiguration().EdgeProvider + "/"
	baseClusterURL := providerURL + "cluster/meh0"
	clusterURL := baseClusterURL + id + "/"
	clusterMemoryURL := clusterURL + "memory-requests"
	fmt.Println("memory url: ", clusterMemoryURL)

	return clusterMemoryURL
}
