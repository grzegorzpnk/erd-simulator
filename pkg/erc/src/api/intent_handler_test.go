// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

// These test cases are to validate the route handler functionalities.
package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"10.254.188.33/matyspi5/erd/pkg/erc/api"
	"10.254.188.33/matyspi5/erd/pkg/erc/pkg/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

// In this example, the SmartPlacementIntentManager exposes two functionalities (CreateSmartPlacementIntent, GetIntent).
// mockSmartPlacementIntentManager implements the mock services for the SmartPlacementIntentManager.
type mockIntentManager struct {
	Items []model.SmartPlacementIntent
	Err   error
}

type test struct {
	name       string
	input      io.Reader
	intent     model.SmartPlacementIntent
	err        error
	statusCode int
	client     *mockIntentManager
}

func init() {
	api.ErJSONFile = "../json-schemas/intent.json"
}

func (m *mockIntentManager) CreateSmartPlacementIntent(intent model.SmartPlacementIntent, project, app, version, deploymentIntentGroup string, failIfExists bool) (model.SmartPlacementIntent, error) {
	if m.Err != nil {
		return model.SmartPlacementIntent{}, m.Err
	}

	return m.Items[0], nil
}

func (m *mockIntentManager) GetSmartPlacementIntent(name, project, app, version, deploymentIntentGroup string) ([]model.SmartPlacementIntent, error) {
	if m.Err != nil {
		return []model.SmartPlacementIntent{}, m.Err
	}

	return m.Items, nil
}

var _ = Describe("IntentHandler",
	func() {
		DescribeTable("Create SmartPlacementIntent",
			func(t test) {
				i := model.SmartPlacementIntent{}
				req := httptest.NewRequest("POST", "/v2/projects/test-project/composite-apps/test-compositeapp/v1/deployment-intent-groups/test-dig/smartPlacementIntents", t.input)
				res := executeRequest(req, api.NewRouter(t.client))
				Expect(res.StatusCode).To(Equal(t.statusCode))
				json.NewDecoder(res.Body).Decode(&i)
				Expect(i).To(Equal(t.intent))
			},
			Entry("successful create",
				test{
					name:       "create",
					statusCode: http.StatusCreated,
					input: bytes.NewBuffer([]byte(`{
						"metadata": {
							"name": "testsmartPlacementIntent",
							"description": "test intent",
							"userData1": "some user data 1",
							"userData2": "some user data 2"
						},
						"spec": {
							"app": "testapp",
							"smartPlacementIntentData": "testIntentData"
						}
					}`)),
					intent: model.SmartPlacementIntent{
						Metadata: model.Metadata{
							Name:        "testSmartPlacementIntent",
							Description: "test intent",
							UserData1:   "some user data 1",
							UserData2:   "some user data 2",
						},
						Spec: model.SmartPlacementIntentSpec{
							AppName:                  "testApp",
							SmartPlacementIntentData: model.SmartPlacementIntentStruct{},
						},
					},
					err: nil,
					client: &mockIntentManager{
						Err:   nil,
						Items: populateTestData(),
					},
				},
			),
			// Add more entries to cover multiple create success/ failure scenarios.
		)
		DescribeTable("Get SmartPlacementIntent",
			func(t test) {
				i := model.SmartPlacementIntent{}
				req := httptest.NewRequest("GET", "/v2/projects/test-project/composite-apps/test-compositeapp/v1/deployment-intent-groups/test-dig/smartPlacementIntents/"+t.name, nil)
				res := executeRequest(req, api.NewRouter(t.client))
				Expect(res.StatusCode).To(Equal(t.statusCode))
				json.NewDecoder(res.Body).Decode(&i)
				Expect(i).To(Equal(t.intent))
			},
			Entry("successful get",
				test{
					name:       "get",
					statusCode: http.StatusOK,
					err:        nil,
					intent: model.SmartPlacementIntent{
						Metadata: model.Metadata{
							Name:        "testSmartPlacementIntent",
							Description: "test intent",
							UserData1:   "some user data 1",
							UserData2:   "some user data 2",
						},
						Spec: model.SmartPlacementIntentSpec{
							AppName:                  "testApp",
							SmartPlacementIntentData: model.SmartPlacementIntentStruct{},
						},
					},
					client: &mockIntentManager{
						Err:   nil,
						Items: populateTestData(),
					},
				},
			),
			// Add more entries to cover multiple get success/ failure scenarios.
		)
		// Add more tests based on the handler functionalities.
	},
)

func populateTestData() []model.SmartPlacementIntent {
	return []model.SmartPlacementIntent{
		{
			Metadata: model.Metadata{
				Name:        "testSmartPlacementIntent",
				Description: "test intent",
				UserData1:   "some user data 1",
				UserData2:   "some user data 2",
			},
			Spec: model.SmartPlacementIntentSpec{
				AppName:                  "testApp",
				SmartPlacementIntentData: model.SmartPlacementIntentStruct{},
			},
		},
		{
			Metadata: model.Metadata{
				Name:        "newSmartPlacementIntent",
				Description: "new intent",
				UserData1:   "some user data 1",
				UserData2:   "some user data 2",
			},
			Spec: model.SmartPlacementIntentSpec{
				AppName:                  "newApp",
				SmartPlacementIntentData: model.SmartPlacementIntentStruct{},
			},
		},
		// Add more data based on the test scenarios.
	}
}
