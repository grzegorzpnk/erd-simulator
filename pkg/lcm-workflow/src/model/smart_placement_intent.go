package model

import "10.254.188.33/matyspi5/erd/pkg/lcm-workflow/src/types"

type SmartPlacementIntent struct {
	Metadata Metadata                 `json:"metadata,omitempty"`
	Spec     SmartPlacementIntentSpec `json:"spec,omitempty"`
}

type Metadata struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	UserData1   string `json:"userData1,omitempty"`
	UserData2   string `json:"userData2,omitempty"`
}

type SmartPlacementIntentSpec struct {
	AppName                  string                     `json:"app"`
	SmartPlacementIntentData SmartPlacementIntentStruct `json:"smartPlacementIntentData"`
}

type SmartPlacementIntentStruct struct {
	TargetCell        types.CellId `json:"targetCell"`
	AppCpuReq         float64      `json:"appCpuReq"`
	AppMemReq         float64      `json:"appMemReq"`
	ConstraintsList   Constraints  `json:"constraintsList,omitempty"`
	ParametersWeights Weights      `json:"parametersWeights,omitempty"`
}

type Constraints struct {
	LatencyMax        float64 `json:"latencyMax,omitempty"`
	CpuUtilizationMax float64 `json:"cpuUtilizationMax,omitempty"`
	MemUtilizationMax float64 `json:"memUtilizationMax,omitempty"`
}

type Weights struct {
	LatencyWeight        float64 `json:"latencyWeight"`
	ResourcesWeight      float64 `json:"resourcesWeight"`
	CpuUtilizationWeight float64 `json:"cpuUtilizationWeight"`
	MemUtilizationWeight float64 `json:"memUtilizationWeight"`
}
