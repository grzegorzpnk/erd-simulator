package results

type AppType string

const (
	CG  AppType = "10"
	V2X AppType = "15"
	UAV AppType = "30"
)

type ResultCounter map[AppType]int

type ExperimentStrategy string

const (
	StrOptimal   ExperimentStrategy = "optimal"
	StrHeuristic ExperimentStrategy = "heuristic"
)

type ExperimentObjective string

const (
	ObjLb       ExperimentObjective = "lb"
	ObjLatency  ExperimentObjective = "latency"
	ObjHybrid   ExperimentObjective = "hybrid"
	ObjHybridIf ExperimentObjective = "hybrid-if"
)

type Client struct {
	expResults []ExpResult
}

func (c *Client) GetResults() []ExpResult {
	return c.expResults
}

func (c *Client) AppendResult(result ExpResult) {
	c.expResults = append(c.expResults, result)
}

func NewClient() *Client {
	return &Client{
		expResults: []ExpResult{},
	}
}

type ExpResult struct {
	Metadata ExpResultsMeta `json:"metadata"`
	Data     ExpResultsData `json:"data"`
}

type ExpResultsMeta struct {
	Strategy  ExperimentStrategy
	Objective ExperimentObjective
}

type ExpResultsData struct {
	Erd  ErdResults  `json:"erdResults,omitempty"`
	Topo TopoResults `json:"topoResults,omitempty"`
}

type ErdResults struct {
	Failed     ResultCounter   `json:"relocation-failed"`
	Successful ResultCounter   `json:"relocation-successful"`
	Redundant  ResultCounter   `json:"relocation-redundant"`
	Skipped    ResultCounter   `json:"relocation-skipped"`
	EvalTimes  EvaluationTimes `json:"evaluation-times"`
}

type TopoResults struct {
	MecHostsResults []MecHostResults
}

// EvaluationTimes times should be in milliseconds
type EvaluationTimes struct {
	Failed     []int `json:"failed"`
	Successful []int `json:"successful"`
	Redundant  []int `json:"redundant"`
	Skipped    []int `json:"skipped"`
}

type MecHostResults struct {
	Identity        MecIdentity      `json:"identity"`
	CpuResources    ClusterResources `json:"cpu_resources,omitempty"`
	MemoryResources ClusterResources `json:"memory_resources,omitempty"`
}

type MecIdentity struct {
	Provider string      `json:"provider"`
	Cluster  string      `json:"cluster"`
	Location MecLocation `json:"location"`
}

type MecType int

const (
	MecLocal MecType = iota
	MecRegional
	MecCentral
)

type MecLocation struct {
	Level  MecType `json:"type"`
	Region string  `json:"region"`
	Zone   string  `json:"zone,omitempty"`
	//city-level eg. west, "" if type different from MecLocal and MecRegional
	LocalZone string `json:"local-zone,omitempty"` // eg. wroclaw, "" if type different from MecLocal
}

type ClusterResources struct {
	Used        float64 `json:"used"`        // How many cpu/memory used (value)
	Capacity    float64 `json:"capacity"`    // How many cpu/memory available (value)
	Utilization float64 `json:"utilization"` // How much is cpu/memory utilized (percentage)
}
