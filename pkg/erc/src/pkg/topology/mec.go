package topology

type MecInfo string
type MecType int

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

type MecResources struct {
	Latency float64    `json:"latency"`
	Cpu     MecResInfo `json:"cpu"`
	Memory  MecResInfo `json:"memory"`
}

type MecLocation struct {
	Type      MecType `json:"type"`
	Region    string  `json:"region"`               // eg. poland
	Zone      string  `json:"zone,omitempty"`       // eg. west, "" if type different from MecLocal and MecRegional
	LocalZone string  `json:"local-zone,omitempty"` // eg. wroclaw, "" if type different from MecLocal
}

type MecHost struct {
	Identity   MecIdentity  `json:"identity"`
	Resources  MecResources `json:"resources,omitempty"`
	Neighbours []*MecHost
}

func (mh *MecHost) BuildClusterEmcoFQDN() string {
	return mh.Identity.Provider + "+" + mh.Identity.Cluster
}

func (mh *MecHost) GetLatency() float64 {
	return mh.Resources.Latency
}

func (mh *MecHost) SetLatency(latency float64) {
	mh.Resources.Latency = latency
}

func (mh *MecHost) GetCpuUsed() float64 {
	return mh.Resources.Cpu.Used
}

func (mh *MecHost) GetCpuAllocatable() float64 {
	return mh.Resources.Cpu.Allocatable
}

func (mh *MecHost) GetCpuUtilization() float64 {
	return mh.Resources.Cpu.Utilization
}

func (mh *MecHost) SetCpuUtilization(cpuUtilization float64) {
	mh.Resources.Cpu.Utilization = cpuUtilization
}

func (mh *MecHost) GetMemUsed() float64 {
	return mh.Resources.Memory.Used
}

func (mh *MecHost) GetMemAllocatable() float64 {
	return mh.Resources.Memory.Allocatable
}

func (mh *MecHost) GetMemUtilization() float64 {
	return mh.Resources.Memory.Utilization
}

func (mh *MecHost) SetMemUtilization(memUtilization float64) {
	mh.Resources.Memory.Utilization = memUtilization
}

func (mh *MecHost) GetNeighbours() []*MecHost {
	return mh.Neighbours
}

func (mh *MecHost) AddNeighbour(mec MecHost) {
	mh.Neighbours = append(mh.Neighbours, &mec)
}
