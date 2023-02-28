package results

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	log "simu/src/logger"
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

	// 0 - lb, 1 - lat, 2 - hyb, 3 - ear, 4 - heu | 0 - local, 1 - regional, 2 - central
	values[0][0] = c.GetMecUtilizationAggregated(ExpOptimal, "lb", MecLocal, resource)
	values[0][1] = c.GetMecUtilizationAggregated(ExpOptimal, "lb", MecRegional, resource)
	values[0][2] = c.GetMecUtilizationAggregated(ExpOptimal, "lb", MecCentral, resource)

	values[1][0] = c.GetMecUtilizationAggregated(ExpOptimal, "latency", MecLocal, resource)
	values[1][1] = c.GetMecUtilizationAggregated(ExpOptimal, "latency", MecRegional, resource)
	values[1][2] = c.GetMecUtilizationAggregated(ExpOptimal, "latency", MecCentral, resource)

	values[2][0] = c.GetMecUtilizationAggregated(ExpOptimal, "hybrid", MecLocal, resource)
	values[2][1] = c.GetMecUtilizationAggregated(ExpOptimal, "hybrid", MecRegional, resource)
	values[2][2] = c.GetMecUtilizationAggregated(ExpOptimal, "hybrid", MecCentral, resource)

	values[3][0] = c.GetMecUtilizationAggregated(ExpEarHeuristic, "hybrid", MecLocal, resource)
	values[3][1] = c.GetMecUtilizationAggregated(ExpEarHeuristic, "hybrid", MecRegional, resource)
	values[3][2] = c.GetMecUtilizationAggregated(ExpEarHeuristic, "hybrid", MecCentral, resource)

	values[4][0] = c.GetMecUtilizationAggregated(ExpHeuristic, "hybrid", MecLocal, resource)
	values[4][1] = c.GetMecUtilizationAggregated(ExpHeuristic, "hybrid", MecRegional, resource)
	values[4][2] = c.GetMecUtilizationAggregated(ExpHeuristic, "hybrid", MecCentral, resource)

	err := c.genRatesPkgAggregatedMecs(resource, values, basePath)
	if err != nil {
		log.Errorf("Error: %v", err)
	}

	return nil
}

