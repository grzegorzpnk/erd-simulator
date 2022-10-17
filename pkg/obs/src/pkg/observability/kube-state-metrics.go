package observability

import (
	"10.254.188.33/matyspi5/erd/pkg/obs/src/config"
	log "10.254.188.33/matyspi5/erd/pkg/obs/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/obs/src/pkg/promql"
	"errors"
	"fmt"
	"math"
	"time"
)

type Cluster struct {
	Provider       string
	Name           string
	CpuRequests    float64
	CpuLimits      float64
	CpuAllocatable float64
	MemRequests    float64
	MemLimits      float64
	MemAllocatable float64
	LastUpdate     time.Time
	Mocked         bool
}

// NewCluster method is called to register clusters from configuration file.
func NewCluster(provider, name string, mocked bool, lastUpdate time.Time, cpuReq, cpuLim, cpuAlloc, memReq, memLim, memAlloc float64) Cluster {
	return Cluster{
		Provider:       provider,
		Name:           name,
		CpuRequests:    cpuReq,
		CpuLimits:      cpuLim,
		CpuAllocatable: cpuAlloc,
		MemRequests:    memReq,
		MemLimits:      memLim,
		MemAllocatable: memAlloc,
		LastUpdate:     lastUpdate,
		Mocked:         mocked,
	}
}

func (c *Cluster) MockCluster() {
	c.Mocked = true
	c.CpuAllocatable = 1
	c.CpuRequests = 1
	c.CpuLimits = 1
	c.MemAllocatable = 1
	c.MemRequests = 1
	c.MemLimits = 1
}

// Setters

func (c *Cluster) SetCpuReq(val float64) {
	c.CpuRequests = math.Round(val*100) / 100
}

func (c *Cluster) SetCpuLim(val float64) {
	c.CpuLimits = math.Round(val*100) / 100
}

func (c *Cluster) SetCpuAlloc(val float64) {
	c.CpuAllocatable = math.Round(val*100) / 100
}

func (c *Cluster) SetMemReq(val float64) {
	c.MemRequests = math.Round(val*100) / 100
}

func (c *Cluster) SetMemLim(val float64) {
	c.MemLimits = math.Round(val*100) / 100
}

func (c *Cluster) SetMemAlloc(val float64) {
	c.MemAllocatable = math.Round(val*100) / 100
}

// Getters

func (c *Cluster) GetCpuReq() float64 {
	return math.Round(c.CpuRequests*100) / 100
}

func (c *Cluster) GetCpuLim() float64 {
	return math.Round(c.CpuLimits*100) / 100
}

func (c *Cluster) GetCpuAlloc() float64 {
	return math.Round(c.CpuAllocatable*100) / 100
}

func (c *Cluster) GetCpuReqUtilization() float64 {
	return math.Round(100*100*(c.CpuRequests/c.CpuAllocatable)) / 100
}

func (c *Cluster) GetMemReq() float64 {
	return math.Round(c.MemRequests*100) / 100
}

func (c *Cluster) GetMemLim() float64 {
	return math.Round(c.MemLimits*100) / 100
}

func (c *Cluster) GetMemAlloc() float64 {
	return math.Round(c.MemAllocatable*100) / 100
}

func (c *Cluster) GetMemReqUtilization() float64 {
	return math.Round(100*100*(c.MemRequests/c.MemAllocatable)) / 100
}

type ClustersInfo struct {
	client   promql.PromQL
	clusters []Cluster
}

