package metrics

import (
	log "10.254.188.33/matyspi5/erd/pkg/nmt/src/logger"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type ClusterResources struct {
	Used        float64 `json:"used"`        // How many cpu/memory used (value)
	Capacity    float64 `json:"capacity"`    // How many cpu/memory available (value)
	Utilization float64 `json:"utilization"` // How much is cpu/memory utilized (percentage)
}

type NetworkMetrics struct {
	//ms
	Latency float64 `json:"latency"`
}

func (cm *ClusterResources) UpdateClusterMetrics(clusterResources ClusterResources) {
	cm.Utilization = clusterResources.Utilization
	cm.Used = clusterResources.Used
	cm.Capacity = clusterResources.Capacity
}

func (nm *NetworkMetrics) UpdateNetworkMetrics(networkMetrics NetworkMetrics) {
	nm.Latency = networkMetrics.Latency
}

func GetLatency(endpoint string) (float64, error) {

	lat, _ := httpGetLatency(endpoint)
	//log.Infof("Resp: %v", lat)

	return lat, nil
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

func BuildCellLatencyURL(endpoint, CellId, MecName, MecProvider string) string {
	//http://10.254.185.50:32138/v1/obs/ltc/cell/1/mec/edge-provider+mec1/latency-ms
	//http://10.254.185.50:32138/v1/obs/ltc/source/edge-provider+mec1/target/edge-provider+mec2/latency-ms

	var latencyURL string
	baseURL := endpoint + "/v1/obs/ltc/source/"
	cellURL := baseURL + CellId + "/target/"
	mecURL := cellURL + MecProvider + "+" + MecName
	latencyURL = mecURL + "/latency-ms"
	//log.Infof("latency url: ", latencyURL)

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
	//log.Infof("latency url: ", latencyURL)

	return latencyURL
}
