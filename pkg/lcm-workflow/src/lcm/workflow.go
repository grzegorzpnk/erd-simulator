package lcm

import (
	"fmt"
	"os"
	"time"

	eta "gitlab.com/project-emco/core/emco-base/src/workflowmgr/pkg/emcotemporalapi"
	wf "go.temporal.io/sdk/workflow"
)

// special name that matches all activities
const ALL_ACTIVITIES = "all-activities"

// NeededParams should be treated as a const
var NeededParams = []string{ // parameters needed for this workflow
	"emcoOrchEndpoint", "project", "compositeApp", "compositeAppVersion", "deploymentIntentGroup", "targetAppName",
	"innotUrl", "plcControllerUrl", "rClientName", "rClientPort", "rWfClientName"}

// LcmWorkflow is a Temporal workflow that should be run after application instantiation
// It should listen for different types of notifications and serve them. For now it listens
// for CELL_ID_CHANGED notifications, calls ER Plc Controller and then invokes ER Workflow
func LcmWorkflow(ctx wf.Context, wfParam *eta.WorkflowParams) (*WorkflowParams, error) {
	// List all activities for this workflow
	activityNames := []string{
		"SubCellChangedNotification",
		"GetCellChangedNotification",
		"GenerateSmartPlacementIntent",
		"CallPlacementController",
		"GenerateRelocateWfIntent",
		"CallTemporalWfController",
	}

	// Set current state and define workflow queries
	currentState := "started" // name of ongoing activity, "started" or "completed"
	queryType := "current-state"
	err := wf.SetQueryHandler(ctx, queryType, func() (string, error) {
		return currentState, nil
	})
	if err != nil {
		currentState = "failed to register current state query handler"
		return nil, err
	}

	// Check that we got "all-activities" params
	allActivitiesParams, ok := wfParam.ActivityParams[ALL_ACTIVITIES]
	if !ok {
		err := fmt.Errorf("LcmWorkflow: expect %s parameters", ALL_ACTIVITIES)
		_, _ = fmt.Fprintf(os.Stderr, err.Error())
		return nil, err
	}

	if err := validateParams(allActivitiesParams); err != nil {
		return nil, err
	}

	// Print activity options from the workflow parameters.
	optsMap := wfParam.ActivityOpts
	actsWithOpts := make([]string, 0, len(optsMap))
	for actName := range optsMap {
		actsWithOpts = append(actsWithOpts, actName)
	}
	fmt.Printf("LcmWorkflow: got activity options for %#v\n", actsWithOpts)

	// Create a separate context for each activity based on its activity options.
	ctxMap, err := getActivityContextMap(ctx, activityNames, optsMap)
	if err != nil {
		return nil, err
	}

	wfParams := WorkflowParams{InParams: allActivitiesParams}

	currentState = "sub-cell-changed-notification"
	ctx1 := ctxMap["SubCellChangedNotification"]
	err = wf.ExecuteActivity(ctx1, SubCellChangedNotification, wfParams).Get(ctx1, &wfParams)
	if err != nil {
		wferr := fmt.Errorf("SubCellChangedNotification failed: %s", err.Error())
		_, _ = fmt.Fprintf(os.Stderr, wferr.Error())
		return nil, wferr
	}

	for {
		currentState = "get-cell-changed-notification"
		ctx2 := ctxMap["GetCellChangedNotification"]
		err = wf.ExecuteActivity(ctx2, GetCellChangedNotification, wfParams).Get(ctx2, &wfParams)
		if err != nil {
			wferr := fmt.Errorf("GetCellChangedNotification failed: %s", err.Error())
			_, _ = fmt.Fprintf(os.Stderr, wferr.Error())
			return nil, wferr
		}

		currentState = "generate-smart-placement-intent"
		ctx3 := ctxMap["GenerateSmartPlacementIntent"]
		err = wf.ExecuteActivity(ctx3, GenerateSmartPlacementIntent, wfParams).Get(ctx3, &wfParams)
		if err != nil {
			wferr := fmt.Errorf("GenerateSmartPlacementIntent failed: %s", err.Error())
			fmt.Fprintf(os.Stderr, wferr.Error())
			return nil, wferr
		}
		currentState = "call-placement-controller"
		ctx4 := ctxMap["CallPlacementController"]
		err = wf.ExecuteActivity(ctx4, CallPlacementController, wfParams).Get(ctx4, &wfParams)
		if err != nil {
			wferr := fmt.Errorf("CallPlacementController failed: %s", err.Error())
			fmt.Fprintf(os.Stderr, wferr.Error())
			// if we fail in this activity, just continue
			continue
			//return nil, wferr
		}

		currentState = "generate-relocate-wf-intent"
		ctx5 := ctxMap["GenerateRelocateWfIntent"]
		err = wf.ExecuteActivity(ctx5, GenerateRelocateWfIntent, wfParams).Get(ctx5, &wfParams)
		if err != nil {
			wferr := fmt.Errorf("CreateTemporalERIntent failed: %s", err.Error())
			fmt.Fprintf(os.Stderr, wferr.Error())
			return nil, wferr
		}

		currentState = "call-temporal-wf-controller"
		ctx6 := ctxMap["CallTemporalWfController"]
		err = wf.ExecuteActivity(ctx6, CallTemporalWfController, wfParams).Get(ctx6, &wfParams)
		if err != nil {
			wferr := fmt.Errorf("CallTemporalWfController failed: %s", err.Error())
			fmt.Fprintf(os.Stderr, wferr.Error())
			return nil, wferr
		}

	}
	currentState = "completed"

	//fmt.Printf("After all activities: wfParams = %#v\n", wfParams)
	fmt.Printf("Workflow ended successfully!\n")

	return &wfParams, nil
}

