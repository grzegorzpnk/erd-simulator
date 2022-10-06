package metrics

import (
	log "10.254.188.33/matyspi5/erd/pkg/nmt/src/logger"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

/*type ClusterResources struct {
	//Latency float64    `json:"latency"`
	CpuUsage    ClusterResources `json:"cpu"`
	MemoryUsage ClusterResources `json:"memory"`
}*/

type ClusterResources struct {
	Used        float64 `json:"used"`        // How many cpu/memory used (value)
	Allocatable float64 `json:"allocatable"` // How many cpu/memory available (value)
	Utilization float64 `json:"utilization"` // How much is cpu/memory utilized (percentage)
}

type NetworkMetrics struct {
	//ms
	Latency float64 `json:"latency"`
}

func (cm *ClusterResources) UpdateClusterMetrics(clusterResources ClusterResources) {
	cm.Utilization = clusterResources.Utilization
	cm.Used = clusterResources.Used
	cm.Allocatable = clusterResources.Allocatable
}

func (nm *NetworkMetrics) UpdateNetworkMetrics(networkMetrics NetworkMetrics) {
	nm.Latency = networkMetrics.Latency
}

/*
func ClusterMetricsUpdateTest(g *Graph) {

	endpoint := config.GetConfiguration().ClusterControllerEndpoint

	for {
		// update metrics for MEC Clusters

		_, err := getClusterMetrics("mec1", endpoint)
		if err != nil {
			log.Errorf(err.Error())
		}

		time.Sleep(5 * time.Second)

	}
}
*/

/*
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
}*/

/*
func NetworkMetricsUpdateTest(g *Graph) {

	endpoint := config.GetConfiguration().ClusterControllerEndpoint
	latencyURL := buildLatencyURL(endpoint, "1", "1")

	for {
		// update metrics for MEC Clusters

		getNetworkMetricsNotification(latencyURL, &Edge{Source: "mec1", Target: "1"}, g)
		time.Sleep(5 * time.Second)

	}
}
*/

func GetClusterMetrics(clusterName, clusterProvider, endpoint string) (ClusterResources, error) {

	cr := httpGet(endpoint)
	log.Infof("Resp: %v", cr)

	return cr, nil
}

/*
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
	//	log.Infof("Latency: %v\n", latencyStr)
	nm.Latency, _ = strconv.ParseFloat(latencyStr, 32)

	return nm
}
*/
func httpGet(endpoint string) ClusterResources {

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
	var cr ClusterResources
	json.Unmarshal(body, &cr)

	/*sb := string(body)

	//delete whitespaces
	ret := strings.TrimSpace(sb)*/

	return cr

}

/*
func buildLatencyURL(endpoint, cellID, MECID string) string {
	//http://10.254.185.50:32138/v1/obs/ltc/cell/1/mec/edge-provider+mec1/latency-ms

	var latencyURL string
	//todo
	baseURL := endpoint + "/v1/obs/ltc/cell/"
	cellURL := baseURL + cellID + "/mec/"
	mecURL := cellURL + config.GetConfiguration().EdgeProvider + "+" + MECID
	latencyURL = mecURL + "/latency-ms"
	log.Infof("latency url: ", latencyURL)

	return latencyURL
}*/

func BuildCpuUrl(clusterName, clusterProvider, endpoint string) string {

	//http://10.254.185.50:32138/v1/obs/ksm/provider/{provider}/cluster/{cluster}/cpu

	baseURL := endpoint + "/v1/obs/ksm/provider/"
	providerURL := baseURL + clusterProvider + "/"
	baseClusterURL := providerURL + "cluster/"
	clusterURL := baseClusterURL + clusterName + "/"
	clusterCPUURL := clusterURL + "cpu"
	//log.Infof("cpu url: ", clusterCPUURL)

	return clusterCPUURL
}

func BuildMemoryUrl(clusterName, clusterProvider, endpoint string) string {

	////  http://10.254.185.50:32138/v1/obs/ksm/provider/{provider}/cluster/{cluster}/memory
	baseURL := endpoint + "/v1/obs/ksm/provider/"
	providerURL := baseURL + clusterProvider + "/"
	baseClusterURL := providerURL + "cluster/"
	clusterURL := baseClusterURL + clusterName + "/"
	clusterMemoryURL := clusterURL + "memory"
	//log.Infof("memory url: ", clusterMemoryURL)

	return clusterMemoryURL
}
