package djikstra

import (
	"fmt"
	"math"
	"sync"
)

func GetShortestPath(startNode *Node, endNode *Node, g *ItemGraph) ([]string, int) {
	visited := make(map[string]bool)
	dist := make(map[string]int)
	prev := make(map[string]string)
	//pq := make(PriorityQueue, 1)
	//heap.Init(&pq)
	q := NodeQueue{}
	pq := q.NewQ()
	start := Vertex{
		Node:     startNode,
		Distance: 0,
	}
	for _, nval := range g.Nodes {
		dist[nval.Value] = math.MaxInt64
	}
	dist[startNode.Value] = start.Distance
	pq.Enqueue(start)
	//im := 0
	for !pq.IsEmpty() {
		v := pq.Dequeue()
		if visited[v.Node.Value] {
			continue
		}
		visited[v.Node.Value] = true
		near := g.Edges[*v.Node]

		for _, val := range near {
			if !visited[val.Node.Value] {
				if dist[v.Node.Value]+val.Weight < dist[val.Node.Value] {
					store := Vertex{
						Node:     val.Node,
						Distance: dist[v.Node.Value] + val.Weight,
					}
					dist[val.Node.Value] = dist[v.Node.Value] + val.Weight
					//prev[val.Node.Value] = fmt.Sprintf("->%s", v.Node.Value)
					prev[val.Node.Value] = v.Node.Value
					pq.Enqueue(store)
				}
				//visited[val.Node.value] = true
			}
		}
	}
	fmt.Println(dist)
	fmt.Println(prev)
	pathval := prev[endNode.Value]
	var finalArr []string
	finalArr = append(finalArr, endNode.Value)
	for pathval != startNode.Value {
		finalArr = append(finalArr, pathval)
		pathval = prev[pathval]
	}
	finalArr = append(finalArr, pathval)
	fmt.Println(finalArr)
	for i, j := 0, len(finalArr)-1; i < j; i, j = i+1, j-1 {
		finalArr[i], finalArr[j] = finalArr[j], finalArr[i]
	}
	return finalArr, dist[endNode.Value]

}

type Node struct {
	Value string
}

type Edge struct {
	Node   *Node
	Weight int
}

type Vertex struct {
	Node     *Node
	Distance int
}

/*type ItemGraph struct {
	Nodes []*Node
	Edges map[Node][]*Edge
	Lock  sync.RWMutex
}*/

type PriorityQueue []*Vertex

type NodeQueue struct {
	Items []Vertex
	Lock  sync.RWMutex
}

type InputGraph struct {
	Graph []InputData `json:"graph"`
	From  string      `json:"from"`
	To    string      `json:"to"`
}

type InputData struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Weight      int    `json:"weight"`
}
