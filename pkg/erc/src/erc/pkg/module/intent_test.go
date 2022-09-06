// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

// These test cases are to validate the module level functionalities.
package module_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"10.254.188.33/matyspi5/erd/pkg/erc/pkg/model"
	"10.254.188.33/matyspi5/erd/pkg/erc/pkg/module"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/db"
	orchmodule "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/module"
)

var _ = Describe("ErcIntent",
	func() {
		var (
			project   orchmodule.Project
			pClient   *orchmodule.ProjectClient
			app       orchmodule.CompositeApp
			aClient   *orchmodule.CompositeAppClient
			diGroup   orchmodule.DeploymentIntentGroup
			digClient *orchmodule.DeploymentIntentGroupClient
			intent    model.SmartPlacementIntent
			iClient   *module.SmartPlacementIntentClient
			mdb       *db.NewMockDB
		)

		BeforeEach(
			func() {
				pClient = orchmodule.NewProjectClient()
				project = orchmodule.Project{
					MetaData: orchmodule.ProjectMetaData{
						Name: "testProj",
					},
				}

				aClient = orchmodule.NewCompositeAppClient()
				app = orchmodule.CompositeApp{
					Metadata: orchmodule.CompositeAppMetaData{
						Name: "app",
					},
					Spec: orchmodule.CompositeAppSpec{
						Version: "v1",
					},
				}

				digClient = orchmodule.NewDeploymentIntentGroupClient()
				diGroup = orchmodule.DeploymentIntentGroup{
					MetaData: orchmodule.DepMetaData{
						Name: "diGroup",
					},
					Spec: orchmodule.DepSpecData{
						Profile:      "profilename",
						Version:      "testver",
						LogicalCloud: "logCloud",
					},
				}

				iClient = module.NewIntentClient()
				intent = model.SmartPlacementIntent{
					Metadata: model.Metadata{
						Name: "smartPlacementIntentName",
					},
				}

				mdb = new(db.NewMockDB)
				mdb.Err = nil
				db.DBconn = mdb
			},
		)

		Describe("Create  SmartPlacementIntent",
			func() {
				It("successful creation of  SmartPlacementIntent",
					func() {
						// set up prerequisites
						_, err := (*pClient).CreateProject(project, false)
						Expect(err).To(BeNil())
						_, err = (*aClient).CreateCompositeApp(app, "testProj", false)
						Expect(err).To(BeNil())
						_, _, err = (*digClient).CreateDeploymentIntentGroup(diGroup, "testProj", "app", "v1", true)
						Expect(err).To(BeNil())

						// test  intent creation
						_, err = (*iClient).CreateSmartPlacementIntent(intent, "testProj", "app", "v1", "diGroup", true)
						Expect(err).To(BeNil())
					},
				)
			},
		)

		Describe("Get  SmartPlacementIntent",
			func() {
				It("successful get of SmartPlacementIntent",
					func() {
						// set up prerequisites
						_, err := (*pClient).CreateProject(project, false)
						Expect(err).To(BeNil())
						_, err = (*aClient).CreateCompositeApp(app, "testProj", false)
						Expect(err).To(BeNil())
						_, _, err = (*digClient).CreateDeploymentIntentGroup(diGroup, "testProj", "app", "v1", true)
						Expect(err).To(BeNil())

						// test  intent creation
						_, err = (*iClient).CreateSmartPlacementIntent(intent, "testProj", "app", "v1", "diGroup", false)
						Expect(err).To(BeNil())

						_, err = (*iClient).GetSmartPlacementIntent("smartPlacementIntentName", "testProj", "app", "v1", "diGroup")
						Expect(err).To(BeNil())
					},
				)
			},
		)
	},
)
