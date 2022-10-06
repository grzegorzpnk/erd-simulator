package mec_topology

import (
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/config"
	log "10.254.188.33/matyspi5/erd/pkg/nmt/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/metrics"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/model"
	"errors"

	"fmt"
	"time"
)

type Graph struct {
	MecHosts []*model.MecHost
	Edges    []*model.Edge
}

func (g *Graph) GetMecHost(clusterName, clusterProvider string) *model.MecHost {
	//getVertexHandler return a pointer to the mecHost with a specific name and provider

	for i, v := range g.MecHosts {
		if v.Identity.Cluster == clusterName && v.Identity.Provider == clusterProvider {
			return g.MecHosts[i]
		}
	}
	return nil
}

//AddMecHost method that adds MEC Host to the list of MEC HOSTS of given g Graph
func (g *Graph) AddMecHost(mecHost model.MecHost) {
	if g.CheckGraphContainsVertex(mecHost) {
		err := errors.New(fmt.Sprintf("Vertex %v+%v not added beacuse already exist vertex with the same name and provider id", mecHost.Identity.Provider, mecHost.Identity.Cluster))
		log.Errorf(err.Error())
	} else {
		g.MecHosts = append(g.MecHosts, &mecHost)
		log.Infof("Added new mec host:  %v\n", mecHost)
	}
}

func (g *Graph) AddEdge(edge model.Edge) {

	//get vertex
	fromMECHost := g.GetMecHost(edge.SourceVertexName, edge.SourceVertexProviderName)
	toMECHost := g.GetMecHost(edge.TargetVertexName, edge.TargetVertexProviderName)

	//check error
	if fromMECHost == nil || toMECHost == nil {
		err := fmt.Errorf("Invalid edge- at least one of MEC Host doesn't exist (%v, %v <--> %v, %v)\n", edge.SourceVertexName, edge.SourceVertexProviderName, edge.TargetVertexName, edge.TargetVertexProviderName)
		log.Errorf(err.Error())
	} else if g.CheckAlreadExistLink(edge) {
		err := errors.New(fmt.Sprintf("Edge between (%v--%v) already exist", edge.SourceVertexName, edge.TargetVertexName))
		log.Warnf(err.Error())
	} else {
		//add neighbours at vertexes instances
		fromMECHost.Neighbours = append(fromMECHost.Neighbours, toMECHost.Identity)
		toMECHost.Neighbours = append(toMECHost.Neighbours, fromMECHost.Identity)

		//add edge at  Edges list
		g.Edges = append(g.Edges, &edge)
		log.Infof("New Edge added : %v %v --- %v %v \n", edge.SourceVertexName, edge.SourceVertexProviderName, edge.TargetVertexName, edge.TargetVertexProviderName)
	}
}

//PrintGraph logs graph: Nodes and Links
func (g *Graph) PrintGraph() {

	log.Infof("Graph: ")
	//print vertexes
	for _, v := range g.MecHosts {
		log.Infof("Vertex: %v %v", v.Identity.Cluster, v.Identity.Provider)
		log.Info(*v)
	}

	//print edges
	for _, v := range g.Edges {
		log.Infof("Edge between: %v and %v\n", v.SourceVertexName, v.TargetVertexName)
	}

}

//ClustersResourcesUpdate is a gorutine function
func (g *Graph) ClustersResourcesUpdate() {

	endpoint := config.GetConfiguration().ClusterControllerEndpoint

	for {
		// update metrics for MEC Clusters
		for _, v := range g.MecHosts {

			clusterCPUURL := metrics.BuildCpuUrl(v.Identity.Cluster, v.Identity.Provider, endpoint)
			clusterMemoryURL := metrics.BuildMemoryUrl(v.Identity.Cluster, v.Identity.Provider, endpoint)
			log.Infof("update for cluster %v\n", v.Identity.Cluster)
			log.Infof("cpu latest update:")
			cpuCr, err := metrics.GetClusterMetrics(v.Identity.Cluster, v.Identity.Provider, clusterCPUURL)
			if err != nil {
				log.Errorf(err.Error())
			}
			log.Infof("memory latest update:")
			memoryCr, err := metrics.GetClusterMetrics(v.Identity.Cluster, v.Identity.Provider, clusterMemoryURL)
			if err != nil {
				log.Errorf(err.Error())
			}

			g.GetMecHost(v.Identity.Cluster, v.Identity.Provider).CpuResources.UpdateClusterMetrics(cpuCr)
			g.GetMecHost(v.Identity.Cluster, v.Identity.Provider).MemoryResources.UpdateClusterMetrics(memoryCr)
			log.Infof("Controller updates cluster metrics for Mec Host: %v\n", v.Identity.Cluster)

		}
		time.Sleep(1 * time.Second)
	}
}
