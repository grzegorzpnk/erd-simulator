package mec_topology

import (
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/config"
	log "10.254.188.33/matyspi5/erd/pkg/nmt/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/metrics"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/model"
	"errors"
	"math/rand"

	"fmt"
	"time"
)

type Graph struct {
	MecHosts     []*model.MecHost
	Edges        []*model.Edge
	NetworkCells []*model.Cell
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

func (g *Graph) GetCell(cellId string) *model.Cell {

	for i, v := range g.NetworkCells {
		if v.Id == cellId {
			return g.NetworkCells[i]
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

func (g *Graph) AddLink(edge model.Edge) {

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
			//	log.Infof("update for cluster %v\n", v.Identity.Cluster)
			//	log.Infof("cpu latest update:")
			cpuCr, err := metrics.GetClusterMetrics(clusterCPUURL)
			if err != nil {
				log.Errorf(err.Error())
			}
			//	log.Infof("memory latest update:")
			memoryCr, err := metrics.GetClusterMetrics(clusterMemoryURL)
			if err != nil {
				log.Errorf(err.Error())
			}

			g.GetMecHost(v.Identity.Cluster, v.Identity.Provider).CpuResources.UpdateClusterMetrics(cpuCr)
			g.GetMecHost(v.Identity.Cluster, v.Identity.Provider).MemoryResources.UpdateClusterMetrics(memoryCr)
			//	log.Infof("Controller updates cluster metrics for Mec Host: %v\n", v.Identity.Cluster)

		}
		time.Sleep(1 * time.Second)
	}
}

func (g *Graph) NetworkMetricsUpdate(fetchLatencyFromObserver bool) {

	endpoint := config.GetConfiguration().ClusterControllerEndpoint

	if !fetchLatencyFromObserver {
		rand.Seed(4)
		//update cell-mecs latencies - this is kept by MEC hosts
		for i, v := range g.MecHosts {
			for k, _ := range v.SupportingCells {
				latency := 8 + 2*rand.Float64()
				g.MecHosts[i].SupportingCells[k].Latency = latency
			}
		}

		// update metrics for MEC Clusters
		for d, b := range g.Edges {
			source := g.GetMecHost(b.SourceVertexName, b.SourceVertexProviderName)
			target := g.GetMecHost(b.TargetVertexName, b.TargetVertexProviderName)
			latency, err := generateLatency(*source, *target)
			if err != nil {
				log.Errorf(err.Error())
			} else {
				g.Edges[d].EdgeMetrics.Latency = latency
			}
		}
	} else {
		for {

			//update cell-mecs latencies - this is kept by MEC hosts
			for i, v := range g.MecHosts {
				for k, c := range v.SupportingCells {
					cellLatenyUrl := metrics.BuildCellLatencyURL(endpoint, c.Id, v.Identity.Cluster, v.Identity.Provider)
					latency, err := metrics.GetLatency(cellLatenyUrl)
					if err != nil {
						log.Errorf(err.Error())
					} else {
						g.MecHosts[i].SupportingCells[k].Latency = latency
					}
				}
			}

			//update mecs-mecs latencies - this is kept by Edges

			// update metrics for MEC Clusters
			for d, b := range g.Edges {

				mecLatenyUrl := metrics.BuildMECLatencyURL(endpoint, b.TargetVertexName, b.TargetVertexProviderName, b.SourceVertexName, b.SourceVertexProviderName)
				latency, err := metrics.GetLatency(mecLatenyUrl)
				if err != nil {
					log.Errorf(err.Error())
				} else {
					g.Edges[d].EdgeMetrics.Latency = latency
				}

			}

			time.Sleep(10 * time.Second)
		}
	}
}

func generateLatency(sNode, tNode interface{}) (float64, error) {
	rand.Seed(4)
	var latency float64

	// copied from observer. Ofc inhere it's impossible to provider cell-id, but
	// I'm going to leave it for now in such form
	switch source := sNode.(type) {
	//case model.Cell:
	//	switch target := tNode.(type) {
	//	case model.Cell:
	//		// source -> cell; target -> cell
	//		// not allowed
	//	case model.MecHost:
	//		// source -> cell; target -> mecHost
	//		if target.Identity.Location.Level == 0 {
	//			latency = 0.5
	//		} else {
	//			latency = -404
	//		}
	//	}
	case model.MecHost:
		switch target := tNode.(type) {
		//case model.Cell:
		//	// source -> mecHost; target -> cell
		//	if source.Identity.Location.Level == 0 {
		//		latency = 0.5
		//	} else {
		//		latency = -404
		//	}
		case model.MecHost:
			if source.Identity == target.Identity {
				return 0, errors.New("target == source: not supported")
			}
			// source -> mecHost; target -> mecHost
			levelDiff := float64(source.Identity.Location.Level - target.Identity.Location.Level)

			switch levelDiff {
			case 0:
				if source.Identity.Location.Level == 0 {
					if source.Identity.Location.LocalZone == target.Identity.Location.LocalZone {
						latency = 1.0 // Latency between clusters at N-level in the same LocalZone
					} else {
						if (source.Identity.Cluster == "mec17" || source.Identity.Cluster == "mec18") && (target.Identity.Cluster == "mec19" || target.Identity.Cluster == "mec20") ||
							(target.Identity.Cluster == "mec17" || target.Identity.Cluster == "mec18") && (source.Identity.Cluster == "mec19" || source.Identity.Cluster == "mec20") {
							latency = 2.8
						} else {
							latency = 4.0
						}
					}
				} else if source.Identity.Location.Level == 1 {
					if source.Identity.Location.Zone == target.Identity.Location.Zone {
						latency = 1.0 // Latency between clusters at N+1-level in the same Zone
					} else {
						latency = 2.8 // TODO: is it ok? Not included on the model
					}
				} else if source.Identity.Location.Level == 2 {
					if source.Identity.Location.Region == target.Identity.Location.Region {
						latency = 1 // Latency between clusters at N+2-level in the same Region
					} else {
						latency = 10 // TODO: Not included on the model, but not relevant
					}
				}
			case 1:
				if source.Identity.Location.Level == 2 {
					if target.Identity.Cluster == "mec6" || target.Identity.Cluster == "mec7" {
						latency = 0.2 // TODO: for now it's 0.2 but consider random value <0, 0.2>
					} else {
						latency = 2.8
					}
					latency += 4 // DPI time for level N+2
				} else if source.Identity.Location.Level == 1 {
					if ((source.Identity.Cluster == "mec3" || source.Identity.Cluster == "mec4" || source.Identity.Cluster == "mec5") &&
						(target.Identity.Cluster == "mec15" || target.Identity.Cluster == "mec16" || target.Identity.Cluster == "mec17" || target.Identity.Cluster == "mec18")) ||
						((source.Identity.Cluster == "mec6" || source.Identity.Cluster == "mec7") &&
							(target.Identity.Cluster == "mec19" || target.Identity.Cluster == "mec20" || target.Identity.Cluster == "mec21" || target.Identity.Cluster == "mec22")) {
						latency = 3
					} else {
						latency = 4
					}
					latency += 4 // DPI time for level N+1
				}
			case -1:
				if source.Identity.Location.Level == 1 {
					if source.Identity.Cluster == "mec6" || source.Identity.Cluster == "mec7" {
						latency = 0.2 // TODO: for now it's 0.2 but consider random value <0, 0.2>
					} else {
						latency = 4
					}
					latency += 4 // DPI time for level N+2
				} else if source.Identity.Location.Level == 0 {
					if ((target.Identity.Cluster == "mec3" || target.Identity.Cluster == "mec4" || target.Identity.Cluster == "mec5") &&
						(source.Identity.Cluster == "mec15" || source.Identity.Cluster == "mec16" || source.Identity.Cluster == "mec17" || source.Identity.Cluster == "mec18")) ||
						((target.Identity.Cluster == "mec6" || target.Identity.Cluster == "mec7") &&
							(source.Identity.Cluster == "mec19" || source.Identity.Cluster == "mec20" || source.Identity.Cluster == "mec21" || source.Identity.Cluster == "mec22")) {
						latency = 3
					} else {
						latency = 4
					}
					latency += 4 // DPI time for level N+1
				}
			case 2, -2:
				// unsupported, no direct links
				latency = -404
			}
		}
	}
	log.Infof("Latency %v <--> %v is %v", sNode, tNode, latency)
	if latency <= 0 {
		return latency, errors.New("could not generate latency properly")
	}
	return latency, nil
}
