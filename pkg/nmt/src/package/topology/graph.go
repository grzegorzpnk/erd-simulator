package topology

import (
	"fmt"
)

type Graph struct {
	MecHosts []*MecHost
	Edges    []*Edge
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

func (g *Graph) AddMecHost(mecHost MecHost) {
	if CheckGraphContainsVertex(g.MecHosts, mecHost) {
		err := fmt.Errorf("Vertex %v not added beacuse already exist vertex with the same name and provider id \n", mecHost.Identity.ClusterName, mecHost.Identity.Provider)
		fmt.Println(err.Error())
	} else {
		g.MecHosts = append(g.MecHosts, &mecHost)
		fmt.Printf("Added new mec host:  %v\n", mecHost)
	}
}

func (g *Graph) AddEdge(edge Edge) {

	//get vertex
	fromVertex := g.GetMecHost(edge.SourceVertexName, edge.SourceVertexProviderName)
	toVertex := g.GetMecHost(edge.TargetVertexName, edge.TargetVertexProviderName)

	//check error
	if fromVertex == nil || toVertex == nil {
		err := fmt.Errorf("Invalid edge- at least one of Vertex not exists (%v<-->%v)\n", edge.SourceVertexName, edge.TargetVertexName)
		fmt.Println(err.Error())
	} else if CheckGraphContainsNeighbour(fromVertex.Neighbours, edge.TargetVertexName) || CheckGraphContainsNeighbour(toVertex.Neighbours, edge.SourceVertexName) {
		err := fmt.Errorf("Edge between (%v--%v) already exist\n", edge.SourceVertexName, edge.TargetVertexName)
		fmt.Println(err.Error())
	} else {
		//add edge at vertexes instances
		fromVertex.Neighbours = append(fromVertex.Neighbours, edge.TargetVertexName)
		toVertex.Neighbours = append(toVertex.Neighbours, edge.SourceVertexName)

		//add edge at  Edges list
		g.Edges = append(g.Edges, &edge)
		fmt.Printf("New Edge added : %v %v --- %v %v \n", edge.SourceVertexName, edge.TargetVertexName)
	}
}

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
