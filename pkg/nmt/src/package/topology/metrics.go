package topology

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

	cm.CpuUsage = clusterMetrics.CpuUsage
	cm.RamUsage = clusterMetrics.RamUsage
	cm.MemoryUsage = clusterMetrics.MemoryUsage

}