func (c *Client) GenerateChartPkgApps(chartType ChartType, basePath string) error {

	lbIterator, latencyIterator, hybridIterator, earIterator, hybIterator := 0, 0, 0, 0, 0
	valuesManyIterations := map[int][][]float64{
		0: initializeEmpty2DArray(),
		1: initializeEmpty2DArray(),
		2: initializeEmpty2DArray(),
		3: initializeEmpty2DArray(),
		4: initializeEmpty2DArray(),
	}

	switch chartType {
	case RelocationRates:
		for _, result := range c.GetResults() {
			switch result.Metadata.Type {
			case ExpOptimal:
				switch result.Metadata.Strategy {
				case "lb":
					if iteratorConstraint(lbIterator) {
						continue
					}
					valuesManyIterations[lbIterator][0][0] = float64(result.Data.Erd.Successful[CG]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[lbIterator][0][1] = float64(result.Data.Erd.Successful[V2X]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[lbIterator][0][2] = float64(result.Data.Erd.Successful[UAV]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[lbIterator][0][3] = 0 // TODO: confidence
					lbIterator++
				case "latency":
					if iteratorConstraint(latencyIterator) {
						continue
					}
					valuesManyIterations[latencyIterator][1][0] = float64(result.Data.Erd.Successful[CG]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[latencyIterator][1][1] = float64(result.Data.Erd.Successful[V2X]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[latencyIterator][1][2] = float64(result.Data.Erd.Successful[UAV]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[latencyIterator][1][3] = 0 // TODO: confidence
					latencyIterator++
				case "hybrid":
					if iteratorConstraint(hybridIterator) {
						continue
					}
					valuesManyIterations[hybridIterator][2][0] = float64(result.Data.Erd.Successful[CG]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[hybridIterator][2][1] = float64(result.Data.Erd.Successful[V2X]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[hybridIterator][2][2] = float64(result.Data.Erd.Successful[UAV]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[hybridIterator][2][3] = 0 // TODO: confidence
					hybridIterator++
				}
			case ExpHeuristic:
				if iteratorConstraint(earIterator) {
					continue
				}
				valuesManyIterations[earIterator][3][0] = float64(result.Data.Erd.Successful[CG]) / float64(result.Metadata.Movements) * 100
				valuesManyIterations[earIterator][3][1] = float64(result.Data.Erd.Successful[V2X]) / float64(result.Metadata.Movements) * 100
				valuesManyIterations[earIterator][3][2] = float64(result.Data.Erd.Successful[UAV]) / float64(result.Metadata.Movements) * 100
				valuesManyIterations[earIterator][3][3] = 0 // TODO: confidence
				earIterator++

			case ExpEarHeuristic:
				if iteratorConstraint(hybIterator) {
					continue
				}
				valuesManyIterations[hybIterator][4][0] = float64(result.Data.Erd.Successful[CG]) / float64(result.Metadata.Movements) * 100
				valuesManyIterations[hybIterator][4][1] = float64(result.Data.Erd.Successful[V2X]) / float64(result.Metadata.Movements) * 100
				valuesManyIterations[hybIterator][4][2] = float64(result.Data.Erd.Successful[UAV]) / float64(result.Metadata.Movements) * 100
				valuesManyIterations[hybIterator][4][3] = 0 // TODO: confidence
				hybIterator++
			}
		}

		err := c.genRatesPkgManyIterations("relocations", valuesManyIterations, basePath)
		return err
	case SkippedRates:
		for _, result := range c.GetResults() {
			switch result.Metadata.Type {
			case ExpOptimal:
				switch result.Metadata.Strategy {
				case "lb":
					if iteratorConstraint(lbIterator) {
						continue
					}
					valuesManyIterations[lbIterator][0][0] = float64(result.Data.Erd.Skipped[CG]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[lbIterator][0][1] = float64(result.Data.Erd.Skipped[V2X]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[lbIterator][0][2] = float64(result.Data.Erd.Skipped[UAV]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[lbIterator][0][3] = 0 // TODO: confidence
					lbIterator++
				case "latency":
					if iteratorConstraint(latencyIterator) {
						continue
					}
					valuesManyIterations[latencyIterator][1][0] = float64(result.Data.Erd.Skipped[CG]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[latencyIterator][1][1] = float64(result.Data.Erd.Skipped[V2X]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[latencyIterator][1][2] = float64(result.Data.Erd.Skipped[UAV]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[latencyIterator][1][3] = 0 // TODO: confidence
					latencyIterator++
				case "hybrid":
					if iteratorConstraint(hybridIterator) {
						continue
					}
					valuesManyIterations[hybridIterator][2][0] = float64(result.Data.Erd.Skipped[CG]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[hybridIterator][2][1] = float64(result.Data.Erd.Skipped[V2X]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[hybridIterator][2][2] = float64(result.Data.Erd.Skipped[UAV]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[hybridIterator][2][3] = 0 // TODO: confidence
					hybridIterator++
				}
			case ExpHeuristic:
				if iteratorConstraint(earIterator) {
					continue
				}
				valuesManyIterations[earIterator][3][0] = float64(result.Data.Erd.Skipped[CG]) / float64(result.Metadata.Movements) * 100
				valuesManyIterations[earIterator][3][1] = float64(result.Data.Erd.Skipped[V2X]) / float64(result.Metadata.Movements) * 100
				valuesManyIterations[earIterator][3][2] = float64(result.Data.Erd.Skipped[UAV]) / float64(result.Metadata.Movements) * 100
				valuesManyIterations[earIterator][3][3] = 0 // TODO: confidence
				earIterator++

			case ExpEarHeuristic:
				if iteratorConstraint(hybIterator) {
					continue
				}
				valuesManyIterations[hybIterator][4][0] = float64(result.Data.Erd.Skipped[CG]) / float64(result.Metadata.Movements) * 100
				valuesManyIterations[hybIterator][4][1] = float64(result.Data.Erd.Skipped[V2X]) / float64(result.Metadata.Movements) * 100
				valuesManyIterations[hybIterator][4][2] = float64(result.Data.Erd.Skipped[UAV]) / float64(result.Metadata.Movements) * 100
				valuesManyIterations[hybIterator][4][3] = 0 // TODO: confidence
				hybIterator++
			}
		}

		err := c.genRatesPkgManyIterations("skipped", valuesManyIterations, basePath)
		return err
	case RedundantRates:
		for _, result := range c.GetResults() {
			switch result.Metadata.Type {
			case ExpOptimal:
				switch result.Metadata.Strategy {
				case "lb":
					if iteratorConstraint(lbIterator) {
						continue
					}
					valuesManyIterations[lbIterator][0][0] = float64(result.Data.Erd.Redundant[CG]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[lbIterator][0][1] = float64(result.Data.Erd.Redundant[V2X]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[lbIterator][0][2] = float64(result.Data.Erd.Redundant[UAV]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[lbIterator][0][3] = 0 // TODO: confidence
					lbIterator++
				case "latency":
					if iteratorConstraint(latencyIterator) {
						continue
					}
					valuesManyIterations[latencyIterator][1][0] = float64(result.Data.Erd.Redundant[CG]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[latencyIterator][1][1] = float64(result.Data.Erd.Redundant[V2X]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[latencyIterator][1][2] = float64(result.Data.Erd.Redundant[UAV]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[latencyIterator][1][3] = 0 // TODO: confidence
					latencyIterator++
				case "hybrid":
					if iteratorConstraint(hybridIterator) {
						continue
					}
					valuesManyIterations[hybridIterator][2][0] = float64(result.Data.Erd.Redundant[CG]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[hybridIterator][2][1] = float64(result.Data.Erd.Redundant[V2X]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[hybridIterator][2][2] = float64(result.Data.Erd.Redundant[UAV]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[hybridIterator][2][3] = 0 // TODO: confidence
					hybridIterator++
				}
			case ExpHeuristic:
				if iteratorConstraint(earIterator) {
					continue
				}
				valuesManyIterations[earIterator][3][0] = float64(result.Data.Erd.Redundant[CG]) / float64(result.Metadata.Movements) * 100
				valuesManyIterations[earIterator][3][1] = float64(result.Data.Erd.Redundant[V2X]) / float64(result.Metadata.Movements) * 100
				valuesManyIterations[earIterator][3][2] = float64(result.Data.Erd.Redundant[UAV]) / float64(result.Metadata.Movements) * 100
				valuesManyIterations[earIterator][3][3] = 0 // TODO: confidence
				earIterator++

			case ExpEarHeuristic:
				if iteratorConstraint(hybIterator) {
					continue
				}
				valuesManyIterations[hybIterator][4][0] = float64(result.Data.Erd.Redundant[CG]) / float64(result.Metadata.Movements) * 100
				valuesManyIterations[hybIterator][4][1] = float64(result.Data.Erd.Redundant[V2X]) / float64(result.Metadata.Movements) * 100
				valuesManyIterations[hybIterator][4][2] = float64(result.Data.Erd.Redundant[UAV]) / float64(result.Metadata.Movements) * 100
				valuesManyIterations[hybIterator][4][3] = 0 // TODO: confidence
				hybIterator++
			}
		}

		err := c.genRatesPkgManyIterations("redundant", valuesManyIterations, basePath)
		return err
	case FailedRates:
		for _, result := range c.GetResults() {
			switch result.Metadata.Type {
			case ExpOptimal:
				switch result.Metadata.Strategy {
				case "lb":
					if iteratorConstraint(lbIterator) {
						continue
					}
					valuesManyIterations[lbIterator][0][0] = float64(result.Data.Erd.Failed[CG]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[lbIterator][0][1] = float64(result.Data.Erd.Failed[V2X]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[lbIterator][0][2] = float64(result.Data.Erd.Failed[UAV]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[lbIterator][0][3] = 0 // TODO: confidence
					lbIterator++
				case "latency":
					if iteratorConstraint(latencyIterator) {
						continue
					}
					valuesManyIterations[latencyIterator][1][0] = float64(result.Data.Erd.Failed[CG]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[latencyIterator][1][1] = float64(result.Data.Erd.Failed[V2X]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[latencyIterator][1][2] = float64(result.Data.Erd.Failed[UAV]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[latencyIterator][1][3] = 0 // TODO: confidence
					latencyIterator++
				case "hybrid":
					if iteratorConstraint(hybridIterator) {
						continue
					}
					valuesManyIterations[hybridIterator][2][0] = float64(result.Data.Erd.Failed[CG]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[hybridIterator][2][1] = float64(result.Data.Erd.Failed[V2X]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[hybridIterator][2][2] = float64(result.Data.Erd.Failed[UAV]) / float64(result.Metadata.Movements) * 100
					valuesManyIterations[hybridIterator][2][3] = 0 // TODO: confidence
					hybridIterator++
				}
			case ExpHeuristic:
				if iteratorConstraint(earIterator) {
					continue
				}
				valuesManyIterations[earIterator][3][0] = float64(result.Data.Erd.Failed[CG]) / float64(result.Metadata.Movements) * 100
				valuesManyIterations[earIterator][3][1] = float64(result.Data.Erd.Failed[V2X]) / float64(result.Metadata.Movements) * 100
				valuesManyIterations[earIterator][3][2] = float64(result.Data.Erd.Failed[UAV]) / float64(result.Metadata.Movements) * 100
				valuesManyIterations[earIterator][3][3] = 0 // TODO: confidence
				earIterator++

			case ExpEarHeuristic:
				if iteratorConstraint(hybIterator) {
					continue
				}
				valuesManyIterations[hybIterator][4][0] = float64(result.Data.Erd.Failed[CG]) / float64(result.Metadata.Movements) * 100
				valuesManyIterations[hybIterator][4][1] = float64(result.Data.Erd.Failed[V2X]) / float64(result.Metadata.Movements) * 100
				valuesManyIterations[hybIterator][4][2] = float64(result.Data.Erd.Failed[UAV]) / float64(result.Metadata.Movements) * 100
				valuesManyIterations[hybIterator][4][3] = 0 // TODO: confidence
				hybIterator++
			}
		}

		err := c.genRatesPkgManyIterations("failed", valuesManyIterations, basePath)
		return err
	case RelocationTriggeringRates:
		values := initializeEmpty2DArray()
		values[0][0] = c.GetRateValue(ExpOptimal, "lb", "triggering", CG)
		values[0][1] = c.GetRateValue(ExpOptimal, "lb", "triggering", V2X)
		values[0][2] = c.GetRateValue(ExpOptimal, "lb", "triggering", UAV)
		values[0][3] = 0 // TODO: confidence

		values[1][0] = c.GetRateValue(ExpOptimal, "latency", "triggering", CG)
		values[1][1] = c.GetRateValue(ExpOptimal, "latency", "triggering", V2X)
		values[1][2] = c.GetRateValue(ExpOptimal, "latency", "triggering", UAV)
		values[1][3] = 0 // TODO: confidence

		values[2][0] = c.GetRateValue(ExpOptimal, "hybrid", "triggering", CG)
		values[2][1] = c.GetRateValue(ExpOptimal, "hybrid", "triggering", V2X)
		values[2][2] = c.GetRateValue(ExpOptimal, "hybrid", "triggering", UAV)
		values[2][3] = 0 // TODO: confidence

		values[3][0] = c.GetRateValue(ExpHeuristic, "hybrid", "triggering", CG)
		values[3][1] = c.GetRateValue(ExpHeuristic, "hybrid", "triggering", V2X)
		values[3][2] = c.GetRateValue(ExpHeuristic, "hybrid", "triggering", UAV)
		values[3][3] = 0 // TODO: confidence

		values[4][0] = c.GetRateValue(ExpEarHeuristic, "hybrid", "triggering", CG)
		values[4][1] = c.GetRateValue(ExpEarHeuristic, "hybrid", "triggering", CG)
		values[4][2] = c.GetRateValue(ExpEarHeuristic, "hybrid", "triggering", CG)
		values[4][3] = 0 // TODO: confidence

		err := c.genRatesPkgAggregatedApps("triggering", values, basePath)
		return err
	case RelocationRejectionRates:
		values := initializeEmpty2DArray()
		values[0][0] = c.GetRateValue(ExpOptimal, "lb", "failed", CG)
		values[0][1] = c.GetRateValue(ExpOptimal, "lb", "failed", V2X)
		values[0][2] = c.GetRateValue(ExpOptimal, "lb", "failed", UAV)
		values[0][3] = 0 // TODO: confidence

		values[1][0] = c.GetRateValue(ExpOptimal, "latency", "failed", CG)
		values[1][1] = c.GetRateValue(ExpOptimal, "latency", "failed", V2X)
		values[1][2] = c.GetRateValue(ExpOptimal, "latency", "failed", UAV)
		values[1][3] = 0 // TODO: confidence

		values[2][0] = c.GetRateValue(ExpOptimal, "hybrid", "failed", CG)
		values[2][1] = c.GetRateValue(ExpOptimal, "hybrid", "failed", V2X)
		values[2][2] = c.GetRateValue(ExpOptimal, "hybrid", "failed", UAV)
		values[2][3] = 0 // TODO: confidence

		values[3][0] = c.GetRateValue(ExpHeuristic, "hybrid", "failed", CG)
		values[3][1] = c.GetRateValue(ExpHeuristic, "hybrid", "failed", V2X)
		values[3][2] = c.GetRateValue(ExpHeuristic, "hybrid", "failed", UAV)
		values[3][3] = 0 // TODO: confidence

		values[4][0] = c.GetRateValue(ExpEarHeuristic, "hybrid", "failed", CG)
		values[4][1] = c.GetRateValue(ExpEarHeuristic, "hybrid", "failed", CG)
		values[4][2] = c.GetRateValue(ExpEarHeuristic, "hybrid", "failed", CG)
		values[4][3] = 0 // TODO: confidence

		err := c.genRatesPkgAggregatedApps("rejection", values, basePath)
		return err
	case RelocationSuccessfulSearchRates:
		values := initializeEmpty2DArray()
		values[0][0] = c.GetRateValue(ExpOptimal, "lb", "successful", CG)
		values[0][1] = c.GetRateValue(ExpOptimal, "lb", "successful", V2X)
		values[0][2] = c.GetRateValue(ExpOptimal, "lb", "successful", UAV)
		values[0][3] = 0 // TODO: confidence

		values[1][0] = c.GetRateValue(ExpOptimal, "latency", "successful", CG)
		values[1][1] = c.GetRateValue(ExpOptimal, "latency", "successful", V2X)
		values[1][2] = c.GetRateValue(ExpOptimal, "latency", "successful", UAV)
		values[1][3] = 0 // TODO: confidence

		values[2][0] = c.GetRateValue(ExpOptimal, "hybrid", "successful", CG)
		values[2][1] = c.GetRateValue(ExpOptimal, "hybrid", "successful", V2X)
		values[2][2] = c.GetRateValue(ExpOptimal, "hybrid", "successful", UAV)
		values[2][3] = 0 // TODO: confidence

		values[3][0] = c.GetRateValue(ExpHeuristic, "hybrid", "successful", CG)
		values[3][1] = c.GetRateValue(ExpHeuristic, "hybrid", "successful", V2X)
		values[3][2] = c.GetRateValue(ExpHeuristic, "hybrid", "successful", UAV)
		values[3][3] = 0 // TODO: confidence

		values[4][0] = c.GetRateValue(ExpEarHeuristic, "hybrid", "successful", CG)
		values[4][1] = c.GetRateValue(ExpEarHeuristic, "hybrid", "successful", CG)
		values[4][2] = c.GetRateValue(ExpEarHeuristic, "hybrid", "successful", CG)
		values[4][3] = 0 // TODO: confidence

		err := c.genRatesPkgAggregatedApps("successful-search", values, basePath)
		return err
	default:
		return errors.New("chartType not found")
	}
}

func (c *Client) genRatesPkgAggregatedMecs(resType string, values [][]float64, basePath string) error {
	mecTypeLabels := []string{"City-Level", "Regional-Level", "International-Level"}

	iterFiles := []string{"o-lb.dat", "o-lat.dat", "o-hyb.dat", "ear.dat", "h-hyb.dat"}

	pkgPath := resType + "-aggregated-rates"
	scriptName := "_" + resType + "-aggregated.sh"

	title := "<TODO xLabel>"
	yLabel := fmt.Sprintf("Edge Server %s Utilization [%c]", strings.ToTitle(resType), '%')

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

	_, err = script.WriteString(generateAggregatedRatesScriptMecs(resType, title, yLabel, iterFiles))
	if err != nil {
		return err
	}
	return err
}

func (c *Client) genRatesPkgAggregatedApps(ratesType string, values [][]float64, basePath string) error {
	expLabels := []string{"O-LoadBalancing", "O-Latency", "O-Hybrid", "EAR-Heuristic", "H-Hybrid"}

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

func createResFileContent(labels []string) string {
	var fileContent string
	return fileContent
}