// InitializeClustersInfo reads providers and clusters from config file. Then fetches cpu/memory requests/limits
// from Prometheus data-source and saves to the ClusterInfo. This information are updated in separate goroutine.
func (ci *ClustersInfo) InitializeClustersInfo(client promql.PromQL) {
	var cpuReq, cpuLim, cpuAll, memReq, memLim, memAll float64
	var err error

	log.Info("[KSM] Initializing ClustersInfo...")
	ci.client = client

	ClusterSet := config.ReadTopologyConfigFile("mecTopology.json")

	for _, cluster := range ClusterSet {
		if cpuReq, err = ci.client.GetCurrentRequests(promql.CPU, cluster.Identity.Cluster); err != nil {
			log.Warnf("Type[%v], Cluster[%v]. Could not get current requests. Reason: %v.", promql.CPU, cluster.BuildClusterEmcoFQDN(), err)
		}
		if cpuLim, err = ci.client.GetCurrentLimits(promql.CPU, cluster.Identity.Cluster); err != nil {
			log.Warnf("Type[%v], Cluster[%v]. Could not get current limits. Reason: %v.", promql.CPU, cluster.BuildClusterEmcoFQDN(), err)
		}
		if cpuAll, err = ci.client.GetAllocatable(promql.CPU, cluster.Identity.Cluster); err != nil {
			log.Warnf("Type[%v], Cluster[%v]. Could not get allocatable. Reason: %v.", promql.CPU, cluster.BuildClusterEmcoFQDN(), err)
		}
		if memReq, err = ci.client.GetCurrentRequests(promql.MEMORY, cluster.Identity.Cluster); err != nil {
			log.Warnf("Type[%v], Cluster[%v]. Could not get current requests. Reason: %v.", promql.MEMORY, cluster.BuildClusterEmcoFQDN(), err)
		}
		if memLim, err = ci.client.GetCurrentLimits(promql.MEMORY, cluster.Identity.Cluster); err != nil {
			log.Warnf("Type[%v], Cluster[%v]. Could not get current limits. Reason: %v.", promql.MEMORY, cluster.BuildClusterEmcoFQDN(), err)
		}
		if memAll, err = ci.client.GetAllocatable(promql.MEMORY, cluster.Identity.Cluster); err != nil {
			log.Warnf("Type[%v], Cluster[%v]. Could not get allocatable. Reason: %v.", promql.MEMORY, cluster.BuildClusterEmcoFQDN(), err)
		}

		if memAll > 0 && cpuAll > 0 {
			cl := NewCluster(cluster.Identity.Provider, cluster.Identity.Cluster, false, time.Now(), cpuReq, cpuLim, cpuAll, memReq, memLim, memAll)
			ci.clusters = append(ci.clusters, cl)
			log.Infof("Adding REAL cluster to the ClustersInfo [%v+%v]", cl.Provider, cl.Name)
		} else {
			cl := NewCluster(cluster.Identity.Provider, cluster.Identity.Cluster, true, time.Now(), 1, 1, 1, 1, 1, 1)
			ci.clusters = append(ci.clusters, cl)
			log.Infof("Adding MOCKED cluster to the ClustersInfo [%v+%v]", cl.Provider, cl.Name)
		}

	}

	go ci.updateClustersInfo()
}

func (ci *ClustersInfo) clusterUnMockable(id int) bool {
	var cpuReq, cpuLim, cpuAll, memReq, memLim, memAll float64
	var err error

	cl := ci.clusters[id]
	cluster := cl.Name
	if cpuReq, err = ci.client.GetCurrentRequests(promql.CPU, cluster); err != nil {
		return false
	}
	if cpuLim, err = ci.client.GetCurrentLimits(promql.CPU, cluster); err != nil {
		return false
	}
	if cpuAll, err = ci.client.GetAllocatable(promql.CPU, cluster); err != nil {
		return false
	}
	if memReq, err = ci.client.GetCurrentRequests(promql.MEMORY, cluster); err != nil {
		return false
	}
	if memLim, err = ci.client.GetCurrentLimits(promql.MEMORY, cluster); err != nil {
		return false
	}
	if memAll, err = ci.client.GetAllocatable(promql.MEMORY, cluster); err != nil {
		return false
	}

	ci.clusters[id].LastUpdate = time.Now()
	ci.clusters[id].Mocked = false
	ci.clusters[id].SetCpuReq(cpuReq)
	ci.clusters[id].SetCpuLim(cpuLim)
	ci.clusters[id].SetCpuAlloc(cpuAll)
	ci.clusters[id].SetMemReq(memReq)
	ci.clusters[id].SetMemLim(memLim)
	ci.clusters[id].SetMemAlloc(memAll)

	return true
}

