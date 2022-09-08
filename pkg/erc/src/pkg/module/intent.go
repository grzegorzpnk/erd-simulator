// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package module

import (
	topology "10.254.188.33/matyspi5/erd/pkg/erc/mock"
	"10.254.188.33/matyspi5/erd/pkg/erc/pkg/model"
	"fmt"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/db"

	"github.com/pkg/errors"
)

// SmartPlacementIntentClient implements the SmartPlacementIntentManager.
// It will also be used to maintain some localized state.
type SmartPlacementIntentClient struct {
	dbInfo DBInfo
}

func NewIntentClient() *SmartPlacementIntentClient {
	return &SmartPlacementIntentClient{
		dbInfo: DBInfo{
			collection: "resources", // should remain the same
			tag:        "data",      // should remain the same
		},
	}
}

// SmartPlacementIntentManager a manager is an interface for exposing the client's functionalities.
// You can have multiple managers based on the requirement and its implementation.
// In this example, SmartPlacementIntentManager exposes the SmartPlacementIntentClient functionalities.
type SmartPlacementIntentManager interface {
	CreateSmartPlacementIntent(intent model.SmartPlacementIntent, project, app, version, deploymentIntentGroup string, failIfExists bool) (model.SmartPlacementIntent, error)
	GetSmartPlacementIntent(name, project, app, version, deploymentIntentGroup string) ([]model.SmartPlacementIntent, error)
	ServeSmartPlacementIntentOutsideEMCO(intent model.SmartPlacementIntent) (string, string, error)
}

func (i *SmartPlacementIntentClient) ServeSmartPlacementIntentOutsideEMCO(intent model.SmartPlacementIntent) (string, string, error) {
	intentData := intent.Spec.SmartPlacementIntentData

	c := intentData.ConstraintsList
	pw := intentData.ParametersWeights
	fmt.Printf("intent: %+v\n", intent)

	// x1 -> latency 			| weight = w1
	// x2 -> cpu utilization 	| weight = w2
	// x3 -> mem utilization	| weight = w3
	// minimize w1 * x1 + w2 * x2 + w3 * x3
	// x1 < latencyMax
	// x2 < cpuUtilMax
	// x3 < memUtilMax

	constraintsMet := func(latency, cpuUtilization, memUtilization float64) bool {
		if latency > c.LatencyMax {
			return false
		} else if cpuUtilization > c.CpuUtilizationMax {
			return false
		} else if memUtilization > c.MemUtilizationMax {
			return false
		} else {
			return true
		}
	}
	normalize := func(latency, cpuUtilization, memUtilization float64) (float64, float64, float64) {
		return latency / c.LatencyMax, cpuUtilization / c.CpuUtilizationMax, memUtilization / c.MemUtilizationMax
	}

	compute := func(latency, cpuUtilization, memUtilization float64) (float64, bool) {
		if !constraintsMet(latency, cpuUtilization, memUtilization) {
			return 3, false
		}
		latency, cpuUtilization, memUtilization = normalize(latency, cpuUtilization, memUtilization)

		return pw.LatencyWeight*latency + pw.CpuUtilizationWeight*cpuUtilization + pw.MemUtilizationWeight*memUtilization, true
	}

	MecHosts := topology.GetMecHostsByCellId(intentData.TargetCell)

	bestMecHost := MecHosts[0]
	var bestOk bool
	var best float64
	for _, mecHost := range MecHosts {
		best, bestOk = compute(bestMecHost.Latency, bestMecHost.CpuUtilization, bestMecHost.MemUtilization)
		current, currentOk := compute(mecHost.Latency, mecHost.CpuUtilization, mecHost.MemUtilization)

		if currentOk && (best > current) {
			bestMecHost = mecHost
		}
	}
	if !bestOk {
		return "", "", errors.New("No clusters met all constraints!")
	}
	return bestMecHost.ProviderName, bestMecHost.ClusterName, nil
}

// CreateSmartPlacementIntent insert a new SmartPlacementIntent in the database
func (i *SmartPlacementIntentClient) CreateSmartPlacementIntent(intent model.SmartPlacementIntent, project, app, version, deploymentIntentGroup string, failIfExists bool) (model.SmartPlacementIntent, error) {
	// Construct key and tag to select the entry.
	key := model.SmartPlacementIntentKey{
		Project:               project,
		CompositeApp:          app,
		CompositeAppVersion:   version,
		DeploymentIntentGroup: deploymentIntentGroup,
		SmartPlacementIntent:  intent.Metadata.Name,
	}

	// Check if this SmartPlacementIntent already exists.
	intents, err := i.GetSmartPlacementIntent(intent.Metadata.Name, project, app, version, deploymentIntentGroup)
	if err == nil &&
		len(intents) > 0 &&
		intents[0].Metadata.Name == intent.Metadata.Name &&
		failIfExists {
		return model.SmartPlacementIntent{}, errors.New("SmartPlacementIntent already exists")
	}

	err = db.DBconn.Insert(i.dbInfo.collection, key, nil, i.dbInfo.tag, intent)
	if err != nil {
		return model.SmartPlacementIntent{}, err
	}

	return intent, nil
}

// GetSmartPlacementIntent returns the SmartPlacementIntent for the corresponding name
func (i *SmartPlacementIntentClient) GetSmartPlacementIntent(name, project, app, version, deploymentIntentGroup string) ([]model.SmartPlacementIntent, error) {
	// Construct key and tag to select the entry.
	key := model.SmartPlacementIntentKey{
		Project:               project,
		CompositeApp:          app,
		CompositeAppVersion:   version,
		DeploymentIntentGroup: deploymentIntentGroup,
		SmartPlacementIntent:  name,
	}

	values, err := db.DBconn.Find(i.dbInfo.collection, key, i.dbInfo.tag)
	if err != nil {
		return []model.SmartPlacementIntent{}, err
	}

	if len(values) == 0 {
		return []model.SmartPlacementIntent{}, errors.New("SmartPlacementIntent not found")
	}

	intents := []model.SmartPlacementIntent{}

	for _, v := range values {
		i := model.SmartPlacementIntent{}
		err = db.DBconn.Unmarshal(v, &i)
		if err != nil {
			return []model.SmartPlacementIntent{}, errors.Wrap(err, "Unmarshalling Value")
		}

		intents = append(intents, i)
	}

	return intents, nil
}
