package results

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	log "simu/src/logger"
	"strconv"
)

type AppType string

const (
	CG  AppType = "10"
	V2X AppType = "15"
	UAV AppType = "30"
)

type ResultCounter map[AppType]int

type ExperimentType string

const (
	ExpOptimal      ExperimentType = "optimal"
	ExpHeuristic    ExperimentType = "heuristic"
	ExpEarHeuristic ExperimentType = "ear-heuristic"
	ExpMLMasked     ExperimentType = "ml-masked"
	ExpMLNonMasked  ExperimentType = "ml-non-masked"
	ExpNotExists    ExperimentType = ""
)

func GetExpTypes() []ExperimentType {
	return []ExperimentType{ExpOptimal, ExpHeuristic, ExpEarHeuristic, ExpMLMasked, ExpMLNonMasked}
}

type AppCounter struct {
	Cg  int `json:"cg"`
	V2x int `json:"v2x"`
	Uav int `json:"uav"`
}

func (ac *AppCounter) GetTotal() int {
	return ac.Cg + ac.V2x + ac.Uav
}

func (ac *AppCounter) GetTotalAsString() string {
	return strconv.Itoa(ac.Cg + ac.V2x + ac.Uav)
}

type ExpResult struct {
	Metadata ExpResultsMeta `json:"metadata"`
	Data     ExpResultsData `json:"data"`
}

type ExpResultsMeta struct {
	ExperimentId int            `json:"experiment-id,omitempty"`
	Type         ExperimentType `json:"type"`
	Strategy     string         `json:"strategy"`
	Apps         AppCounter     `json:"apps-number"`
	Movements    int            `json:"movements"`
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
	MecHostsResults []MecHostResults `json:"mec-hosts"`
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

func getHttpRespBody(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		getErr := fmt.Errorf("HTTP GET failed for URL %s.\nError: %s\n",
			url, err)
		fmt.Fprintf(os.Stderr, getErr.Error())
		return nil, getErr
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		getErr := fmt.Errorf("HTTP GET returned status code %s for URL %s.\n",
			resp.Status, url)
		fmt.Fprintf(os.Stderr, getErr.Error())
		return nil, getErr
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln()
	}

	return b, nil
}
