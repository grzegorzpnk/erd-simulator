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

	// TODO: Consider EAR-Heuristic as an alternative for Heuristic
	//values[1][0] = c.GetMecUtilizationAggregated(ExpEarHeuristic, StrHybrid, MecLocal, resource)
	//values[1][1] = c.GetMecUtilizationAggregated(ExpEarHeuristic, StrHybrid, MecRegional, resource)
	//values[1][2] = c.GetMecUtilizationAggregated(ExpEarHeuristic, StrHybrid, MecCentral, resource)

	values[2][0] = c.GetMecUtilizationAggregated(model.ExpMLMasked, model.StrML, MecLocal, resource)
	values[2][1] = c.GetMecUtilizationAggregated(model.ExpMLMasked, model.StrML, MecRegional, resource)
	values[2][2] = c.GetMecUtilizationAggregated(model.ExpMLMasked, model.StrML, MecCentral, resource)

	values[3][0] = c.GetMecUtilizationAggregated(model.ExpMLNonMasked, model.StrML, MecLocal, resource)
	values[3][1] = c.GetMecUtilizationAggregated(model.ExpMLNonMasked, model.StrML, MecRegional, resource)
	values[3][2] = c.GetMecUtilizationAggregated(model.ExpMLNonMasked, model.StrML, MecCentral, resource)

	err := c.genRatesPkgAggregatedMecs(resource, values, basePath)
	if err != nil {
		log.Errorf("Error: %v", err)
	}

	return nil
}

