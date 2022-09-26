package topology

type Edge struct {
	SourceVertexName         string         `json:"source"`
	SourceVertexProviderName string         `json:"sourceVertexProviderName,omitempty"`
	TargetVertexName         string         `json:"target"`
	TargetVertexProviderName string         `json:"targetVertexProviderName,omitempty"`
	EdgeMetrics              NetworkMetrics `json:"edgeMetrics"`
}
