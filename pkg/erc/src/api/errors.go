// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package api

import (
	"net/http"

	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/apierror"
)

var apiErrors = []apierror.APIError{
	{ID: "SmartPlacementIntent not found", Message: "SmartPlacementIntent not found", Status: http.StatusNotFound},
	{ID: "SmartPlacementIntent already exists", Message: "SmartPlacementIntent already exists", Status: http.StatusConflict},
}
