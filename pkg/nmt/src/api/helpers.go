package api

import mec_topology "nmt/src/package/mec-topology"

func containsAnyEdge(vertex mec_topology.MecHost) bool {

	if vertex.Neighbours != nil {
		return true
	} else {
		return false
	}

}
