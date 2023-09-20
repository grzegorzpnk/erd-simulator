package results

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	log "simu/src/logger"
	"simu/src/pkg/model"
	"strings"
)

type ChartType int

const (
	RelocationRates ChartType = iota
	FailedRates
	SkippedRates
	RedundantRates
	RelocationTriggeringRates
	RelocationRejectionRates
	RelocationSuccessfulSearchRates
	ResCpu
	ResMemory
)

func (c *Client) GenerateChartPkgMecs(chartType ChartType, basePath string) error {
	var resource string

	switch chartType {
	case ResCpu:
		resource = "cpu"
	case ResMemory:
		resource = "memory"
	}

	var values = [][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}}

	values[0][0] = c.GetMecUtilizationAggregated(model.ExpOptimal, model.StrHybrid, MecLocal, resource)
	values[0][1] = c.GetMecUtilizationAggregated(model.ExpOptimal, model.StrHybrid, MecRegional, resource)
	values[0][2] = c.GetMecUtilizationAggregated(model.ExpOptimal, model.StrHybrid, MecCentral, resource)

	values[1][0] = c.GetMecUtilizationAggregated(model.ExpHeuristic, model.StrHybrid, MecLocal, resource)
	values[1][1] = c.GetMecUtilizationAggregated(model.ExpHeuristic, model.StrHybrid, MecRegional, resource)
	values[1][2] = c.GetMecUtilizationAggregated(model.ExpHeuristic, model.StrHybrid, MecCentral, resource)

	values[2][0] = c.GetMecUtilizationAggregated(model.ExpEarHeuristic, model.StrHybrid, MecLocal, resource)
	values[2][1] = c.GetMecUtilizationAggregated(model.ExpEarHeuristic, model.StrHybrid, MecRegional, resource)
	values[2][2] = c.GetMecUtilizationAggregated(model.ExpEarHeuristic, model.StrHybrid, MecCentral, resource)

	values[3][0] = c.GetMecUtilizationAggregated(model.ExpMLMasked, model.StrML, MecLocal, resource)
	values[3][1] = c.GetMecUtilizationAggregated(model.ExpMLMasked, model.StrML, MecRegional, resource)
	values[3][2] = c.GetMecUtilizationAggregated(model.ExpMLMasked, model.StrML, MecCentral, resource)

	values[4][0] = c.GetMecUtilizationAggregated(model.ExpMLNonMasked, model.StrML, MecLocal, resource)
	values[4][1] = c.GetMecUtilizationAggregated(model.ExpMLNonMasked, model.StrML, MecRegional, resource)
	values[4][2] = c.GetMecUtilizationAggregated(model.ExpMLNonMasked, model.StrML, MecCentral, resource)

	err := c.genRatesPkgAggregatedMecs(resource, values, basePath, false)
	if err != nil {
		log.Errorf("Error: %v", err)
	}

	return nil
}

