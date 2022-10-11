package metrics

import (
	log "10.254.188.33/matyspi5/erd/pkg/nmt/src/logger"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
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

/*//gorutine function
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
*/
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

func GetClusterMetrics(endpoint string) (ClusterResources, error) {

	cr := httpGet(endpoint)
	//	log.Infof("Resp: %v", cr)

	return cr, nil
}

func GetLatency(endpoint string) (float64, error) {

	lat, _ := httpGetLatency(endpoint)
	//log.Infof("Resp: %v", lat)

	return lat, nil
}

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

	return cr

}

func httpGetLatency(endpoint string) (float64, error) {

	resp, err := http.Get(endpoint)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)
	ret := strings.TrimSpace(sb)
	tmp, _ := strconv.ParseFloat(ret, 64)

	return tmp, nil
}

func Float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func BuildCellLatencyURL(endpoint, CellId, MecName, MecProvider string) string {
	//http://10.254.185.50:32138/v1/obs/ltc/cell/1/mec/edge-provider+mec1/latency-ms
	//http://10.254.185.50:32138/v1/obs/ltc/source/edge-provider+mec1/target/edge-provider+mec2/latency-ms

	var latencyURL string
	baseURL := endpoint + "/v1/obs/ltc/source/"
	cellURL := baseURL + CellId + "/target/"
	mecURL := cellURL + MecProvider + "+" + MecName
	latencyURL = mecURL + "/latency-ms"
	log.Infof("latency url: ", latencyURL)

	return latencyURL
}

func BuildMECLatencyURL(endpoint, targetMEC, targetProvider, SourceMEC, SourceProvider string) string {
	//http://10.254.185.50:32138/v1/obs/ltc/cell/1/mec/edge-provider+mec1/latency-ms
	//http://10.254.185.50:32138/v1/obs/ltc/source/edge-provider+mec1/target/edge-provider+mec2/latency-ms

	var latencyURL string
	baseURL := endpoint + "/v1/obs/ltc/source/"
	cellURL := baseURL + SourceProvider + "+" + SourceMEC + "/target/"
	mecURL := cellURL + targetProvider + "+" + targetMEC
	latencyURL = mecURL + "/latency-ms"
	log.Infof("latency url: ", latencyURL)

	return latencyURL
}

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
