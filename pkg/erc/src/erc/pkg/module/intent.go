// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package module

import (
	"10.254.188.33/matyspi5/erd/pkg/erc/pkg/model"
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
