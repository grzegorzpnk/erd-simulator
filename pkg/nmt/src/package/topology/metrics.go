package topology

import (
	"bufio"
	"fmt"
	"net/http"
	"nmt/src/config"
	"time"
)

type ClusterMetrics struct {
	CpuUsage    int `json:"cpuUsage"`
	MemoryUsage int `json:"memoryUsage"`
	RamUsage    int `json:"ramUsage"`
}

type NetworkMetrics struct {

	//ms
	Latency    float32 `json:"latency"`
	PacketDrop int     `json:"packetDrop"`
}

func (cm *ClusterMetrics) UpdateClusterMetrics(clusterMetrics ClusterMetrics) {

	/*cm.CpuUsage = clusterMetrics.CpuUsage
	cm.RamUsage = clusterMetrics.RamUsage
	cm.MemoryUsage = clusterMetrics.MemoryUsage*/
}

func (nm *NetworkMetrics) UpdateNetworkMetrics(networkMetrics NetworkMetrics) {

	nm.Latency = networkMetrics.Latency
}

func (graph *Graph) TopologyMetricsUpdate() {

	endpoint := config.GetConfiguration().ClusterControllerEndpoint

	for {

		for i, v := range graph.Vertices {

			cm, err := getClusterMetricsNotification(v.Id, endpoint)
			if err != nil {
				fmt.Errorf(err.Error())
			}
			graph.Vertices[i].VertexMetrics.UpdateClusterMetrics(cm)

		}
		time.Sleep(1 * time.Second)
	}

}

func getClusterMetricsNotification(id int, endpoint string) (ClusterMetrics, error) {

	var cm ClusterMetrics

	baseURL := endpoint + "v1/obs/ksm/provider/edge-provider/"
	providerURL := baseURL + "orange/"
	baseClusterURL := providerURL + "cluster/meh"
	clusterURL := baseClusterURL + string(id) + "/"
	clusterCPUURL := clusterURL + "cpu-requests"
	fmt.Println(clusterCPUURL)

	resp, err := http.Get(clusterCPUURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	/*//func getNotification(sub db.Subscriber, amfEndpoint string) (types.AmfCreatedEventSubscription, error) {
		var respBody types.AmfCreatedEventSubscription
		reqBody, err := json.Marshal(sub.BodyRequest)

		resp, err := http.Post(amfEndpoint, "text/plain", bytes.NewBuffer(reqBody))
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("Could not get notification for: EventType: %s, AMF endpoint: %s",
				sub.AmfEventType, sub.Endpoint))
			return types.AmfCreatedEventSubscription{}, err
		}

		body, err := ioutil.ReadAll(resp.Body)

		err = json.Unmarshal(body, &respBody)
		if err != nil {
			err = errors.Wrap(err, "Failed to unmarshal body")
			log.Errorf("[SUBSCRIPTION][ID=%v] Error: %v", err)
			return types.AmfCreatedEventSubscription{}, err
		}
		return respBody, nil
	}*/

	return cm, nil
}
