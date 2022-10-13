package main

import "fmt"

func main() {

	var inputGraph InputGraph
	inputGraph.Graph = make([]InputData, 20)

	inputGraph.Graph[0].Source = "mec1"
	inputGraph.Graph[0].Destination = "mec2"
	inputGraph.Graph[0].Weight = 2

	inputGraph.Graph[1].Source = "mec1"
	inputGraph.Graph[1].Destination = "mec4"
	inputGraph.Graph[1].Weight = 1

	inputGraph.Graph[2].Source = "mec1"
	inputGraph.Graph[2].Destination = "mec3"
	inputGraph.Graph[2].Weight = 5

	inputGraph.Graph[3].Source = "mec2"
	inputGraph.Graph[3].Destination = "mec4"
	inputGraph.Graph[3].Weight = 2

	inputGraph.Graph[4].Source = "mec4"
	inputGraph.Graph[4].Destination = "mec3"
	inputGraph.Graph[4].Weight = 3

	inputGraph.Graph[5].Source = "mec4"
	inputGraph.Graph[5].Destination = "mec5"
	inputGraph.Graph[5].Weight = 1

	inputGraph.Graph[6].Source = "mec2"
	inputGraph.Graph[6].Destination = "mec3"
	inputGraph.Graph[6].Weight = 3

	inputGraph.Graph[7].Source = "mec5"
	inputGraph.Graph[7].Destination = "mec3"
	inputGraph.Graph[7].Weight = 1

	inputGraph.Graph[8].Source = "mec6"
	inputGraph.Graph[8].Destination = "mec3"
	inputGraph.Graph[8].Weight = 5

	inputGraph.Graph[9].Source = "mec5"
	inputGraph.Graph[9].Destination = "mec6"
	inputGraph.Graph[9].Weight = 2

	inputGraph.Graph[10].Source = "mec1"
	inputGraph.Graph[10].Destination = "mec2"
	inputGraph.Graph[10].Weight = 2

	/*
		inputGraph.From = "mec1"
		inputGraph.To = "mec3"
	*/
	itemGraph := CreateGraph(inputGraph)

	var value int
	var path []string
	path, value = GetShortestPath(&Node{"mec1"}, &Node{"mec3"}, itemGraph)

	fmt.Printf("path value: %v, path: %v", value, path)

}
