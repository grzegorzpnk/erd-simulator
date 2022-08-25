package topology

import (
	"fmt"
)

type Graph struct {
	Vertices []*Vertex
	Edges    []*Edge
}

type Vertex struct {
	Id            int            `json:"id"`
	Type          string         `json:"type"` //MEC or CELL
	Neighbours    []int          `json:"neighbours"`
	VertexMetrics ClusterMetrics `json:"vertexMetrics"`
}

type Edge struct {
	Source      int            `json:"source"`
	Target      int            `json:"target"`
	EdgeMetrics NetworkMetrics `json:"edgeMetrics"`
}

func (g *Graph) GetVertex(k int) *Vertex {
	//getVertexHandler return a pointer to the Vertex with a key int

	for i, v := range g.Vertices {
		if v.Id == k {
			return g.Vertices[i]
		}
	}
	return nil
}

func (g *Graph) AddVertex(vertex Vertex) {
	if ContainsVertex(g.Vertices, vertex.Id) {
		err := fmt.Errorf("Vertex %v not added beacuse already exist vertex with the same id\n", vertex.Id)
		fmt.Println(err.Error())
	} else {
		g.Vertices = append(g.Vertices, &vertex)
		fmt.Printf("Added new vertex  %v\n", vertex)
	}
}

func (g *Graph) AddEdge(edge Edge) {

	//get vertex
	fromVertex := g.GetVertex(edge.Source)
	toVertex := g.GetVertex(edge.Target)
	//check error
	if fromVertex == nil || toVertex == nil {
		err := fmt.Errorf("Invalid edge- at least one of Vertex not exists (%v<-->%v)\n", edge.Source, edge.Target)
		fmt.Println(err.Error())
	} else if fromVertex.Type == toVertex.Type {
		err := fmt.Errorf("You cannot connect two Vertexes at the same type:  %v !\n", fromVertex.Type)
		fmt.Println(err.Error())
	} else if containsInt(fromVertex.Neighbours, edge.Target) || containsInt(toVertex.Neighbours, edge.Source) {
		err := fmt.Errorf("Edge between (%v--%v) already exist\n", edge.Source, edge.Target)
		fmt.Println(err.Error())
	} else {
		//add edge at vertexes instances
		fromVertex.Neighbours = append(fromVertex.Neighbours, edge.Target)
		toVertex.Neighbours = append(toVertex.Neighbours, edge.Source)

		//add edge at  Edges list
		g.Edges = append(g.Edges, &edge)
		fmt.Printf("New Edge added : %v --- %v \n", edge.Source, edge.Target)
	}
}

func (g *Graph) PrintGraph() {

	//print vertexes
	for _, v := range g.Vertices {
		fmt.Printf("\nVertex: %v : ", v.Id)
		fmt.Println(v)
		for _, v := range v.Neighbours {
			fmt.Printf("%v  ", v)
		}
	}
	fmt.Println()

	//print edges
	for _, v := range g.Edges {
		fmt.Printf("Edge between: %v and %v\n", v.Source, v.Target)
	}

}

func ContainsVertex(s []*Vertex, k int) bool {

	for _, v := range s {
		if k == v.Id {
			return true
		}
	}
	return false

}

func containsInt(s []int, k int) bool {

	for _, v := range s {
		if k == v {
			return true
		}
	}
	return false

}

/*func getEdge (sourceId, TargetId int) *Edge {

/*	for i,v := range graph.Edges{
		if v.Source == sourceId & v.Target == TargetId
			return gra
		}
	}
*/