func (c *Client) GenerateChartPkgApps(chartType ChartType, basePath string) error {

	switch chartType {
	case RelocationTriggeringRates:
		values := initializeEmpty2DArray()

		values[0][0] = c.GetRateValue(model.ExpOptimal, model.StrHybrid, "triggering", model.CG)
		values[0][1] = c.GetRateValue(model.ExpOptimal, model.StrHybrid, "triggering", model.V2X)
		values[0][2] = c.GetRateValue(model.ExpOptimal, model.StrHybrid, "triggering", model.UAV)
		values[0][3] = 0 // TODO: confidence

		values[1][0] = c.GetRateValue(model.ExpHeuristic, model.StrHybrid, "triggering", model.CG)
		values[1][1] = c.GetRateValue(model.ExpHeuristic, model.StrHybrid, "triggering", model.V2X)
		values[1][2] = c.GetRateValue(model.ExpHeuristic, model.StrHybrid, "triggering", model.UAV)
		values[1][3] = 0 // TODO: confidence

		// TODO: As an alternative to Heuristic we can use EAR-Heuristic?
		//values[1][0] = c.GetRateValue(ExpEarHeuristic, StrHybrid, "triggering", CG)
		//values[1][1] = c.GetRateValue(ExpEarHeuristic, StrHybrid, "triggering", V2X)
		//values[1][2] = c.GetRateValue(ExpEarHeuristic, StrHybrid, "triggering", UAV)
		//values[1][3] = 0 // TODO: confidence

		values[2][0] = c.GetRateValue(model.ExpMLMasked, model.StrML, "triggering", model.CG)
		values[2][1] = c.GetRateValue(model.ExpMLMasked, model.StrML, "triggering", model.V2X)
		values[2][2] = c.GetRateValue(model.ExpMLMasked, model.StrML, "triggering", model.UAV)
		values[2][3] = 0 // TODO: confidence

		values[3][0] = c.GetRateValue(model.ExpMLNonMasked, model.StrML, "triggering", model.CG)
		values[3][1] = c.GetRateValue(model.ExpMLNonMasked, model.StrML, "triggering", model.V2X)
		values[3][2] = c.GetRateValue(model.ExpMLNonMasked, model.StrML, "triggering", model.UAV)
		values[3][3] = 0 // TODO: confidence

		err := c.genRatesPkgAggregatedApps("triggering", values, basePath)
		return err
	case RelocationRejectionRates:
		values := initializeEmpty2DArray()

		values[0][0] = c.GetRateValue(model.ExpOptimal, model.StrHybrid, "failed", model.CG)
		values[0][1] = c.GetRateValue(model.ExpOptimal, model.StrHybrid, "failed", model.V2X)
		values[0][2] = c.GetRateValue(model.ExpOptimal, model.StrHybrid, "failed", model.UAV)
		values[0][3] = 0 // TODO: confidence

		values[1][0] = c.GetRateValue(model.ExpHeuristic, model.StrHybrid, "failed", model.CG)
		values[1][1] = c.GetRateValue(model.ExpHeuristic, model.StrHybrid, "failed", model.V2X)
		values[1][2] = c.GetRateValue(model.ExpHeuristic, model.StrHybrid, "failed", model.UAV)
		values[1][3] = 0 // TODO: confidence

		// TODO: As an alternative to Heuristic we can use EAR-Heuristic?
		//values[1][0] = c.GetRateValue(ExpEarHeuristic, StrHybrid, "failed", CG)
		//values[1][1] = c.GetRateValue(ExpEarHeuristic, StrHybrid, "failed", V2X)
		//values[1][2] = c.GetRateValue(ExpEarHeuristic, StrHybrid, "failed", UAV)
		//values[1][3] = 0 // TODO: confidence

		values[2][0] = c.GetRateValue(model.ExpMLMasked, model.StrML, "failed", model.CG)
		values[2][1] = c.GetRateValue(model.ExpMLMasked, model.StrML, "failed", model.V2X)
		values[2][2] = c.GetRateValue(model.ExpMLMasked, model.StrML, "failed", model.UAV)
		values[2][3] = 0 // TODO: confidence

		values[3][0] = c.GetRateValue(model.ExpMLNonMasked, model.StrML, "failed", model.CG)
		values[3][1] = c.GetRateValue(model.ExpMLNonMasked, model.StrML, "failed", model.V2X)
		values[3][2] = c.GetRateValue(model.ExpMLNonMasked, model.StrML, "failed", model.UAV)
		values[3][3] = 0 // TODO: confidence

		err := c.genRatesPkgAggregatedApps("rejection", values, basePath)
		return err
	case RelocationSuccessfulSearchRates:
		values := initializeEmpty2DArray()

		values[0][0] = c.GetRateValue(model.ExpOptimal, model.StrHybrid, "successful", model.CG)
		values[0][1] = c.GetRateValue(model.ExpOptimal, model.StrHybrid, "successful", model.V2X)
		values[0][2] = c.GetRateValue(model.ExpOptimal, model.StrHybrid, "successful", model.UAV)
		values[0][3] = 0 // TODO: confidence

		values[1][0] = c.GetRateValue(model.ExpHeuristic, model.StrHybrid, "successful", model.CG)
		values[1][1] = c.GetRateValue(model.ExpHeuristic, model.StrHybrid, "successful", model.V2X)
		values[1][2] = c.GetRateValue(model.ExpHeuristic, model.StrHybrid, "successful", model.UAV)
		values[1][3] = 0 // TODO: confidence

		// TODO: As an alternative to Heuristic we can use EAR-Heuristic?
		//values[1][0] = c.GetRateValue(ExpEarHeuristic, StrHybrid, "successful", CG)
		//values[1][1] = c.GetRateValue(ExpEarHeuristic, StrHybrid, "successful", V2X)
		//values[1][2] = c.GetRateValue(ExpEarHeuristic, StrHybrid, "successful", UAV)
		//values[1][3] = 0 // TODO: confidence

		values[2][0] = c.GetRateValue(model.ExpMLMasked, model.StrML, "successful", model.CG)
		values[2][1] = c.GetRateValue(model.ExpMLMasked, model.StrML, "successful", model.V2X)
		values[2][2] = c.GetRateValue(model.ExpMLMasked, model.StrML, "successful", model.UAV)
		values[2][3] = 0 // TODO: confidence

		values[3][0] = c.GetRateValue(model.ExpMLNonMasked, model.StrML, "successful", model.CG)
		values[3][1] = c.GetRateValue(model.ExpMLNonMasked, model.StrML, "successful", model.V2X)
		values[3][2] = c.GetRateValue(model.ExpMLNonMasked, model.StrML, "successful", model.UAV)
		values[3][3] = 0 // TODO: confidence

		err := c.genRatesPkgAggregatedApps("successful-search", values, basePath)
		return err
	default:
		return errors.New("chartType not found")
	}
}

func (c *Client) genRatesPkgAggregatedMecs(resType string, values [][]float64, basePath string) error {
	mecTypeLabels := []string{"City-Level", "Regional-Level", "International-Level"}

	iterFiles := []string{"algorithm1.dat", "algorithm2.dat", "algorithm3.dat", "algorithm4.dat"}

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

	_, err = script.WriteString(generateAggregatedRatesScriptMecs(resType, title, yLabel, iterFiles))
	if err != nil {
		return err
	}
	return err
}

func (c *Client) genRatesPkgAggregatedApps(ratesType string, values [][]float64, basePath string) error {
	expLabels := []string{"Optimal-Hybrid", "Heuristic-Hybrid", "ML-Masked", "ML-NonMasked"}

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
