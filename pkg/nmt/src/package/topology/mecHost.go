package topology

type MecInfo string
type MecType int

type MecHost struct {
	Identity   MecIdentity  `json:"identity"`
	Resources  MecResources `json:"resources,omitempty"`
	Neighbours []string
}

const (
	MecCpuUtil MecInfo = "/cpu"
	MecMemUtil MecInfo = "/memory"
)

const (
	MecLocal MecType = iota
	MecRegional
	MecCentral
)

type MecIdentity struct {
	Provider    string      `json:"provider"`
	ClusterName string      `json:"cluster"`
	Location    MecLocation `json:"location"`
}

type MecResInfo struct {
	Used        float64 `json:"used"`        // How many cpu/memory used (value)
	Allocatable float64 `json:"allocatable"` // How many cpu/memory available (value)
	Utilization float64 `json:"utilization"` // How much is cpu/memory utilized (percentage)
}

type MecLocation struct {
	Type      MecType `json:"type"`
	Region    string  `json:"region"`               // eg. poland
	Zone      string  `json:"zone,omitempty"`       // eg. west, "" if type different from MecLocal and MecRegional
	LocalZone string  `json:"local-zone,omitempty"` // eg. wroclaw, "" if type different from MecLocal
}