func (c *Client) GenerateChartPkgMecsICC(chartType ChartType, basePath string) error {
	var resource string

	switch chartType {
	case ResCpu:
		resource = "cpu"
	case ResMemory:
		resource = "memory"
	}

	var values = [][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}}

	values[0][0] = c.GetMecUtilizationAggregated(model.ExpOptimal, model.StrLB, MecLocal, resource)
	values[0][1] = c.GetMecUtilizationAggregated(model.ExpOptimal, model.StrLB, MecRegional, resource)
	values[0][2] = c.GetMecUtilizationAggregated(model.ExpOptimal, model.StrLB, MecCentral, resource)

	values[1][0] = c.GetMecUtilizationAggregated(model.ExpOptimal, model.StrLatency, MecLocal, resource)
	values[1][1] = c.GetMecUtilizationAggregated(model.ExpOptimal, model.StrLatency, MecRegional, resource)
	values[1][2] = c.GetMecUtilizationAggregated(model.ExpOptimal, model.StrLatency, MecCentral, resource)

	values[2][0] = c.GetMecUtilizationAggregated(model.ExpOptimal, model.StrHybrid, MecLocal, resource)
	values[2][1] = c.GetMecUtilizationAggregated(model.ExpOptimal, model.StrHybrid, MecRegional, resource)
	values[2][2] = c.GetMecUtilizationAggregated(model.ExpOptimal, model.StrHybrid, MecCentral, resource)

	values[3][0] = c.GetMecUtilizationAggregated(model.ExpEarHeuristic, model.StrHybrid, MecLocal, resource)
	values[3][1] = c.GetMecUtilizationAggregated(model.ExpEarHeuristic, model.StrHybrid, MecRegional, resource)
	values[3][2] = c.GetMecUtilizationAggregated(model.ExpEarHeuristic, model.StrHybrid, MecCentral, resource)

	values[4][0] = c.GetMecUtilizationAggregated(model.ExpEarHeuristic, model.StrLB, MecLocal, resource)
	values[4][1] = c.GetMecUtilizationAggregated(model.ExpEarHeuristic, model.StrLB, MecRegional, resource)
	values[4][2] = c.GetMecUtilizationAggregated(model.ExpEarHeuristic, model.StrLB, MecCentral, resource)

	log.Infof("values length: %v, size: %v", len(values), len(values[0]))

	values[5][0] = c.GetMecUtilizationAggregated(model.ExpHeuristic, model.StrHybrid, MecLocal, resource)
	values[5][1] = c.GetMecUtilizationAggregated(model.ExpHeuristic, model.StrHybrid, MecRegional, resource)
	values[5][2] = c.GetMecUtilizationAggregated(model.ExpHeuristic, model.StrHybrid, MecCentral, resource)

	err := c.genRatesPkgAggregatedMecs(resource, values, basePath, false)
	if err != nil {
		log.Errorf("Error: %v", err)
	}

	return nil
}

func (c *Client) GenerateChartPkgMecsICCTunning(chartType ChartType, basePath string) error {
	var resource string

	switch chartType {
	case ResCpu:
		resource = "cpu"
	case ResMemory:
		resource = "memory"
	}

	var values = [][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}}

	values[0][0] = c.GetMecUtilizationAggregated(model.ExpEarHeuristic, model.StrHybrid, MecLocal, resource)
	values[0][1] = c.GetMecUtilizationAggregated(model.ExpEarHeuristic, model.StrHybrid, MecRegional, resource)
	values[0][2] = c.GetMecUtilizationAggregated(model.ExpEarHeuristic, model.StrHybrid, MecCentral, resource)

	values[1][0] = c.GetMecUtilizationAggregated(model.ExpEarHeuristic, model.StrLB, MecLocal, resource)
	values[1][1] = c.GetMecUtilizationAggregated(model.ExpEarHeuristic, model.StrLB, MecRegional, resource)
	values[1][2] = c.GetMecUtilizationAggregated(model.ExpEarHeuristic, model.StrLB, MecCentral, resource)

	values[2][0] = c.GetMecUtilizationAggregated(model.ExpEarHeuristic, model.StrLatency, MecLocal, resource)
	values[2][1] = c.GetMecUtilizationAggregated(model.ExpEarHeuristic, model.StrLatency, MecRegional, resource)
	values[2][2] = c.GetMecUtilizationAggregated(model.ExpEarHeuristic, model.StrLatency, MecCentral, resource)

	values[3][0] = c.GetMecUtilizationAggregated(model.ExpEarHeuristic, model.Str7L3R, MecLocal, resource)
	values[3][1] = c.GetMecUtilizationAggregated(model.ExpEarHeuristic, model.Str7L3R, MecRegional, resource)
	values[3][2] = c.GetMecUtilizationAggregated(model.ExpEarHeuristic, model.Str7L3R, MecCentral, resource)

	values[4][0] = c.GetMecUtilizationAggregated(model.ExpEarHeuristic, model.Str3L7R, MecLocal, resource)
	values[4][1] = c.GetMecUtilizationAggregated(model.ExpEarHeuristic, model.Str3L7R, MecRegional, resource)
	values[4][2] = c.GetMecUtilizationAggregated(model.ExpEarHeuristic, model.Str3L7R, MecCentral, resource)

	err := c.genRatesPkgAggregatedMecs(resource, values, basePath, true)
	if err != nil {
		log.Errorf("Error: %v", err)
	}

	return nil
}

