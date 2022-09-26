package topology

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"nmt/src/config"
	"strconv"
	"strings"
	"time"
)

type ClusterMetrics struct {
	CpuUsage    float64 `json:"cpuUsage"`
	MemoryUsage float64 `json:"memoryUsage"`
}

type MecResources struct {
	//Latency float64    `json:"latency"`
	Cpu    MecResInfo `json:"cpu"`
	Memory MecResInfo `json:"memory"`
}

type NetworkMetrics struct {

	//ms
	Latency float64 `json:"latency"`
}

func (cm *ClusterMetrics) UpdateClusterMetrics(clusterMetrics ClusterMetrics) {

	cm.CpuUsage = clusterMetrics.CpuUsage
	cm.MemoryUsage = clusterMetrics.MemoryUsage
}

func (nm *NetworkMetrics) UpdateNetworkMetrics(networkMetrics NetworkMetrics) {

	nm.Latency = networkMetrics.Latency
}

//gorutine function
func ClustersMetricsUpdate(g *Graph) {

	endpoint := config.GetConfiguration().ClusterControllerEndpoint

	for {
		// update metrics for MEC Clusters
		for _, v := range g.MecHosts {
			if v.Type == "MEC" {
				cm, err := getClusterMetricsNotification(v.Id, endpoint)
				if err != nil {
					fmt.Errorf(err.Error())
				}
				g.GetVertex(v.Id).VertexMetrics.UpdateClusterMetrics(cm)
				//g.MecHosts[i].VertexMetrics.UpdateClusterMetrics(cm)
				fmt.Printf("Controller updates cluster metrics for vertex ID: %v\n", v.Id)
			}
		}
		time.Sleep(1 * time.Second)
	}

}

func ClusterMetricsUpdateTest(g *Graph) {

	endpoint := config.GetConfiguration().ClusterControllerEndpoint

	for {
		// update metrics for MEC Clusters

		_, err := getClusterMetricsNotification("mec1", endpoint)
		if err != nil {
			fmt.Errorf(err.Error())
		}

		time.Sleep(5 * time.Second)

	}
}

//gorutine function
func NetworkMetricsUpdate(g *Graph) {

	endpoint := config.GetConfiguration().ClusterControllerEndpoint

	for {
		//Update metrics for each EDGE
		for j, k := range g.Edges {
			nm := getNetworkMetricsNotification(endpoint, k, g)
			g.Edges[j].EdgeMetrics.UpdateNetworkMetrics(nm)
		}
		time.Sleep(1 * time.Second)
	}
}

func NetworkMetricsUpdateTest(g *Graph) {

	endpoint := config.GetConfiguration().ClusterControllerEndpoint
	latencyURL := buildLatencyURL(endpoint, "1", "1")

	for {
		// update metrics for MEC Clusters

		getNetworkMetricsNotification(latencyURL, &Edge{Source: "mec1", Target: "1"}, g)
		time.Sleep(5 * time.Second)

	}
}

//this function takes clusterID and requests to receive latest info about ClusterName Metrics at the end it returns ClusterMetrics object
func getClusterMetricsNotification(clusterId, endpoint string) (ClusterMetrics, error) {

	var cm ClusterMetrics

	//CPU----------------------------------------------------------------------------
	clusterCPUURL := buildCpuUrl(clusterId, endpoint)
	clusterMemoryURL := buildMemoryUrl(clusterId, endpoint)

	//get current CPU
	cpuStr := httpGet(clusterCPUURL)

	log.Printf("CPU resp: %v", cpuStr)
	//fmt.Printf("CPU Response :%v\n", cpuStr)
	cm.CpuUsage, _ = strconv.ParseFloat(cpuStr, 32)

	//MEMORY-------------------------------------------------------------------------
	MemoryStr := httpGet(clusterMemoryURL)
	log.Printf("Memory resp: %v", MemoryStr)
	//fmt.Printf("Memory: %v\n", MemoryStr)
	cm.MemoryUsage, _ = strconv.ParseFloat(MemoryStr, 32)

	return cm, nil
}

//to be tested
func getNetworkMetricsNotification(endpoint string, edge *Edge, g *Graph) NetworkMetrics {
	var nm NetworkMetrics
	var cellID, mecID string

	if g.GetVertex(edge.Source).Type == "CELL" {
		cellID = edge.Source
		mecID = edge.Target
	} else {
		cellID = edge.Target
		mecID = edge.Source
	}
	latencyURL := buildLatencyURL(endpoint, cellID, mecID)
	latencyStr := httpGet(latencyURL)

	//get current latency
	log.Printf("Latency resp: %v ", latencyStr)
	//	fmt.Printf("Latency: %v\n", latencyStr)
	nm.Latency, _ = strconv.ParseFloat(latencyStr, 32)

	return nm
}

func httpGet(endpoint string) string {

	resp, err := http.Get(endpoint)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)

	//delete whitespaces
	ret := strings.TrimSpace(sb)

	return ret

}

func buildLatencyURL(endpoint, cellID, MECID string) string {
	//http://10.254.185.50:32138/v1/obs/ltc/cell/1/mec/edge-provider+mec1/latency-ms

	var latencyURL string
	//todo
	baseURL := endpoint + "/v1/obs/ltc/cell/"
	cellURL := baseURL + cellID + "/mec/"
	mecURL := cellURL + config.GetConfiguration().EdgeProvider + "+" + MECID
	latencyURL = mecURL + "/latency-ms"
	fmt.Println("latency url: ", latencyURL)

	return latencyURL
}

func buildCpuUrl(mecId, endpoint string) string {
	//http://10.254.185.50:32138/v1/obs/ksm/provider/edge-provider/cluster/mec1/cpu-requests

	baseURL := endpoint + "/v1/obs/ksm/provider/"
	providerURL := baseURL + config.GetConfiguration().EdgeProvider + "/"
	baseClusterURL := providerURL + "cluster/"
	clusterURL := baseClusterURL + mecId + "/"
	clusterCPUURL := clusterURL + "cpu-requests"
	//fmt.Println("cpu url: ", clusterCPUURL)

	return clusterCPUURL
}

func buildMemoryUrl(id, endpoint string) string {
	//np. http://10.254.185.50:32138/v1/obs/ksm/provider/edge-provider/cluster/mec2/memory-requests

	baseURL := endpoint + "/v1/obs/ksm/provider/"
	providerURL := baseURL + config.GetConfiguration().EdgeProvider + "/"
	baseClusterURL := providerURL + "cluster/"
	clusterURL := baseClusterURL + id + "/"
	clusterMemoryURL := clusterURL + "memory-requests"
	//fmt.Println("memory url: ", clusterMemoryURL)

	return clusterMemoryURL
}
