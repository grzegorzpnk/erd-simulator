package mec_topology

import (
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/djikstra"
	log "10.254.188.33/matyspi5/erd/pkg/nmt/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/model"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
)

type Graph struct {
	MecHosts                 []*model.MecHost
	Edges                    []*model.Edge
	NetworkCells             []*model.Cell
	Application              []*model.MECApp
	ImmutableApplicationList []model.MECApp
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

func (g *Graph) NetworkMetricsUpdate() {

	rand.Seed(42)
	//update cell-mecs latencies - this is kept by MEC hosts
	for i, v := range g.MecHosts {
		for k, _ := range v.SupportingCells {
			latency := 4 + 2*rand.Float64()
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
}

func (g *Graph) FindShortestPathsConfigurationForMECs() {

	for _, v := range g.MecHosts {
		for i := 1; i <= len(g.NetworkCells); i++ {
			cell := g.GetCell(strconv.Itoa(i))
			latency, _ := g.ShortestPath(cell, v)
			v.LatencyVector = append(v.LatencyVector, latency)
		}
	}
}

func generateLatency(sNode, tNode interface{}) (float64, error) {
	rand.Seed(42)
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
						if (source.Identity.Cluster == "mec13" || source.Identity.Cluster == "mec14") && (target.Identity.Cluster == "mec15" || target.Identity.Cluster == "mec16") ||
							(target.Identity.Cluster == "mec13" || target.Identity.Cluster == "mec14") && (source.Identity.Cluster == "mec15" || source.Identity.Cluster == "mec16") ||
							(source.Identity.Cluster == "mec21" || source.Identity.Cluster == "mec22") && (target.Identity.Cluster == "mec27" || target.Identity.Cluster == "mec28") ||
							(target.Identity.Cluster == "mec21" || target.Identity.Cluster == "mec22") && (source.Identity.Cluster == "mec27" || source.Identity.Cluster == "mec28") {
							latency = 5.0 // 2.8
						} else {
							latency = 5.0 // 4.0
						}
					}
				} else if source.Identity.Location.Level == 1 {
					if source.Identity.Location.Zone == target.Identity.Location.Zone {
						latency = 1.0 // Latency between clusters at N+1-level in the same Zone
					} else {
						latency = 5.0 // 2.0
					}
				} else if source.Identity.Location.Level == 2 {
					if source.Identity.Location.Region == target.Identity.Location.Region {
						latency = 1 // Latency between clusters at N+2-level in the same Region
					} else {
						latency = 10
					}
				}
			case 1:
				if source.Identity.Location.Level == 2 {
					if target.Identity.Cluster == "mec5" || target.Identity.Cluster == "mec6" {
						latency = 3
					} else {
						latency = 2.8
					}
					latency += 4 // DPI time for level N+2
				} else if source.Identity.Location.Level == 1 {
					if ((source.Identity.Cluster == "mec2" || source.Identity.Cluster == "mec3" || source.Identity.Cluster == "mec4") &&
						(target.Identity.Cluster == "mec11" || target.Identity.Cluster == "mec12" || target.Identity.Cluster == "mec13" || target.Identity.Cluster == "mec14")) ||
						((source.Identity.Cluster == "mec5" || source.Identity.Cluster == "mec6") &&
							(target.Identity.Cluster == "mec15" || target.Identity.Cluster == "mec16" || target.Identity.Cluster == "mec17" || target.Identity.Cluster == "mec18")) {
						latency = 3
					} else {
						latency = 4
					}
					latency += 6 // DPI time for level N+1
				}
			case -1:
				if source.Identity.Location.Level == 1 {
					if source.Identity.Cluster == "mec6" || source.Identity.Cluster == "mec7" {
						latency = 3
					} else {
						latency = 4
					}
					latency += 4 // DPI time for level N+2
				} else if source.Identity.Location.Level == 0 {
					if ((target.Identity.Cluster == "mec2" || target.Identity.Cluster == "mec3" || target.Identity.Cluster == "mec4") &&
						(source.Identity.Cluster == "mec11" || source.Identity.Cluster == "mec12" || source.Identity.Cluster == "mec13" || source.Identity.Cluster == "mec14")) ||
						((target.Identity.Cluster == "mec5" || target.Identity.Cluster == "mec6") &&
							(source.Identity.Cluster == "mec15" || source.Identity.Cluster == "mec16" || source.Identity.Cluster == "mec17" || source.Identity.Cluster == "mec18")) {
						latency = 3
					} else {
						latency = 4
					}
					latency += 6 // DPI time for level N+1
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

func (g *Graph) DeclareApplications(ac AppCounter) {

	//v2x, drones, cg := ac.V2x, ac.Uav, ac.Cg
	cg, v2x, drones := ac.Cg, ac.V2x, ac.Uav

	fmt.Printf("Number of declared apps: %v, where: %v of v2x, %v of drones and %v of video\n", ac.GetTotal(), v2x, drones, cg)
	//todo: values should be parametrized, not hardcoded!
	for i := 0; i < cg; i++ {
		var app model.MECApp
		app.Id = strconv.Itoa(i + 1)
		app.Requirements.RequestedLatency = 10
		app.GeneratreResourceRequirements()
		g.Application = append(g.Application, &app)
	}
	for i := cg; i < (cg + v2x); i++ {
		var app model.MECApp
		app.Id = strconv.Itoa(i + 1)
		app.Requirements.RequestedLatency = 15
		app.GeneratreResourceRequirements()
		g.Application = append(g.Application, &app)
	}
	for i := (cg + v2x); i < (cg + v2x + drones); i++ {
		var app model.MECApp
		app.Id = strconv.Itoa(i + 1)
		app.Requirements.RequestedLatency = 30
		app.GeneratreResourceRequirements()
		g.Application = append(g.Application, &app)
	}

	//fmt.Printf("Apps without clusters:\n")
	//for i := 0; i < len(g.Application); i++ {
	//	g.Application[i].PrintApplication()
	//}

}

//func to instanitate all apps from application's declaration lists
func (g *Graph) InstantiateAllDefinedApps() error {

	for _, v := range g.Application {

		mecHost := g.GetMecHost(v.ClusterId, "orange")
		err := mecHost.InstantiateApp(*v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Graph) UninstallAllApps() {

	//todo: I think this function does not uninstall apps, but jus a copy of apps !
	//uninstall on MEC Hosts
	for _, v := range g.MecHosts {
		if len(v.MECApps) != 0 {
			//log.Infof("Cluster: %v", v.Identity.Cluster)
			for _, x := range v.MECApps {
				v.UninstallApp(x)
			}
		}
	}
	log.Infof("Deleted all apps, Resources at clusters:")
	//for _, v := range g.MecHosts {
	//	log.Infof("Cluster: %v, CPU Util: %v, MEM Util: %v", v.Identity.Cluster, v.CpuResources.Utilization, v.MemoryResources.Utilization)
	//}

}

func (g *Graph) DeleteAllDeclaredApps() {

	g.Application = nil
	g.ImmutableApplicationList = nil
}

func (g *Graph) FindInitialClusters() (bool, []model.MecHost) {

	var mecHostSource []model.MecHost
	var cells = map[int]int{}

	//PREREQUESTIES
	//in order to work on copied list of clusters ( cause we need to update resources in process of initial finidng placement)
	for _, v := range g.MecHosts {
		mecHostSource = append(mecHostSource, *v)
	}
	//another copy for each serach iteration ( if one search will fail, wee need to repeat on still fresh data)
	var mecHostsSourcesTmp = make([]model.MecHost, len(mecHostSource))

	var cnt = 0
	var successullyFound = 0
	search := true
	for search {
		if cnt != 0 {
			fmt.Printf("\nIn previous iteration, found place for %v apps", successullyFound)
		}
		successullyFound = 0
		fmt.Println("[DEBUG] Starting search.")
		cnt++
		if cnt > 100 {
			//if you failed more than 20 time let's break the function and return false
			fmt.Printf("Cannot identify initial clusters after %v trials!\n", cnt)
			return false, nil
		}
		copy(mecHostsSourcesTmp, mecHostSource)
		cells = GenerateRandomCellsForUsers(len(g.Application), *g)

		search = false
		for index, edgeApp := range g.Application {
			startCell := g.GetCell(strconv.Itoa(cells[index+1]))
			cmh, err := FindCanidateMec(*edgeApp, startCell, mecHostsSourcesTmp, g)
			if err != nil {
				search = true
				fmt.Printf("Could not find candidate mec for App[%v] due to: %v\n", edgeApp.Id, err.Error())
				break
			}
			mecHostsSourcesTmp = updateMecResourcesInfo(mecHostsSourcesTmp, cmh, *edgeApp)
			edgeApp.ClusterId = cmh.Identity.Cluster
			edgeApp.UserLocation = startCell.Id
			successullyFound += 1
		}

	}

	fmt.Printf("Found after %v iterations", cnt)
	return true, mecHostsSourcesTmp

}

func (g *Graph) ShortestPath(startCell *model.Cell, destCluster *model.MecHost) (float64, error) {

	var min float64

	if destCluster == nil {
		log.Fatalln("destination MEC host not recognized!")
		err := "cannot find destination mec, shortest path finds failed"
		return 0, errors.New(err)

	}
	//check if they are direct neighbours, if so the latency is just between start and stop node
	if destCluster.CheckMECsupportsCell(startCell.Id) {
		min = destCluster.GetCell(startCell.Id).Latency
		//log.Infof("direct nodes, latency between cell: %v and mec: [%v+%v], is: %v", startCell.Id, destCluster.Identity.Provider, destCluster.Identity.Cluster, destCluster.GetCell(startCell.Id).Latency)

	} else {
		// if not, we have to calculate path between all MEC clusters that are in the same local zone as cell, to the target cluster, the final latency is a sum of the calculated one + between started mec and cell
		var startClusters []model.MecHost

		for _, v := range g.MecHosts {
			if v.Identity.Location.LocalZone == startCell.LocalZone {
				startClusters = append(startClusters, *v)
			}
		}

		var inputGraph djikstra.InputGraph
		//max 1000 MEC HOSTS
		inputGraph.Graph = make([]djikstra.InputData, 1000)

		//add all mec hosts to temp graph todo: should be only subset of graph nodes
		for i, v := range g.Edges {
			inputGraph.Graph[i].Source = v.SourceVertexName
			inputGraph.Graph[i].Destination = v.TargetVertexName
			inputGraph.Graph[i].Weight = v.EdgeMetrics.Latency
		}
		itemGraph := djikstra.CreateGraph(inputGraph)

		//calculate shortest path between all []startClusters and stopNode, where startClusters is a list of cluster directly associated with cell
		results := make([]ShortestPathResult, 0)

		for _, v := range startClusters {

			startNd := djikstra.Node{v.Identity.Cluster}
			stopNd := djikstra.Node{destCluster.Identity.Cluster}

			var resultTmp ShortestPathResult
			resultTmp.path, resultTmp.latencyResults = djikstra.GetShortestPath(&startNd, &stopNd, itemGraph)

			//add latency between cell and start MEC host
			resultTmp.latencyResults += g.GetMecHost(v.Identity.Cluster, v.Identity.Provider).GetCell(startCell.Id).Latency

			results = append(results, resultTmp)

		}

		//find minimal value
		min = results[0].latencyResults
		for _, v := range results {
			if v.latencyResults < min {
				min = v.latencyResults
			}
		}

		//log.Infof("indirect nodes, latency between cell: %v and mec: [%v], is: %v", startCell.Id, destCluster.Identity.Cluster, min)
	}
	return min, nil
}

func (g *Graph) UpdateAppCluster(app model.MECApp, destCluster *model.MecHost) error {

	for i, v := range g.Application {

		if v.Id == app.Id {
			g.Application[i].ClusterId = destCluster.Identity.Cluster
			g.Application[i].UserLocation = app.UserLocation
			return nil
		}
	}

	err := errors.New("Cannot find app on declaration list!")
	return err

}
