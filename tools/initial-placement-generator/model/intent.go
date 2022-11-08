// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

// Package model contains the data model necessary for the implementations.
// In this example, SmartPlacementIntent
package model

type CellId string

type Cluster struct {
	Provider string `json:"provider"`
	Cluster  string `json:"cluster"`
}

// SmartPlacementIntent defines the high level structure of a network chain document
type SmartPlacementIntent struct {
	Metadata         Metadata `json:"metadata,omitempty"`
	CurrentPlacement Cluster
	Spec             SmartPlacementIntentSpec `json:"spec,omitempty"`
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
	LatencyMax        float64 `json:"latencyMax"`
	CpuUtilizationMax float64 `json:"cpuUtilizationMax"`
	MemUtilizationMax float64 `json:"memUtilizationMax"`
}

type Weights struct {
	LatencyWeight        float64 `json:"latencyWeight"`
	ResourcesWeight      float64 `json:"resourcesWeight"`
	CpuUtilizationWeight float64 `json:"cpuUtilizationWeight"`
	MemUtilizationWeight float64 `json:"memUtilizationWeight"`
}

// SmartPlacementIntentKey is the key structure that is used in the database
type SmartPlacementIntentKey struct {
	Project               string `json:"project"`
	CompositeApp          string `json:"compositeApp"`
	CompositeAppVersion   string `json:"compositeAppVersion"`
	DeploymentIntentGroup string `json:"deploymentIntentGroup"`
	SmartPlacementIntent  string `json:"smartPlacementIntent"`
}