func (c *Client) GenerateChartPkgApps(chartType ChartType, basePath string) error {

	switch chartType {
	case RelocationTriggeringRates:
		values := initializeEmpty2DArray()

		values[0][0] = c.GetRateValueAllIter(model.ExpOptimal, model.StrHybrid, "triggering", model.CG)
		values[0][1] = c.GetRateValueAllIter(model.ExpOptimal, model.StrHybrid, "triggering", model.V2X)
		values[0][2] = c.GetRateValueAllIter(model.ExpOptimal, model.StrHybrid, "triggering", model.UAV)
		values[0][3] = c.GetConfidenceValue(model.ExpOptimal, model.StrHybrid, "triggering")

		values[1][0] = c.GetRateValueAllIter(model.ExpHeuristic, model.StrHybrid, "triggering", model.CG)
		values[1][1] = c.GetRateValueAllIter(model.ExpHeuristic, model.StrHybrid, "triggering", model.V2X)
		values[1][2] = c.GetRateValueAllIter(model.ExpHeuristic, model.StrHybrid, "triggering", model.UAV)
		values[1][3] = c.GetConfidenceValue(model.ExpHeuristic, model.StrHybrid, "triggering")

		values[2][0] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrHybrid, "triggering", model.CG)
		values[2][1] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrHybrid, "triggering", model.V2X)
		values[2][2] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrHybrid, "triggering", model.UAV)
		values[2][3] = c.GetConfidenceValue(model.ExpEarHeuristic, model.StrHybrid, "triggering")

		values[3][0] = c.GetRateValueAllIter(model.ExpMLMasked, model.StrML, "triggering", model.CG)
		values[3][1] = c.GetRateValueAllIter(model.ExpMLMasked, model.StrML, "triggering", model.V2X)
		values[3][2] = c.GetRateValueAllIter(model.ExpMLMasked, model.StrML, "triggering", model.UAV)
		values[3][3] = c.GetConfidenceValue(model.ExpMLMasked, model.StrML, "triggering")

		values[4][0] = c.GetRateValueAllIter(model.ExpMLNonMasked, model.StrML, "triggering", model.CG)
		values[4][1] = c.GetRateValueAllIter(model.ExpMLNonMasked, model.StrML, "triggering", model.V2X)
		values[4][2] = c.GetRateValueAllIter(model.ExpMLNonMasked, model.StrML, "triggering", model.UAV)
		values[4][3] = c.GetConfidenceValue(model.ExpMLNonMasked, model.StrML, "triggering")

		err := c.genRatesPkgAggregatedApps("triggering", values, basePath)
		return err
	case RelocationRejectionRates:
		values := initializeEmpty2DArray()

		values[0][0] = c.GetRateValueAllIter(model.ExpOptimal, model.StrHybrid, "failed", model.CG)
		values[0][1] = c.GetRateValueAllIter(model.ExpOptimal, model.StrHybrid, "failed", model.V2X)
		values[0][2] = c.GetRateValueAllIter(model.ExpOptimal, model.StrHybrid, "failed", model.UAV)
		values[0][3] = c.GetConfidenceValue(model.ExpOptimal, model.StrHybrid, "failed")

		values[1][0] = c.GetRateValueAllIter(model.ExpHeuristic, model.StrHybrid, "failed", model.CG)
		values[1][1] = c.GetRateValueAllIter(model.ExpHeuristic, model.StrHybrid, "failed", model.V2X)
		values[1][2] = c.GetRateValueAllIter(model.ExpHeuristic, model.StrHybrid, "failed", model.UAV)
		values[1][3] = c.GetConfidenceValue(model.ExpHeuristic, model.StrHybrid, "failed")

		values[2][0] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrHybrid, "failed", model.CG)
		values[2][1] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrHybrid, "failed", model.V2X)
		values[2][2] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrHybrid, "failed", model.UAV)
		values[2][3] = c.GetConfidenceValue(model.ExpEarHeuristic, model.StrHybrid, "failed")

		values[3][0] = c.GetRateValueAllIter(model.ExpMLMasked, model.StrML, "failed", model.CG)
		values[3][1] = c.GetRateValueAllIter(model.ExpMLMasked, model.StrML, "failed", model.V2X)
		values[3][2] = c.GetRateValueAllIter(model.ExpMLMasked, model.StrML, "failed", model.UAV)
		values[3][3] = c.GetConfidenceValue(model.ExpMLMasked, model.StrML, "failed")

		values[4][0] = c.GetRateValueAllIter(model.ExpMLNonMasked, model.StrML, "failed", model.CG)
		values[4][1] = c.GetRateValueAllIter(model.ExpMLNonMasked, model.StrML, "failed", model.V2X)
		values[4][2] = c.GetRateValueAllIter(model.ExpMLNonMasked, model.StrML, "failed", model.UAV)
		values[4][3] = c.GetConfidenceValue(model.ExpMLNonMasked, model.StrML, "failed")

		err := c.genRatesPkgAggregatedApps("rejection", values, basePath)
		return err
	case RelocationSuccessfulSearchRates:
		values := initializeEmpty2DArray()

		values[0][0] = c.GetRateValueAllIter(model.ExpOptimal, model.StrHybrid, "successful", model.CG)
		values[0][1] = c.GetRateValueAllIter(model.ExpOptimal, model.StrHybrid, "successful", model.V2X)
		values[0][2] = c.GetRateValueAllIter(model.ExpOptimal, model.StrHybrid, "successful", model.UAV)
		values[0][3] = c.GetConfidenceValue(model.ExpOptimal, model.StrHybrid, "successful")

		values[1][0] = c.GetRateValueAllIter(model.ExpHeuristic, model.StrHybrid, "successful", model.CG)
		values[1][1] = c.GetRateValueAllIter(model.ExpHeuristic, model.StrHybrid, "successful", model.V2X)
		values[1][2] = c.GetRateValueAllIter(model.ExpHeuristic, model.StrHybrid, "successful", model.UAV)
		values[1][3] = c.GetConfidenceValue(model.ExpHeuristic, model.StrHybrid, "successful")

		values[2][0] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrHybrid, "successful", model.CG)
		values[2][1] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrHybrid, "successful", model.V2X)
		values[2][2] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrHybrid, "successful", model.UAV)
		values[2][3] = c.GetConfidenceValue(model.ExpEarHeuristic, model.StrHybrid, "successful")

		values[3][0] = c.GetRateValueAllIter(model.ExpMLMasked, model.StrML, "successful", model.CG)
		values[3][1] = c.GetRateValueAllIter(model.ExpMLMasked, model.StrML, "successful", model.V2X)
		values[3][2] = c.GetRateValueAllIter(model.ExpMLMasked, model.StrML, "successful", model.UAV)
		values[3][3] = c.GetConfidenceValue(model.ExpMLMasked, model.StrML, "successful")

		values[4][0] = c.GetRateValueAllIter(model.ExpMLNonMasked, model.StrML, "successful", model.CG)
		values[4][1] = c.GetRateValueAllIter(model.ExpMLNonMasked, model.StrML, "successful", model.V2X)
		values[4][2] = c.GetRateValueAllIter(model.ExpMLNonMasked, model.StrML, "successful", model.UAV)
		values[4][3] = c.GetConfidenceValue(model.ExpMLNonMasked, model.StrML, "successful")

		err := c.genRatesPkgAggregatedApps("successful-search", values, basePath)
		return err
	default:
		return errors.New("chartType not found")
	}
}

