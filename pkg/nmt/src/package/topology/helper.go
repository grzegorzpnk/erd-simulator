package topology

func CheckGraphContainsVertex(s []*MecHost, mecHost MecHost) bool {

	for _, v := range s {
		if mecHost.Identity.ClusterName == v.Identity.ClusterName &&
			mecHost.Identity.Provider == v.Identity.ClusterName {
			return true
		}
	}
	return false
}

func CheckGraphContainsNeighbour(s []string, k string) bool {

	for _, v := range s {
		if k == v {
			return true
		}
	}
	return false

}
