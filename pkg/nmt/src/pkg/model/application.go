package model

type MECApp struct {
	Id           string             `json:"id"`
	Requirements RequestedResources `json:"requirements"`
}

type RequestedResources struct {
	RequestedCPU    float64 `json:"requestedCPU"`
	RequestedMEMORY float64 `json:"requestedMEMORY"`
}