// updateClustersInfo fetches current cpu/memory request/limits utilization and updates ClusterInfo
func (ci *ClustersInfo) updateClustersInfo() {
	var cpuReq, cpuLim, cpuAll, memReq, memLim, memAll float64
	var err error

	ci.client.Time = time.Now()

	for id, cl := range ci.clusters {
		if cl.Mocked && ci.client.Time.Sub(cl.LastUpdate) < 5*time.Minute {
			continue
		} else if cl.Mocked && ci.client.Time.Sub(cl.LastUpdate) > 5*time.Minute {
			if !ci.clusterUnMockable(id) {
				continue
			} else {
				log.Infof("Cluster[%v] reachable... Un-mocking")
				continue
			}
		}

		if ci.client.Time.Sub(cl.LastUpdate) > 2*time.Minute {
			log.Warnf("There are problems to update cluster[%v] information since[%v].. Mocking", cl.Name, ci.client.Time.Sub(cl.LastUpdate))
			ci.clusters[id].MockCluster()
			continue
		}

		cluster := cl.Name
		if cpuReq, err = ci.client.GetCurrentRequests(promql.CPU, cluster); err != nil {
			log.Warnf("Type[%v], Cluster[%v], Action[skip-new-values]. Could not get current requests. Reason: %v.", promql.CPU, cluster, err)
			continue
		}
		if cpuLim, err = ci.client.GetCurrentLimits(promql.CPU, cluster); err != nil {
			log.Warnf("Type[%v], Cluster[%v], Action[skip-new-values]. Could not get current limits. Reason: %v.", promql.CPU, cluster, err)
			continue
		}
		if cpuAll, err = ci.client.GetAllocatable(promql.CPU, cluster); err != nil {
			log.Warnf("Type[%v], Cluster[%v], Action[skip-new-values]. Could not get allocatable. Reason: %v.", promql.CPU, cluster, err)
			continue
		}
		if memReq, err = ci.client.GetCurrentRequests(promql.MEMORY, cluster); err != nil {
			log.Warnf("Type[%v], Cluster[%v], Action[skip-new-values]. Could not get current requests. Reason: %v.", promql.MEMORY, cluster, err)
			continue
		}
		if memLim, err = ci.client.GetCurrentLimits(promql.MEMORY, cluster); err != nil {
			log.Warnf("Type[%v], Cluster[%v], Action[skip-new-values]. Could not get current limits. Reason: %v.", promql.MEMORY, cluster, err)
			continue
		}
		if memAll, err = ci.client.GetAllocatable(promql.MEMORY, cluster); err != nil {
			log.Warnf("Type[%v], Cluster[%v], Action[skip-new-values]. Could not get allocatable. Reason: %v.", promql.MEMORY, cluster, err)
			continue
		}

		ci.clusters[id].LastUpdate = time.Now()
		ci.clusters[id].SetCpuReq(cpuReq)
		ci.clusters[id].SetCpuLim(cpuLim)
		ci.clusters[id].SetCpuAlloc(cpuAll)
		ci.clusters[id].SetMemReq(memReq)
		ci.clusters[id].SetMemLim(memLim)
		ci.clusters[id].SetMemAlloc(memAll)
	}

	time.Sleep(5 * time.Second)
	ci.updateClustersInfo()
}

func (ci *ClustersInfo) GetClusterReq(resType promql.Resource, clusterProvider, clusterName string) (float64, error) {
	cluster, err := ci.GetCluster(clusterProvider, clusterName)
	if err != nil {
		err = fmt.Errorf("couldn't get cluster %v requests. Reason: %v", resType, err)
		log.Errorf("[KSM] error: %v", err)
		return -1, err
	}
	if resType == promql.CPU {
		return cluster.GetCpuReq(), nil
	} else if resType == promql.MEMORY {
		return cluster.GetMemReq(), nil
	} else {
		return -1, errors.New(fmt.Sprintf("couldn't get cluster requests. Reason: resType[%v] doesn't exist", resType))
	}
}

func (ci *ClustersInfo) GetClusterLim(resType promql.Resource, clusterProvider, clusterName string) (float64, error) {
	cluster, err := ci.GetCluster(clusterProvider, clusterName)
	if err != nil {
		err = fmt.Errorf("couldn't get cluster %v limits. Reason: %v", resType, err)
		log.Errorf("[KSM] error: %v", err)
		return -1, err
	}
	if resType == promql.CPU {
		return cluster.GetCpuLim(), nil
	} else if resType == promql.MEMORY {
		return cluster.GetMemLim(), nil
	} else {
		return -1, errors.New(fmt.Sprintf("couldn't get cluster limits. Reason: resType[%v] doesn't exist", resType))
	}
}

func (ci *ClustersInfo) GetClusterAlloc(resType promql.Resource, clusterProvider, clusterName string) (float64, error) {
	cluster, err := ci.GetCluster(clusterProvider, clusterName)
	if err != nil {
		err = fmt.Errorf("couldn't get cluster %v allocatable. Reason: %v", resType, err)
		log.Errorf("[KSM] error: %v", err)
		return -1, err
	}
	if resType == promql.CPU {
		return cluster.GetCpuAlloc(), nil
	} else if resType == promql.MEMORY {
		return cluster.GetMemAlloc(), nil
	} else {
		return -1, errors.New(fmt.Sprintf("couldn't get cluster allocatable. Reason: resType[%v] doesn't exist", resType))
	}
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
