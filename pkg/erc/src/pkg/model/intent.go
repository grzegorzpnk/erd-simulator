package model

type CellId string

// CurrentPlacement represents where the application is currently instantiated
type CurrentPlacement struct {
	Provider        string `json:"provider"`
	Cluster         string `json:"cluster"`
	ClusterCapacity int    `json:"omitempty,clusterCapacity"`
}

// SmartPlacementIntent defines the Intent to perform Smart Placement for giver application
type SmartPlacementIntent struct {
	Metadata         Metadata                 `json:"metadata"`
	CurrentPlacement CurrentPlacement         `json:"currentPlacement"`
	Spec             SmartPlacementIntentSpec `json:"spec"`
}

type Metadata struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description,omitempty" yaml:"-"`
	UserData1   string `json:"userData1,omitempty" yaml:"-"`
	UserData2   string `json:"userData2,omitempty" yaml:"-"`
}

type SmartPlacementIntentSpec struct {
	AppName                  string                     `json:"app"`
	SmartPlacementIntentData SmartPlacementIntentStruct `json:"smartPlacementIntentData"`
}

type SmartPlacementIntentStruct struct {
	TargetCell        CellId      `json:"targetCell"`
	AppCpuReq         float64     `json:"appCpuReq"`
	AppMemReq         float64     `json:"appMemReq"`
	ConstraintsList   Constraints `json:"constraintsList"`
	ParametersWeights Weights     `json:"parametersWeights,omitempty"`
}

type Constraints struct {
	LatencyMax float64 `json:"latencyMax"`
}

type Weights struct {
	LatencyWeight        float64 `json:"latencyWeight"`
	ResourcesWeight      float64 `json:"resourcesWeight"`
	CpuUtilizationWeight float64 `json:"cpuUtilizationWeight"`
	MemUtilizationWeight float64 `json:"memUtilizationWeight"`
}

type SmartPlacementIntentKey struct {
	Project               string `json:"project"`
	CompositeApp          string `json:"compositeApp"`
	CompositeAppVersion   string `json:"compositeAppVersion"`
	DeploymentIntentGroup string `json:"deploymentIntentGroup"`
	SmartPlacementIntent  string `json:"smartPlacementIntent"`
}
