package topology

import (
	"fmt"
)

type Graph struct {
	Vertices []*Vertex
	Edges    []*Edge
}

type Vertex struct {
	Id            string         `json:"id"`
	Type          string         `json:"type"` //MEC or CELL
	Neighbours    []string       `json:"neighbours"`
	VertexMetrics ClusterMetrics `json:"vertexMetrics"`
	EdgeProvider  string         `json:"edgeProvider"`
}

type Edge struct {
	Source      string         `json:"source"`
	Target      string         `json:"target"`
	EdgeMetrics NetworkMetrics `json:"edgeMetrics"`
}

func (g *Graph) GetVertex(id string) *Vertex {
	//getVertexHandler return a pointer to the Vertex with a key int

	for i, v := range g.Vertices {
		if v.Id == id {
			return g.Vertices[i]
		}
	}
	return nil
}

func (g *Graph) AddVertex(vertex Vertex) {
	if ContainsVertex(g.Vertices, vertex) {
		err := fmt.Errorf("Vertex %v not added beacuse already exist vertex with the same id and the same type\n", vertex.Id, vertex.Type)
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
	} else if containsNeighbour(fromVertex.Neighbours, edge.Target) || containsNeighbour(toVertex.Neighbours, edge.Source) {
		err := fmt.Errorf("Edge between (%v--%v) already exist\n", edge.Source, edge.Target)
		fmt.Println(err.Error())
	} else {
		//add edge at vertexes instances
		fromVertex.Neighbours = append(fromVertex.Neighbours, edge.Target)
		toVertex.Neighbours = append(toVertex.Neighbours, edge.Source)

		//add edge at  Edges list
		g.Edges = append(g.Edges, &edge)
		fmt.Printf("New Edge added : %v %v --- %v %v \n", edge.Source, fromVertex.Type, edge.Target, toVertex.Type)
	}
}

func (g *Graph) PrintGraph() {

	fmt.Println("Graph: ")
	//print vertexes
	for _, v := range g.Vertices {
		fmt.Printf("\nVertex: %v : ", v.Id)
		fmt.Print(*v)
	}
	fmt.Println()

	//print edges
	for _, v := range g.Edges {
		fmt.Printf("Edge between: %v and %v\n", v.Source, v.Target)
	}

}

func ContainsVertex(s []*Vertex, vertex Vertex) bool {

	for _, v := range s {
		if vertex.Id == v.Id {
			return true
		}
	}
	return false
}

func containsNeighbour(s []string, k string) bool {

	for _, v := range s {
		if k == v {
			return true
		}
	}
	return false

}
