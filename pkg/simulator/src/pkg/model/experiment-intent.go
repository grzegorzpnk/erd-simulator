package model

import (
	"strconv"
)

type ExperimentType string

const (
	ExpOptimal      ExperimentType = "optimal"
	ExpHeuristic    ExperimentType = "heuristic"
	ExpEarHeuristic ExperimentType = "ear-heuristic"
	ExpMLMasked     ExperimentType = "ml-masked"
	ExpMLNonMasked  ExperimentType = "ml-non-masked"
	ExpNotExists    ExperimentType = ""
)

var ExperimentTypes []ExperimentType = []ExperimentType{
	ExpOptimal,
	ExpHeuristic,
	ExpEarHeuristic,
	ExpMLMasked,
	ExpMLNonMasked,
}

type ExperimentStrategy string

const (
	StrLatency ExperimentStrategy = "latency"
	StrLB      ExperimentStrategy = "load-balancing"
	StrHybrid  ExperimentStrategy = "hybrid"
	StrML      ExperimentStrategy = "ml"
	StrML2     ExperimentStrategy = "ml2"
	StrML3     ExperimentStrategy = "ml3"
	StrML4     ExperimentStrategy = "ml4"
	Str7L3R    ExperimentStrategy = "7l3r"
	Str3L7R    ExperimentStrategy = "3l7r"
)

type AppType string

const (
	CG  AppType = "10"
	V2X AppType = "15"
	UAV AppType = "30"
)

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

type ExperimentIntent struct {
	ExperimentType     ExperimentType     `json:"experiment-type"`
	ExperimentStrategy ExperimentStrategy `json:"experiment-strategy"`
	ExperimentDetails  ExperimentDetails  `json:"experiment-details"`
	Weights            Weights            `json:"Weights,omitempty"`
}

type ExperimentDetails struct {
	MovementsInExperiment string     `json:"number-of-movements"`
	InitialAppsNumber     AppCounter `json:"initial-apps-number"`
}