func (c *Client) GenerateChartPkgAppsICC(chartType ChartType, basePath string) error {

	switch chartType {
	case RelocationTriggeringRates:
		values := initializeEmpty2DArray()

		values[0][0] = c.GetRateValueAllIter(model.ExpOptimal, model.StrLB, "triggering", model.CG)
		values[0][1] = c.GetRateValueAllIter(model.ExpOptimal, model.StrLB, "triggering", model.V2X)
		values[0][2] = c.GetRateValueAllIter(model.ExpOptimal, model.StrLB, "triggering", model.UAV)
		values[0][3] = c.GetConfidenceValue(model.ExpOptimal, model.StrLB, "triggering")

		values[1][0] = c.GetRateValueAllIter(model.ExpOptimal, model.StrLatency, "triggering", model.CG)
		values[1][1] = c.GetRateValueAllIter(model.ExpOptimal, model.StrLatency, "triggering", model.V2X)
		values[1][2] = c.GetRateValueAllIter(model.ExpOptimal, model.StrLatency, "triggering", model.UAV)
		values[1][3] = c.GetConfidenceValue(model.ExpOptimal, model.StrLatency, "triggering")

		values[2][0] = c.GetRateValueAllIter(model.ExpOptimal, model.StrHybrid, "triggering", model.CG)
		values[2][1] = c.GetRateValueAllIter(model.ExpOptimal, model.StrHybrid, "triggering", model.V2X)
		values[2][2] = c.GetRateValueAllIter(model.ExpOptimal, model.StrHybrid, "triggering", model.UAV)
		values[2][3] = c.GetConfidenceValue(model.ExpOptimal, model.StrHybrid, "triggering")

		values[3][0] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrHybrid, "triggering", model.CG)
		values[3][1] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrHybrid, "triggering", model.V2X)
		values[3][2] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrHybrid, "triggering", model.UAV)
		values[3][3] = c.GetConfidenceValue(model.ExpEarHeuristic, model.StrHybrid, "triggering")

		values[4][0] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrLB, "triggering", model.CG)
		values[4][1] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrLB, "triggering", model.V2X)
		values[4][2] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrLB, "triggering", model.UAV)
		values[4][3] = c.GetConfidenceValue(model.ExpEarHeuristic, model.StrLB, "triggering")

		values[5][0] = c.GetRateValueAllIter(model.ExpHeuristic, model.StrHybrid, "triggering", model.CG)
		values[5][1] = c.GetRateValueAllIter(model.ExpHeuristic, model.StrHybrid, "triggering", model.V2X)
		values[5][2] = c.GetRateValueAllIter(model.ExpHeuristic, model.StrHybrid, "triggering", model.UAV)
		values[5][3] = c.GetConfidenceValue(model.ExpHeuristic, model.StrHybrid, "triggering")

		err := c.genRatesPkgAggregatedAppsICC("triggering", values, basePath, false)
		return err
	case RelocationRejectionRates:
		values := initializeEmpty2DArray()

		values[0][0] = c.GetRateValueAllIter(model.ExpOptimal, model.StrLB, "failed", model.CG)
		values[0][1] = c.GetRateValueAllIter(model.ExpOptimal, model.StrLB, "failed", model.V2X)
		values[0][2] = c.GetRateValueAllIter(model.ExpOptimal, model.StrLB, "failed", model.UAV)
		values[0][3] = c.GetConfidenceValue(model.ExpOptimal, model.StrLB, "failed")

		values[1][0] = c.GetRateValueAllIter(model.ExpOptimal, model.StrLatency, "failed", model.CG)
		values[1][1] = c.GetRateValueAllIter(model.ExpOptimal, model.StrLatency, "failed", model.V2X)
		values[1][2] = c.GetRateValueAllIter(model.ExpOptimal, model.StrLatency, "failed", model.UAV)
		values[1][3] = c.GetConfidenceValue(model.ExpOptimal, model.StrLatency, "failed")

		values[2][0] = c.GetRateValueAllIter(model.ExpOptimal, model.StrHybrid, "failed", model.CG)
		values[2][1] = c.GetRateValueAllIter(model.ExpOptimal, model.StrHybrid, "failed", model.V2X)
		values[2][2] = c.GetRateValueAllIter(model.ExpOptimal, model.StrHybrid, "failed", model.UAV)
		values[2][3] = c.GetConfidenceValue(model.ExpOptimal, model.StrHybrid, "failed")

		values[3][0] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrHybrid, "failed", model.CG)
		values[3][1] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrHybrid, "failed", model.V2X)
		values[3][2] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrHybrid, "failed", model.UAV)
		values[3][3] = c.GetConfidenceValue(model.ExpEarHeuristic, model.StrHybrid, "failed")

		values[4][0] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrLB, "failed", model.CG)
		values[4][1] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrLB, "failed", model.V2X)
		values[4][2] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrLB, "failed", model.UAV)
		values[4][3] = c.GetConfidenceValue(model.ExpEarHeuristic, model.StrLB, "failed")

		values[5][0] = c.GetRateValueAllIter(model.ExpHeuristic, model.StrHybrid, "failed", model.CG)
		values[5][1] = c.GetRateValueAllIter(model.ExpHeuristic, model.StrHybrid, "failed", model.V2X)
		values[5][2] = c.GetRateValueAllIter(model.ExpHeuristic, model.StrHybrid, "failed", model.UAV)
		values[5][3] = c.GetConfidenceValue(model.ExpHeuristic, model.StrHybrid, "failed")
		err := c.genRatesPkgAggregatedAppsICC("rejection", values, basePath, false)
		return err

	default:
		return errors.New("chartType not found")
	}
}

