// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package module

import (
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/model"
	log "github.com/sirupsen/logrus"

	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/topology"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/db"
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
	var bestOk bool
	var best float64
	tClient := topology.NewTopologyClient()

	fmt.Printf("Smart Placement Intent: %+v\n", intent)

	targetCell := intent.Spec.SmartPlacementIntentData.TargetCell

	mecHosts, err := tClient.GetMecHostsByCellId(targetCell)
	if err != nil {
		log.Errorf("could not serve Smart Placement Intent: %v", err)
		return "", "", err
	}

	if len(mecHosts) <= 0 {
		return "", "", errors.New(fmt.Sprintf("no mec hosts found for given cell: %v\n.", targetCell))
	}

	bestMecHost := mecHosts[0]
	for _, mecHost := range mecHosts {
		best, bestOk = ComputeObjectiveValue(intent, bestMecHost)
		current, currentOk := ComputeObjectiveValue(intent, mecHost)

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

func ComputeObjectiveValue(i model.SmartPlacementIntent, mec topology.MecHost) (float64, bool) {
	cl := i.Spec.SmartPlacementIntentData.ConstraintsList
	pw := i.Spec.SmartPlacementIntentData.ParametersWeights

	if !CheckConstraintsForGiven(cl, mec) {
		return 3, false
	}
	latency, cpuUtilization, memUtilization := NormalizeMecParameters(cl, mec)

	return pw.LatencyWeight*latency + pw.CpuUtilizationWeight*cpuUtilization + pw.MemUtilizationWeight*memUtilization, true
}

func CheckConstraintsForGiven(cl model.Constraints, mec topology.MecHost) bool {
	if mec.Latency > cl.LatencyMax {
		return false
	} else if mec.CpuUtilization > cl.CpuUtilizationMax {
		return false
	} else if mec.MemUtilization > cl.MemUtilizationMax {
		return false
	} else {
		return true
	}
}

func NormalizeMecParameters(cl model.Constraints, mec topology.MecHost) (float64, float64, float64) {
	return mec.Latency / cl.LatencyMax, mec.CpuUtilization / cl.CpuUtilizationMax, mec.MemUtilization / cl.MemUtilizationMax
}
