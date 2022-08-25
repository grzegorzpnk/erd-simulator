package observability

import (
	"10.254.188.33/matyspi5/obs/src/config"
	log "10.254.188.33/matyspi5/obs/src/logger"
	"10.254.188.33/matyspi5/obs/src/pkg/promql"
	"errors"
	"fmt"
	"math"
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
	ci.client.Time = time.Now()

	for id, cluster := range ci.clusters {
		cpuReq, cpuLim, err := ci.client.GetCpuRequestsLimits(cluster.Name)
		if err != nil {
			log.Errorf("[KSM] error: could not updateClustersInfo. reason: %v", err)
		}
		memReq, memLim, err := ci.client.GetMemoryRequestsLimits(cluster.Name)
		if err != nil {
			log.Errorf("[KSM] error: could not updateClustersInfo. reason: %v", err)
		}

		ci.clusters[id].SetCpuReq(math.Round(cpuReq*100) / 100)
		ci.clusters[id].SetCpuLim(math.Round(cpuLim*100) / 100)
		ci.clusters[id].SetMemReq(math.Round(memReq*100) / 100)
		ci.clusters[id].SetMemLim(math.Round(memLim*100) / 100)
	}

	time.Sleep(5 * time.Second)
	ci.updateClustersInfo()
}

func (ci *ClustersInfo) GetClusterCpuReq(clusterProvider, clusterName string) (float64, error) {
	cluster, err := ci.GetCluster(clusterProvider, clusterName)
	if err != nil {
		err = fmt.Errorf("can't get cluster cpu requests. reason: %v", err)
		log.Errorf("[KSM] error: %v", err)
		return -1, err
	}
	return cluster.GetMemReq(), nil
}

func (ci *ClustersInfo) GetClusterCpuLim(clusterProvider, clusterName string) (float64, error) {
	cluster, err := ci.GetCluster(clusterProvider, clusterName)
	if err != nil {
		err = fmt.Errorf("can't get cluster cpu limits. reason: %v", err)
		log.Errorf("[KSM] error: %v", err)
		return -1, err
	}
	return cluster.GetCpuLim(), nil
}

func (ci *ClustersInfo) GetClusterMemReq(clusterProvider, clusterName string) (float64, error) {
	cluster, err := ci.GetCluster(clusterProvider, clusterName)
	if err != nil {
		err = fmt.Errorf("can't get cluster memory requests. reason: %v", err)
		log.Errorf("[KSM] error: %v", err)
		return -1, err
	}
	return cluster.GetMemReq(), nil
}

func (ci *ClustersInfo) GetClusterMemLim(clusterProvider, clusterName string) (float64, error) {
	cluster, err := ci.GetCluster(clusterProvider, clusterName)
	if err != nil {
		err = fmt.Errorf("can't get cluster memory limits. reason: %v", err)
		log.Errorf("[KSM] error: %v", err)
		return -1, err
	}
	return cluster.GetMemLim(), nil
}

func (ci *ClustersInfo) GetCluster(clusterProvider, clusterName string) (Cluster, error) {
	for _, cluster := range ci.clusters {
		if cluster.Provider == clusterProvider && cluster.Name == clusterName {
			return cluster, nil
		}
	}
	err := errors.New(fmt.Sprintf("cluster %v+%v doesn't exist", clusterProvider, clusterName))
	return Cluster{}, err
}
