package model

type MecInfo string
type MecType int

type MecHost struct {
	Identity        MecIdentity `json:"identity"`
	CpuResources    MecResInfo  `json:"cpu_resources,omitempty"`
	MemoryResources MecResInfo  `json:"memory_resources,omitempty"`
	Neighbours      []MecIdentity
	SupportingCells []Cell `json:"supporting_cells"`
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
	Provider string      `json:"provider"`
	Cluster  string      `json:"cluster"`
	Location MecLocation `json:"location"`
}

type MecResInfo struct {
	Used        float64 `json:"used"`        // How many cpu/memory used (value)
	Allocatable float64 `json:"allocatable"` // How many cpu/memory available (value)
	Utilization float64 `json:"utilization"` // How much is cpu/memory utilized (percentage)
}

type MecLocation struct {
	Level  MecType `json:"type"`
	Region string  `json:"region"`
	Zone   string  `json:"zone,omitempty"`
	//city-level eg. west, "" if type different from MecLocal and MecRegional
	LocalZone string `json:"local-zone,omitempty"` // eg. wroclaw, "" if type different from MecLocal
}

type Cell struct {
	Id        string  `json:"id"`
	Latency   float32 `json:"latency"`
	LocalZone string  `json:"local-zone,omitempty"`
}

func (m *MecHost) GetCell(cellId string) *Cell {

	for i, v := range m.SupportingCells {
		if v.Id == cellId {
			return &m.SupportingCells[i]
		}
	}
	return nil
}

//check if cluster supports given Cell ID
func (mec *MecHost) CheckMECsupportsCell(cellId string) bool {

	for _, v := range mec.SupportingCells {
		if v.Id == cellId {
			return true
			break
		}
	}
	return false
}

func (mh *MecHost) BuildClusterEmcoFQDN() string {
	if mh.Identity.Provider != "" && mh.Identity.Cluster != "" {
		return mh.Identity.Provider + "+" + mh.Identity.Cluster
	}
	return ""
}
