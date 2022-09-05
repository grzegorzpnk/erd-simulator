package module

type AppPriority int

const (
	PRIORITY_LOW AppPriority = iota
	PRIORITY_NORMAL
	PRIORITY_IMPORTANT
	PRIORITY_CRITICAL
)

type ErIntent struct {
	MetaData MetaData `json:"metadata,omitempty"`
	Spec     SpecData `json:"spec,omitempty"`
}

type MetaData struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	UserData1   string `json:"userData1,omitempty"`
	UserData2   string `json:"userData2,omitempty"`
}

type SpecData struct {
	AppName string       `json:"app,omitempty"`
	Intent  IntentStruct `json:"intent,omitempty"`
}

type IntentStruct struct {
	PriorityLevel     AppPriority `json:"priority-level,omitempty"`
	ConstraintsList   Constraints `json:"constraints-list,omitempty"`
	ParametersWeights Weights     `json:"parameters-weights,omitempty"`
}

type Constraints struct {
	LatencyMax        float64 `json:"latency-max,omitempty"`
	CpuUtilizationMax float64 `json:"cpu-utilization-max,omitempty"`
	MemUtilizationMax float64 `json:"mem-utilization-max,omitempty"`
}

type Weights struct {
	LatencyWeight        float64
	CpuUtilizationWeight float64
	MemUtilizationWeight float64
}
