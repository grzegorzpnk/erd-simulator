package mec_topology

func (g *Graph) CheckGraphContainsVertex(mecHost MecHost) bool {

	for _, v := range g.MecHosts {
		if mecHost.Identity.ClusterName == v.Identity.ClusterName &&
			mecHost.Identity.Provider == v.Identity.ClusterName {
			return true
		}
	}
	return false
}

// this func checks in bidirectional way
func (g *Graph) CheckAlreadExistLink(k Edge) bool {

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