func (c *Client) GenerateChartPkgAppsICCTunning(chartType ChartType, basePath string) error {

	switch chartType {
	case RelocationTriggeringRates:
		values := initializeEmpty2DArray()

		values[0][0] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrHybrid, "triggering", model.CG)
		values[0][1] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrHybrid, "triggering", model.V2X)
		values[0][2] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrHybrid, "triggering", model.UAV)
		values[0][3] = c.GetConfidenceValue(model.ExpEarHeuristic, model.StrHybrid, "triggering")

		values[1][0] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrLB, "triggering", model.CG)
		values[1][1] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrLB, "triggering", model.V2X)
		values[1][2] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrLB, "triggering", model.UAV)
		values[1][3] = c.GetConfidenceValue(model.ExpEarHeuristic, model.StrLB, "triggering")

		values[2][0] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrLatency, "triggering", model.CG)
		values[2][1] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrLatency, "triggering", model.V2X)
		values[2][2] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrLatency, "triggering", model.UAV)
		values[2][3] = c.GetConfidenceValue(model.ExpEarHeuristic, model.StrLatency, "triggering")

		values[3][0] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.Str7L3R, "triggering", model.CG)
		values[3][1] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.Str7L3R, "triggering", model.V2X)
		values[3][2] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.Str7L3R, "triggering", model.UAV)
		values[3][3] = c.GetConfidenceValue(model.ExpEarHeuristic, model.Str7L3R, "triggering")

		values[4][0] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.Str3L7R, "triggering", model.CG)
		values[4][1] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.Str3L7R, "triggering", model.V2X)
		values[4][2] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.Str3L7R, "triggering", model.UAV)
		values[4][3] = c.GetConfidenceValue(model.ExpEarHeuristic, model.Str3L7R, "triggering")

		err := c.genRatesPkgAggregatedAppsICC("triggering", values, basePath, true)
		return err
	case RelocationRejectionRates:
		values := initializeEmpty2DArray()

		values[0][0] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrHybrid, "failed", model.CG)
		values[0][1] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrHybrid, "failed", model.V2X)
		values[0][2] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrHybrid, "failed", model.UAV)
		values[0][3] = c.GetConfidenceValue(model.ExpEarHeuristic, model.StrHybrid, "failed")

		values[1][0] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrLB, "failed", model.CG)
		values[1][1] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrLB, "failed", model.V2X)
		values[1][2] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrLB, "failed", model.UAV)
		values[1][3] = c.GetConfidenceValue(model.ExpEarHeuristic, model.StrLB, "failed")

		values[2][0] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrLatency, "failed", model.CG)
		values[2][1] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrLatency, "failed", model.V2X)
		values[2][2] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.StrLatency, "failed", model.UAV)
		values[2][3] = c.GetConfidenceValue(model.ExpEarHeuristic, model.StrLatency, "failed")

		values[3][0] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.Str7L3R, "failed", model.CG)
		values[3][1] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.Str7L3R, "failed", model.V2X)
		values[3][2] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.Str7L3R, "failed", model.UAV)
		values[3][3] = c.GetConfidenceValue(model.ExpEarHeuristic, model.Str7L3R, "failed")

		values[4][0] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.Str3L7R, "failed", model.CG)
		values[4][1] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.Str3L7R, "failed", model.V2X)
		values[4][2] = c.GetRateValueAllIter(model.ExpEarHeuristic, model.Str3L7R, "failed", model.UAV)
		values[4][3] = c.GetConfidenceValue(model.ExpEarHeuristic, model.Str3L7R, "failed")

		err := c.genRatesPkgAggregatedAppsICC("rejection", values, basePath, true)
		return err

	default:
		return errors.New("chartType not found")
	}
}

