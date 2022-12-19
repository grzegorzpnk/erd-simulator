package api

import (
	mec_topology "10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/model"
)

func containsAnyEdge(vertex mec_topology.MecHost) bool {

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

func (g) CheckGraphContainsVertex() bool {

	return false
}
