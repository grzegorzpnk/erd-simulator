package observability

import (
	"10.254.188.33/matyspi5/pmc/src/config"
	log "10.254.188.33/matyspi5/pmc/src/logger"
	"10.254.188.33/matyspi5/pmc/src/pkg/promql"
	"time"
)

type Cluster struct {
	Provider    string
	Name        string
	CpuRequests float64
	CpuLimits   float64
	MemRequests float64
	MemLimits   float64
}

// NewCluster method is called to register clusters from configuration file.
func NewCluster(provider, name string, cpuReq, cpuLim, ramReq, ramLim float64) Cluster {
	return Cluster{
		Provider:    provider,
		Name:        name,
		CpuRequests: cpuReq,
		CpuLimits:   cpuLim,
		MemRequests: ramReq,
		MemLimits:   ramLim,
	}
}

// Setters

func (c *Cluster) SetCpuReq(val float64) {
	c.CpuRequests = val
}

func (c *Cluster) SetCpuLim(val float64) {
	c.CpuLimits = val
}

func (c *Cluster) SetMemReq(val float64) {
	c.MemRequests = val
}

func (c *Cluster) SetMemLim(val float64) {
	c.MemLimits = val
}

// Getters

func (c *Cluster) GetCpuReq() float64 {
	return c.CpuRequests
}

func (c *Cluster) GetCpuLim() float64 {
	return c.CpuLimits
}

func (c *Cluster) GetMemReq() float64 {
	return c.MemRequests
}

func (c *Cluster) GetMemLim() float64 {
	return c.MemLimits
}

type ClustersInfo struct {
	client   promql.PromQL
	clusters []Cluster
}

// InitializeClustersInfo reads providers and clusters from config file. Then fetches cpu/memory requests/limits
// from Prometheus data-source and saves to the ClusterInfo. This information are updated in separate goroutine.
func (ci *ClustersInfo) InitializeClustersInfo(client promql.PromQL) {
	log.Info("[KSM] Initializing ClustersInfo...")
	ci.client = client

	ClusterSets := config.GetConfiguration().Clusters

	for _, clusterSet := range ClusterSets {
		for _, cluster := range clusterSet.Clusters {
			cpuReq, cpuLim, err := ci.client.GetCpuRequestsLimits(cluster)
			if err != nil {
				log.Errorf("Skipping. Error: %v", err)
				continue
			}
			ramReq, ramLim, err := ci.client.GetMemoryRequestsLimits(cluster)
			if err != nil {
				log.Errorf("Skipping. Error: %v", err)
				continue
			}

			cl := NewCluster(clusterSet.Provider, cluster, cpuReq, cpuLim, ramReq, ramLim)
			ci.clusters = append(ci.clusters, cl)
		}
	}

	go ci.updateClustersInfo()
}

// updateClustersInfo fetches current cpu/memory request/limits utilization and updates ClusterInfo
func (ci *ClustersInfo) updateClustersInfo() {
	//log.Info("[KSM] Updating ClustersInfo...")
	ci.client.Time = time.Now()

	for id, cluster := range ci.clusters {
		cpuReq, cpuLim, err := ci.client.GetCpuRequestsLimits(cluster.Name)
		if err != nil {
			log.Errorf("[KSM] Could not updateClustersInfo. Reason: %v", err)
		}
		memReq, memLim, err := ci.client.GetMemoryRequestsLimits(cluster.Name)
		if err != nil {
			log.Errorf("[KSM] Could not updateClustersInfo. Reason: %v", err)
		}

		ci.clusters[id].SetCpuReq(cpuReq)
		ci.clusters[id].SetCpuLim(cpuLim)
		ci.clusters[id].SetMemReq(memReq)
		ci.clusters[id].SetMemLim(memLim)
	}

	//log.Infof("[KSM] Current state: %v", ci.clusters)
	time.Sleep(5 * time.Second)

	ci.updateClustersInfo()
}

func (ci *ClustersInfo) GetClusterCpuReq(clusterProvider, clusterName string) float64 {
	for _, cluster := range ci.clusters {
		if clusterName == cluster.Name && clusterProvider == cluster.Provider {
			return cluster.GetCpuReq()
		}
	}
	return -1
}

func (ci *ClustersInfo) GetClusterCpuLim(clusterProvider, clusterName string) float64 {
	for _, cluster := range ci.clusters {
		if clusterName == cluster.Name && clusterProvider == cluster.Provider {
			return cluster.GetCpuLim()
		}
	}
	return -1
}

func (ci *ClustersInfo) GetClusterMemReq(clusterProvider, clusterName string) float64 {
	for _, cluster := range ci.clusters {
		if clusterName == cluster.Name && clusterProvider == cluster.Provider {
			return cluster.GetMemReq()
		}
	}
	return -1
}

func (ci *ClustersInfo) GetClusterMemLim(clusterProvider, clusterName string) float64 {
	for _, cluster := range ci.clusters {
		if clusterName == cluster.Name && clusterProvider == cluster.Provider {
			return cluster.GetMemLim()
		}
	}
	return -1
}
