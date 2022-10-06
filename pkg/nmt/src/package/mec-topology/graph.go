package mec_topology

import (
	"fmt"
)

type Graph struct {
	MecHosts     []*MecHost
	Edges        []*Edge
	NetworkCells []*Cell
}

func (g *Graph) GetMecHost(clusterName, clusterProvider string) *MecHost {
	//getVertexHandler return a pointer to the mecHost with a specific name and provider

	for i, v := range g.MecHosts {
		if v.Identity.ClusterName == clusterName && v.Identity.Provider == clusterProvider {
			return g.MecHosts[i]
		}
	}
	return nil
}

func (g *Graph) GetCell(cellId string) *Cell {

	for i, v := range g.NetworkCells {
		if v.Id == cellId {
			return g.NetworkCells[i]
		}

	}
	return nil
}

//methods that adds MEC Host to the list of MEC HOSTS of given g Graph
func (g *Graph) AddMecHost(mecHost MecHost) {
	if g.CheckGraphContainsVertex(mecHost) {
		err := fmt.Errorf("Vertex %v not added beacuse already exist vertex with the same name and provider id \n", mecHost.Identity.ClusterName, mecHost.Identity.Provider)
		fmt.Println(err.Error())
	} else {
		g.MecHosts = append(g.MecHosts, &mecHost)
		fmt.Printf("Added new mec host:  %v\n", mecHost)
	}
}

func (g *Graph) AddLink(edge Edge) {

	//get vertex
	fromMECHost := g.GetMecHost(edge.SourceVertexName, edge.SourceVertexProviderName)
	toMECHost := g.GetMecHost(edge.TargetVertexName, edge.TargetVertexProviderName)

	//check error
	if fromMECHost == nil || toMECHost == nil {
		err := fmt.Errorf("Invalid edge- at least one of MEC Host doesn't exist (%v, %v <--> %v, %v)\n", edge.SourceVertexName, edge.SourceVertexProviderName, edge.TargetVertexName, edge.TargetVertexProviderName)
		fmt.Println(err.Error())
	} else if g.CheckAlreadExistLink(edge) {
		err := fmt.Errorf("Edge between (%v--%v) already exist\n", edge.SourceVertexName, edge.TargetVertexName)
		fmt.Println(err.Error())
	} else {
		//add neighbours at vertexes instances
		fromMECHost.Neighbours = append(fromMECHost.Neighbours, toMECHost.Identity)
		toMECHost.Neighbours = append(toMECHost.Neighbours, fromMECHost.Identity)

		//add edge at  Edges list
		g.Edges = append(g.Edges, &edge)
		fmt.Printf("New Edge added : %v %v --- %v %v \n", edge.SourceVertexName, edge.SourceVertexProviderName, edge.TargetVertexName, edge.TargetVertexProviderName)
	}
}

//that function prints graph: Nodes and Links
func (g *Graph) PrintGraph() {

	fmt.Println("Graph: ")
	//print vertexes
	for _, v := range g.MecHosts {
		fmt.Printf("\nVertex: %v %v", v.Identity.ClusterName, v.Identity.Provider)
		fmt.Print(*v)
	}
	fmt.Println()

	//print edges
	for _, v := range g.Edges {
		fmt.Printf("Edge between: %v and %v\n", v.SourceVertexName, v.TargetVertexName)
	}

}
