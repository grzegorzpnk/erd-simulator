package api

import (
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/djikstra"
	mec_topology "10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/mec-topology"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/model"
	"log"

	"fmt"
)

func containsAnyEdge(vertex model.MecHost) bool {

	if vertex.Neighbours != nil {
		return true
	} else {
		return false
	}

}

type ShortestPathResult struct {
	latencyResults float64
	path           []string
}

func ShortestPath(startCell *model.Cell, destCluster *model.MecHost, graph *mec_topology.Graph) float64 {

	var min float64

	if destCluster == nil {
		log.Fatalln("destination MEC host not recognized!")

	}
	//check if they are direct neighbours, if so the latency is just between start and stop node
	if destCluster.CheckMECsupportsCell(startCell.Id) {
		min = destCluster.GetCell(startCell.Id).Latency
		fmt.Printf("direct nodes, latency between cell: %v and mec: [%v+%v], is: %v", startCell.Id, destCluster.Identity.Provider, destCluster.Identity.Cluster, destCluster.GetCell(startCell.Id).Latency)

	} else {
		// if not, we have to calculate path between all MEC clusters that are in the same local zone as cell, to the target cluster, the final latency is a sum of the calculated one + between started mec and cell
		var startClusters []model.MecHost

		for _, v := range graph.MecHosts {
			if v.Identity.Location.LocalZone == startCell.LocalZone {
				startClusters = append(startClusters, *v)
			}
		}

		var inputGraph djikstra.InputGraph
		inputGraph.Graph = make([]djikstra.InputData, 200)

		//add all mec hosts to temp graph todo: should be only subset of graph nodes
		for i, v := range graph.Edges {
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
			resultTmp.latencyResults += graph.GetMecHost(v.Identity.Cluster, v.Identity.Provider).GetCell(startCell.Id).Latency

			results = append(results, resultTmp)

		}

		//find minimal value
		min := results[0].latencyResults
		for _, v := range results {
			if v.latencyResults < min {
				min = v.latencyResults
			}
		}

		for i, v := range results {

			if v.latencyResults == min {
				fmt.Printf("final path is: %v\n", results[i].path)
			}
		}

	}
	return min
}
