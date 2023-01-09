package model

// MecInfo is used to define MEC Host resources types (cpu, memory)
type MecInfo string

const (
	MecCpu MecInfo = "cpu"
	MecMem MecInfo = "memory"
)

// MecType is used to define MEC Host type (3 types defined)
type MecType int

const (
	MecCityLevel MecType = iota
	MecRegionalLevel
	MecInternationalLevel
)

// MecIdentity is an unambiguous MEC Host representation
type MecIdentity struct {
	Provider string      `json:"provider"`
	Cluster  string      `json:"cluster"`
	Location MecLocation `json:"location"`
}

// MecResInfo contains information about MEC Host resource defined using MecInfo
type MecResInfo struct {
	Used        float64 `json:"used"`        // How many cpu/memory used (value)
	Capacity    float64 `json:"capacity"`    // How many cpu/memory available (value)
	Utilization float64 `json:"utilization"` // How much is cpu/memory utilized (percentage)
}

// MecResources is a struct which defined MEC Hosts resources, including CPU, Memory but also E2E Latency
type MecResources struct {
	Latency float64    `json:"latency"`
	Cpu     MecResInfo `json:"cpu"`
	Memory  MecResInfo `json:"memory"`
}

// MecLocation struct is an unambiguous MEC Host location
type MecLocation struct {
	Type      MecType `json:"type"`
	Region    string  `json:"region"`               // eg. poland
	Zone      string  `json:"zone,omitempty"`       // eg. west, "" if type different from MecLocal and MecRegional
	LocalZone string  `json:"local-zone,omitempty"` // eg. wroclaw, "" if type different from MecLocal
}

// MecHost struct is a representation of a MEC Host (Edge Cluster) which contains all the necessary information
type MecHost struct {
	Identity   MecIdentity  `json:"identity"`
	Resources  MecResources `json:"resources,omitempty"`
	Neighbours []*MecHost
}

// BuildClusterEmcoFQDN returns cluster identifier, which is compliant with EMCO representation - provider+cluster
func (mh *MecHost) BuildClusterEmcoFQDN() string {
	if mh.Identity.Provider != "" && mh.Identity.Cluster != "" {
		return mh.Identity.Provider + "+" + mh.Identity.Cluster
	}
	return ""
}

// Getters & Setters

func (mh *MecHost) GetLatency() float64 {
	return mh.Resources.Latency
}

func (mh *MecHost) SetLatency(latency float64) {
	mh.Resources.Latency = latency
}

func (mh *MecHost) GetCpuUsed() float64 {
	return mh.Resources.Cpu.Used
}

func (mh *MecHost) GetCpuCapacity() float64 {
	return mh.Resources.Cpu.Capacity
}

func (mh *MecHost) GetCpuUtilization() float64 {
	return mh.Resources.Cpu.Utilization
}

func (mh *MecHost) SetCpuInfo(cpuInfo MecResInfo) {
	mh.Resources.Cpu = cpuInfo
}

func (mh *MecHost) GetMemUsed() float64 {
	return mh.Resources.Memory.Used
}

func (mh *MecHost) GetMemCapacity() float64 {
	return mh.Resources.Memory.Capacity
}

func (mh *MecHost) GetMemUtilization() float64 {
	return mh.Resources.Memory.Utilization
}

func (mh *MecHost) SetMemInfo(memInfo MecResInfo) {
	mh.Resources.Memory = memInfo
}

func (mh *MecHost) GetNeighbours() []*MecHost {
	return mh.Neighbours
}

func (mh *MecHost) AddNeighbour(mec MecHost) {
	mh.Neighbours = append(mh.Neighbours, &mec)
}