func (c *Client) genRatesPkgAggregatedMecs(resType string, values [][]float64, basePath string, tuning bool) error {
	mecTypeLabels := []string{"City-Level", "Regional-Level", "International-Level"}

	iterFiles := []string{"algorithm1.dat", "algorithm2.dat", "algorithm3.dat", "algorithm4.dat", "algorithm5.dat", "algorithm6.dat"}

	pkgPath := resType + "-aggregated-rates"
	scriptName := "_" + resType + "-aggregated.sh"

	title := "<TODO xLabel>"
	yLabel := fmt.Sprintf("Edge Server %s Utilization [%c]", strings.ToTitle(resType), '%')

	err := os.MkdirAll(basePath+"/"+pkgPath, os.ModePerm)
	if err != nil {
		return err
	}

	for i := range iterFiles {
		iter, err := os.Create(filepath.Join(basePath+"/"+pkgPath, filepath.Base(iterFiles[i])))
		defer iter.Close()
		if err != nil {
			return err
		}

		_, err = iter.WriteString(createIterFileContentMecs(mecTypeLabels, values[i]))
		if err != nil {
			return err
		}
	}

	script, err := os.Create(filepath.Join(basePath+"/"+pkgPath, filepath.Base(scriptName)))
	defer script.Close()
	if err != nil {
		return err
	}

	if tuning == true {
		_, err = script.WriteString(generateAggregatedRatesScriptMecsTuning(resType, title, yLabel, iterFiles))
		if err != nil {
			return err
		}
		return err
	} else {
		_, err = script.WriteString(generateAggregatedRatesScriptMecs(resType, title, yLabel, iterFiles))
		if err != nil {
			return err
		}
		return err
	}

}

