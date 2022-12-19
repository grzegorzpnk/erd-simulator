package mec_topology

import "10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/model"

/*func (g *Graph) CheckGraphContainsVertex(mecHost model.MecHost) bool {

	for _, v := range g.MecHosts {
		if mecHost.Identity.Cluster == v.Identity.Cluster &&
			mecHost.Identity.Provider == v.Identity.Cluster {
			return true
		}
	}
	return false
}
*/

// this func checks in bidirectional way
func (g *Graph) CheckAlreadExistLink(k model.Edge) bool {

	for _, v := range g.Edges {
		if (k.SourceVertexName == v.SourceVertexName &&
			k.SourceVertexProviderName == v.SourceVertexProviderName &&
			k.TargetVertexName == v.TargetVertexName &&
			k.TargetVertexProviderName == v.TargetVertexProviderName) ||
			(k.SourceVertexName == v.TargetVertexName &&
				k.SourceVertexProviderName == v.TargetVertexProviderName &&
				k.TargetVertexName == v.SourceVertexName &&
				k.TargetVertexProviderName == v.SourceVertexProviderName) {
			return true
		}

	}
	return false

}
