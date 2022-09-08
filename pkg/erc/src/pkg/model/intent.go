// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

// Package model contains the data model necessary for the implementations.
// In this example, SmartPlacementIntent
package model

type AppPriority int
type CellId string

const (
	PRIORITY_LOW AppPriority = iota
	PRIORITY_NORMAL
	PRIORITY_IMPORTANT
	PRIORITY_CRITICAL
)

// SmartPlacementIntent defines the high level structure of a network chain document
type SmartPlacementIntent struct {
	Metadata Metadata                 `json:"metadata" yaml:"metadata"`
	Spec     SmartPlacementIntentSpec `json:"spec" yaml:"spec"`
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
	PriorityLevel     AppPriority `json:"priorityLevel"`
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