func (c *Client) genRatesPkgAggregatedApps(ratesType string, values [][]float64, basePath string) error {
	expLabels := []string{"Optimal", "Heuristic", "EAR-Heuristic", "ML-Masked", "ML-NonMasked"}

	iterFile := ratesType + ".dat"

	pkgPath := ratesType + "-aggregated-rates"
	scriptName := "_" + ratesType + "-aggregated.sh"

	xLabel := "Time"
	yLabel := fmt.Sprintf("%s Rate [%c]", strings.ToTitle(ratesType), '%')

	err := os.MkdirAll(basePath+"/"+pkgPath, os.ModePerm)
	if err != nil {
		return err
	}

	iter, err := os.Create(filepath.Join(basePath+"/"+pkgPath, filepath.Base(iterFile)))
	defer iter.Close()
	if err != nil {
		return err
	}

	_, err = iter.WriteString(createIterFileContentApps(0, expLabels, values))
	if err != nil {
		return err
	}

	script, err := os.Create(filepath.Join(basePath+"/"+pkgPath, filepath.Base(scriptName)))
	defer script.Close()
	if err != nil {
		return err
	}

	_, err = script.WriteString(generateAggregatedRatesScriptApps(ratesType, xLabel, yLabel, iterFile))
	if err != nil {
		return err
	}
	return err
}

