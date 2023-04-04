package model

type MLSmartPlacementIntent struct {
	StateApp  SpaceAPP  `json:"space_App"`
	StateMECS SpaceMECs `json:"space_MEC"`
}

// SpaceAPP (for single app)  : 1) Required mvCPU 2) required Memory 3) Required Latency 4) Current MEC 5) Current RAN
type SpaceAPP struct {
	AppCharacteristic [5]int64 `json:"app_characteristic"`
}

type SpaceMECs struct {
	MecCharacteristics [][]int `json:"mec_characteristics"`
}