// getActivityContextMap returns a list of Temporal contexts for each activity.
// Note that this is generic code that is independent of user's app/workflows.
func getActivityContextMap(ctx wf.Context, activityNames []string,
	optsMap map[string]wf.ActivityOptions) (map[string]wf.Context, error) {

	// Validate that all activity names in given workflow params are valid
	allActivitiesFlag := false
	for paramActName := range optsMap {
		found := false
		for _, actName := range activityNames {
			if paramActName == actName {
				found = true
				break
			}
			if paramActName == ALL_ACTIVITIES {
				found = true
				allActivitiesFlag = true
				break
			}
		}
		if !found {
			err := fmt.Errorf("invalid activity name in params: %s", paramActName)
			return nil, err
		}
	}

	// Init each activity-specific context to default or the specified param
	defaultActivityOpts := wf.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctxMap := make(map[string]wf.Context, len(activityNames))
	for _, actName := range activityNames {
		// init to default
		ctxMap[actName] = wf.WithActivityOptions(ctx, defaultActivityOpts)
		// Apply all-activities options if specified
		if allActivitiesFlag {
			fmt.Printf("Applying all-activities options for activity %s\n", actName)
			ctxMap[actName] = wf.WithActivityOptions(ctx, optsMap[ALL_ACTIVITIES])
		}
		// Apply activity-specific options, if specified
		for paramActName := range optsMap {
			if paramActName == actName {
				fmt.Printf("Applying activity-specifc options for %s\n", actName)
				ctxMap[actName] = wf.WithActivityOptions(ctx, optsMap[actName])
			}
		}
	}
	return ctxMap, nil
}

// validateParams verifies that inParams has all needed params for this workflow
func validateParams(inParams map[string]string) error {

	var paramsNotFound []string
	for _, neededParam := range NeededParams {
		found := false
		for inParam := range inParams {
			if neededParam == inParam {
				found = true
			}
		}
		if !found {
			paramsNotFound = append(paramsNotFound, neededParam)
		}
	}

	if len(paramsNotFound) > 0 {
		err := fmt.Errorf("Workflow needs these params: %#v\n", paramsNotFound)
		_, _ = fmt.Fprintf(os.Stderr, err.Error())
		return err
	}

	return nil
}