func (c *Client) genRatesPkgAggregatedAppsICC(ratesType string, values [][]float64, basePath string, tunning bool) error {

	var expLabels []string
	if tunning == true {
		expLabels = []string{"L5R5", "L1R0", "L0R1", "L7R3", "L3R7"}
	} else {
		expLabels = []string{"O-LoadBalancing", "O-Latency", "O-Hybrid", "EAR-Hybrid", "EAR-LB", "H-Hybrid"}
	}

	iterFile := ratesType + ".dat"

	pkgPath := ratesType + "-aggregated-rates"
	scriptName := "_" + ratesType + "-aggregated.sh"

	xLabel := "Time"
	yLabel := fmt.Sprintf("%s Rate [%c]", strings.ToTitle(ratesType), '%')

	err := os.MkdirAll(basePath+"/"+pkgPath, os.ModePerm)
	if err != nil {
		return err
	}

	iter, err := os.Create(filepath.Join(basePath+"/"+pkgPath, filepath.Base(iterFile)))
	defer iter.Close()
	if err != nil {
		return err
	}

	_, err = iter.WriteString(createIterFileContentApps(0, expLabels, values))
	if err != nil {
		return err
	}

	script, err := os.Create(filepath.Join(basePath+"/"+pkgPath, filepath.Base(scriptName)))
	defer script.Close()
	if err != nil {
		return err
	}

	_, err = script.WriteString(generateAggregatedRatesScriptApps(ratesType, xLabel, yLabel, iterFile))
	if err != nil {
		return err
	}
	return err
}

func (c *Client) genRatesPkgManyIterations(ratesType string, val5x4 map[int][][]float64, basePath string) error {
	expLabels := []string{"O-LoadBalancing", "O-Latency", "O-Hybrid", "EAR-Heuristic", "H-Hybrid"}
	iterFiles := []string{"iter0.dat", "iter1.dat", "iter2.dat", "iter3.dat", "iter4.dat"}

	pkgPath := ratesType + "-rates"
	scriptName := "_" + ratesType + ".sh"

	xLabel := "Time"
	yLabel := fmt.Sprintf("%s Rate [%c]", strings.ToTitle(ratesType), '%')

	err := os.MkdirAll(basePath+"/"+pkgPath, os.ModePerm)
	if err != nil {
		return err
	}

	for i := 0; i < 5; i++ {
		iter, err := os.Create(filepath.Join(basePath+"/"+pkgPath, filepath.Base(iterFiles[i])))
		defer iter.Close()
		if err != nil {
			return err
		}

		_, err = iter.WriteString(createIterFileContentApps(i, expLabels, val5x4[i]))
		if err != nil {
			return err
		}
	}

	script, err := os.Create(filepath.Join(basePath+"/"+pkgPath, filepath.Base(scriptName)))
	defer script.Close()
	if err != nil {
		return err
	}

	_, err = script.WriteString(generateRatesScript(ratesType, xLabel, yLabel, iterFiles))
	if err != nil {
		return err
	}

	return err
}

func (c *Client) genResCpuPkg(basePath string) error {
	return nil
}

func (c *Client) genResMemoryPkg(basePath string) error {
	return nil
}

func createIterFileContentMecs(labels []string, val []float64) string {
	fileContent := fmt.Sprintf("#\tmec-type\tutil-rate\n")

	for i := 0; i < len(labels); i++ {
		line := fmt.Sprintf("%s\t%.2f\n", labels[i], val[i])
		fileContent += line
	}

	return fileContent
}

func createIterFileContentApps(iterFileNo int, labels []string, val [][]float64) string {
	fileContent := fmt.Sprintf("#\tcg\tv2x\tuav\tconf\tx\n")

	var startIndex int
	switch iterFileNo {
	case 0:
		startIndex = 0
	case 1:
		startIndex = 6
	case 2:
		startIndex = 12
	case 3:
		startIndex = 18
	case 4:
		startIndex = 24
	}

	for i := 0; i < len(labels); i++ {
		line := fmt.Sprintf("%s\t%.2f\t%.2f\t%.2f\t%.2f\t%d\n", labels[i], val[i][0], val[i][1], val[i][2], val[i][3], startIndex+i)
		fileContent += line
	}

	return fileContent
}
