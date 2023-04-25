package model

type MLSmartPlacementIntent struct {
	State       State `json:"state"`
	CurrentMask Mask  `json:"mask,omitempty"`
} // SpaceAPP (for single app) : 1) Required mvCPU 2) required Memory 3) Required Latency 4) Current MEC 5) Current RAN
type State struct {
	SpaceAPP  [1][5]int `json:"space_App"`
	SpaceMECs [][]int   `json:"space_MEC"`
}

type Mask []int