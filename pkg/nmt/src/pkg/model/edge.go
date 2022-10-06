package model

import "10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/metrics"

type Edge struct {
	SourceVertexName         string                 `json:"source"`
	SourceVertexProviderName string                 `json:"sourceVertexProviderName,omitempty"`
	TargetVertexName         string                 `json:"target"`
	TargetVertexProviderName string                 `json:"targetVertexProviderName,omitempty"`
	EdgeMetrics              metrics.NetworkMetrics `json:"edgeMetrics"`
}
