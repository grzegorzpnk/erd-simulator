package topology

import (
	"bufio"
	"fmt"
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

//currently this function is hardcoded, please make a common topology to make it work
func TopologyMetricsUpdate(g *Graph) {

	endpoint := config.GetConfiguration().ClusterControllerEndpoint

	for {
		//update cluster metrics
		/*		//for i, v := range g.Vertices {
				//todo: hardcoded to test please change for generic
				cm, err := getClusterMetricsNotification("02", endpoint)
				if err != nil {
					fmt.Errorf(err.Error())
				}
				g.Vertices[2].VertexMetrics.UpdateClusterMetrics(cm)

				//}*/

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

		//Update metrics for each EDGE
		for j, k := range g.Edges {
			nm := getNetworkMetricsNotification(endpoint, k)

			g.Edges[j].EdgeMetrics.UpdateNetworkMetrics(nm)
		}
	}

	time.Sleep(1 * time.Second)
}

//this function takes clusterID and requests to receive latest info about Cluster Metrics at the end it returns ClusterMetrics object
func getClusterMetricsNotification(clusterId, endpoint string) (ClusterMetrics, error) {

	var cm ClusterMetrics

	clusterCPUURL := buildCpuUrl(clusterId, endpoint)
	clusterMemoryURL := buildMemoryUrl(clusterId, endpoint)

	//get current CPU
	CPUresp, err := http.Get(clusterCPUURL)
	if err != nil {
		panic(err)
	}
	defer CPUresp.Body.Close()

	fmt.Println("CPU Response status:", CPUresp.Status)

	CPUscanner := bufio.NewScanner(CPUresp.Body)
	cm.CpuUsage, _ = strconv.Atoi(CPUscanner.Text())

	//get current Memory
	MemoryResp, err := http.Get(clusterMemoryURL)
	if err != nil {
		panic(err)
	}
	defer MemoryResp.Body.Close()

	fmt.Println("Memory Response status:", MemoryResp.Status)

	MemoryScanner := bufio.NewScanner(MemoryResp.Body)
	cm.MemoryUsage, _ = strconv.Atoi(MemoryScanner.Text())

	return cm, nil
}

// np. http://10.254.185.50:32138/v1/obs/ltc/cell/1/meh/edge-provider+meh01/latency-ms
func getNetworkMetricsNotification(endpoint string, edge *Edge) NetworkMetrics {
	var nm NetworkMetrics

	//if Vertex[edge.Source].Type == "CELL"
	cellID := strconv.Itoa(edge.Source)
	mecID := strconv.Itoa(edge.Target)
	//else
	//mecID := edge.Source
	//cellID := edge.target
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

func buildCpuUrl(id, endpoint string) string {

	baseURL := endpoint + "v1/obs/ksm/provider/edge-provider/"
	providerURL := baseURL + "orange/"
	baseClusterURL := providerURL + "cluster/meh"
	clusterURL := baseClusterURL + id + "/"
	clusterCPUURL := clusterURL + "cpu-requests"
	fmt.Println("cpu url: ", clusterCPUURL)

	return clusterCPUURL
}

//np. http://10.254.185.50:32138/v1/obs/ksm/provider/edge-provider/cluster/meh02/memory-requests
func buildMemoryUrl(id, endpoint string) string {

	baseURL := endpoint + "v1/obs/ksm/provider/edge-provider/"
	providerURL := baseURL + "orange/"
	baseClusterURL := providerURL + "cluster/meh"
	clusterURL := baseClusterURL + id + "/"
	clusterMemoryURL := clusterURL + "memory-requests"
	fmt.Println("memory url: ", clusterMemoryURL)

	return clusterMemoryURL
}
