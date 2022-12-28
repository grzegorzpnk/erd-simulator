package model

type MECApp struct {
	Id           string             `json:"id"`
	Requirements RequestedResources `json:"requirements"`
	ClusterId    string             `json:"clusterId"`
}

type RequestedResources struct {
	RequestedCPU     float64 `json:"requestedCPU"`
	RequestedMEMORY  float64 `json:"requestedMEMORY"`
	RequestedLatency float64 `json:"requestedLatency"`
}
