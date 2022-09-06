// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package api

import (
	"net/http"

	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/apierror"
)

// Capture the errors returned from the module.
// Format the error message and status code accordingly to reply to the request.
// In this example, send a status code of 404 if the intent is not there in the database.
// ref: https://gitlab.com/project-emco/core/emco-base/-/blob/main/src/hpa-plc/api/apierrors.go
var apiErrors = []apierror.APIError{
	{ID: "SmartPlacementIntent not found", Message: "SmartPlacementIntent not found", Status: http.StatusNotFound},
	{ID: "SmartPlacementIntent already exists", Message: "SmartPlacementIntent already exists", Status: http.StatusConflict},
}
