package results

import (
	"encoding/json"
	"fmt"
	"math"
	"simu/src/config"
	log "simu/src/logger"
	"simu/src/pkg/model"
	"strconv"
	"strings"
)

var experimentId = 1

type Client struct {
	expResults []ExpResult
}

func NewClient() *Client {
	return &Client{
		expResults: []ExpResult{},
	}
}

func (c *Client) GetResults() []ExpResult {
	return c.expResults
}

func (c *Client) AppendResult(result ExpResult) {
	c.expResults = append(c.expResults, result)
}

func (c *Client) IncExpId() {
	experimentId++
}

func (c *Client) GetMecUtilizationSingleExperiment(expId int, mecType MecType) float64 {
	var util float64
	fmt.Println("NOT IMPLEMENTED")
	return util
}

// GetMecUtilizationAggregated returns aggregated value of resource (cpu/memory) utilization for specific:
//   - et		ExperimentType
//   - strategy	string
//   - mt		MecType
func (c *Client) GetMecUtilizationAggregated(et model.ExperimentType, strategy model.ExperimentStrategy, mt MecType, resType string) float64 {
	var okResults []ExpResult

	for _, result := range c.GetResults() {
		if result.Metadata.Type == et && result.Metadata.Strategy == strategy {
			okResults = append(okResults, result)
		}
	}

	log.Infof("[GenerateCharts][GetMecUtilAggr] Found %v results for [Exp:%v][Str:%v][Mec:%v][Res:%v]",
		len(okResults), et, strategy, mt, resType)

	var pc percentageCounter

	for _, result := range okResults {
		switch mt {
		case MecLocal:
			for _, mecHost := range result.Data.Topo.MecHostsResults {
				switch strings.ToLower(resType) {
				case "cpu":
					if mecHost.Identity.Location.Level == MecLocal {
						pc.total += mecHost.CpuResources.Used
						pc.divisor += mecHost.CpuResources.Capacity
					}
				case "memory":
					if mecHost.Identity.Location.Level == MecLocal {
						pc.total += mecHost.MemoryResources.Used
						pc.divisor += mecHost.MemoryResources.Capacity
					}
				}
			}
		case MecRegional:
			for _, mecHost := range result.Data.Topo.MecHostsResults {
				switch strings.ToLower(resType) {
				case "cpu":
					if mecHost.Identity.Location.Level == MecRegional {
						pc.total += mecHost.CpuResources.Used
						pc.divisor += mecHost.CpuResources.Capacity
					}
				case "memory":
					if mecHost.Identity.Location.Level == MecRegional {
						pc.total += mecHost.MemoryResources.Used
						pc.divisor += mecHost.MemoryResources.Capacity
					}
				}
			}
		case MecCentral:
			for _, mecHost := range result.Data.Topo.MecHostsResults {
				switch strings.ToLower(resType) {
				case "cpu":
					if mecHost.Identity.Location.Level == MecCentral {
						pc.total += mecHost.CpuResources.Used
						pc.divisor += mecHost.CpuResources.Capacity
					}
				case "memory":
					if mecHost.Identity.Location.Level == MecCentral {
						pc.total += mecHost.MemoryResources.Used
						pc.divisor += mecHost.MemoryResources.Capacity
					}
				}
			}
		}
	}

	return pc.getPercentage()
}

func (c *Client) GetRateValueForExperiment(expId int, eventType string) float64 {
	var result ExpResult
	var pc percentageCounter

	for _, exp := range c.GetResults() {
		if exp.Metadata.ExperimentId == expId {
			result = exp
			break
		}
	}

	// Make sure, that this condition is OK. If expId is 0, it means that we didn't find experiment with given ID
	if result.Metadata.ExperimentId == 0 {
		log.Errorf("[GenerateCharts][GetRateValueForExperiment] Could not find experiment [ID:%v]", expId)
		return 0
	}

	switch strings.ToLower(eventType) {
	case "successful":
		pc.total += float64(result.Data.Erd.Successful[model.CG])
		pc.total += float64(result.Data.Erd.Redundant[model.CG])

		pc.total += float64(result.Data.Erd.Successful[model.V2X])
		pc.total += float64(result.Data.Erd.Redundant[model.V2X])

		pc.total += float64(result.Data.Erd.Successful[model.UAV])
		pc.total += float64(result.Data.Erd.Redundant[model.UAV])

		pc.divisor += float64(result.Metadata.Movements)
	case "triggering":
		pc.total += float64(result.Data.Erd.Successful[model.CG])

		pc.total += float64(result.Data.Erd.Successful[model.V2X])

		pc.total += float64(result.Data.Erd.Successful[model.UAV])

		pc.divisor += float64(result.Metadata.Movements)

	case "failed":
		pc.total += float64(result.Data.Erd.Failed[model.CG])

		pc.total += float64(result.Data.Erd.Failed[model.V2X])

		pc.total += float64(result.Data.Erd.Failed[model.UAV])

		pc.divisor += float64(result.Metadata.Movements)

	case "redundant":
		pc.total += float64(result.Data.Erd.Redundant[model.CG])

		pc.total += float64(result.Data.Erd.Redundant[model.V2X])

		pc.total += float64(result.Data.Erd.Redundant[model.UAV])

		pc.divisor += float64(result.Metadata.Movements)

	case "skipped":
		pc.total += float64(result.Data.Erd.Skipped[model.CG])

		pc.total += float64(result.Data.Erd.Skipped[model.V2X])

		pc.total += float64(result.Data.Erd.Skipped[model.UAV])

		pc.divisor += float64(result.Metadata.Movements)
	}

	return pc.getPercentage()
}

func (c *Client) GetRateValueAllIter(et model.ExperimentType, strategy model.ExperimentStrategy, eventType string, appType model.AppType) float64 {
	var okResults []ExpResult

	for _, result := range c.GetResults() {
		if result.Metadata.Type == et && result.Metadata.Strategy == strategy {
			okResults = append(okResults, result)
		}
	}

	log.Infof("[GenerateCharts][GetRateValueAllIter] Found %v results for [Exp:%v][Str:%v][Eve:%v][App:%v]",
		len(okResults), et, strategy, eventType, appType)

	var pc percentageCounter

	for _, result := range okResults {
		switch strings.ToLower(eventType) {
		case "successful":
			switch appType {
			case model.CG:
				pc.total += float64(result.Data.Erd.Successful[model.CG])
				pc.total += float64(result.Data.Erd.Redundant[model.CG])
				pc.divisor += float64(result.Metadata.Movements)
			case model.V2X:
				pc.total += float64(result.Data.Erd.Successful[model.V2X])
				pc.total += float64(result.Data.Erd.Redundant[model.V2X])
				pc.divisor += float64(result.Metadata.Movements)
			case model.UAV:
				pc.total += float64(result.Data.Erd.Successful[model.UAV])
				pc.total += float64(result.Data.Erd.Redundant[model.UAV])
				pc.divisor += float64(result.Metadata.Movements)
			}
		case "triggering":
			switch appType {
			case model.CG:
				pc.total += float64(result.Data.Erd.Successful[model.CG])
				pc.divisor += float64(result.Metadata.Movements)
			case model.V2X:
				pc.total += float64(result.Data.Erd.Successful[model.V2X])
				pc.divisor += float64(result.Metadata.Movements)
			case model.UAV:
				pc.total += float64(result.Data.Erd.Successful[model.UAV])
				pc.divisor += float64(result.Metadata.Movements)
			}
		case "failed":
			switch appType {
			case model.CG:
				pc.total += float64(result.Data.Erd.Failed[model.CG])
				pc.divisor += float64(result.Metadata.Movements)
			case model.V2X:
				pc.total += float64(result.Data.Erd.Failed[model.V2X])
				pc.divisor += float64(result.Metadata.Movements)
			case model.UAV:
				pc.total += float64(result.Data.Erd.Failed[model.UAV])
				pc.divisor += float64(result.Metadata.Movements)
			}
		case "redundant":
			switch appType {
			case model.CG:
				pc.total += float64(result.Data.Erd.Redundant[model.CG])
				pc.divisor += float64(result.Metadata.Movements)
			case model.V2X:
				pc.total += float64(result.Data.Erd.Redundant[model.V2X])
				pc.divisor += float64(result.Metadata.Movements)
			case model.UAV:
				pc.total += float64(result.Data.Erd.Redundant[model.UAV])
				pc.divisor += float64(result.Metadata.Movements)
			}
		case "skipped":
			switch appType {
			case model.CG:
				pc.total += float64(result.Data.Erd.Skipped[model.CG])
				pc.divisor += float64(result.Metadata.Movements)
			case model.V2X:
				pc.total += float64(result.Data.Erd.Skipped[model.V2X])
				pc.divisor += float64(result.Metadata.Movements)
			case model.UAV:
				pc.total += float64(result.Data.Erd.Skipped[model.UAV])
				pc.divisor += float64(result.Metadata.Movements)
			}
		}
	}

	return pc.getPercentage()
}

func (c *Client) GetAverageConvTimeAllIter(et model.ExperimentType, strategy model.ExperimentStrategy) float64 {
	var okResults []ExpResult

	for _, result := range c.GetResults() {
		if result.Metadata.Type == et && result.Metadata.Strategy == strategy {
			okResults = append(okResults, result)
		}
	}

	var pc percentageCounter

	for _, result := range okResults {
		for i := 0; i < len(result.Data.Erd.EvalTimes.Failed); i++ {
			pc.total += float64(result.Data.Erd.EvalTimes.Failed[i])
			pc.divisor += 1
		}
		for i := 0; i < len(result.Data.Erd.EvalTimes.Skipped); i++ {
			pc.total += float64(result.Data.Erd.EvalTimes.Skipped[i])
			pc.divisor += 1
		}
		for i := 0; i < len(result.Data.Erd.EvalTimes.Redundant); i++ {
			pc.total += float64(result.Data.Erd.EvalTimes.Redundant[i])
			pc.divisor += 1
		}
		for i := 0; i < len(result.Data.Erd.EvalTimes.Successful); i++ {
			pc.total += float64(result.Data.Erd.EvalTimes.Successful[i])
			pc.divisor += 1
		}

	}

	return pc.getAverage()
}

func (c *Client) GetConfidenceValue(et model.ExperimentType, es model.ExperimentStrategy, evt string) (confidence float64) {

	var okResults []ExpResult
	var xn []float64
	var rootSum float64

	for _, result := range c.GetResults() {
		if result.Metadata.Type == et && result.Metadata.Strategy == es {
			okResults = append(okResults, result)
		}
	}

	log.Infof("[GenerateCharts][GetConfidenceValue] Found %v results for [Exp:%v][Str:%v][Eve:%v][App:ALL]",
		len(okResults), et, es, evt)

	// Calculate average value for all experiments
	cgAvg := c.GetRateValueAllIter(et, es, evt, model.CG)
	v2xAvg := c.GetRateValueAllIter(et, es, evt, model.V2X)
	uavAvg := c.GetRateValueAllIter(et, es, evt, model.UAV)

	// Calculate average for all applications
	avg := cgAvg + v2xAvg + uavAvg
	log.Warnf("[DEBUG][GetConfidenceValue] Avg = %v", avg)

	// Calculate percentage values for each experiment in range (1, N)
	for _, result := range okResults {
		x := c.GetRateValueForExperiment(result.Metadata.ExperimentId, evt)
		xn = append(xn, x)
	}

	// Calculate standard deviation
	for i := 0; i < len(xn); i++ {
		rootSum += math.Pow(xn[i]-avg, 2)
	}

	std := math.Sqrt(rootSum / float64(len(xn)))
	log.Warnf("[DEBUG][GetConfidenceValue] Std = %v", std)

	// Calculate confidence level
	z := 1.96 // z value for 95% confidence is 1.96

	confidence = z * std / math.Sqrt(float64(len(xn)))
	log.Warnf("[DEBUG][GetConfidenceValue] Confidence = %v", confidence)

	return
}

func (c *Client) CollectExperimentStats(exp model.ExperimentIntent) error {

	var ercResults ErdResults
	var topoResults TopoResults

	ercUrl := config.GetConfiguration().ERCEndpoint + "/v2/erc/results"
	topoUrl := config.GetConfiguration().NMTEndpoint + "/v1/topology/mecHosts/metrics"

	ercBody, err := getHttpRespBody(ercUrl)
	if err != nil {
		return err
	}

	err = json.Unmarshal(ercBody, &ercResults)
	if err != nil {
		return err
	}

	topoBody, err := getHttpRespBody(topoUrl)
	if err != nil {
		return err
	}

	err = json.Unmarshal(topoBody, &topoResults.MecHostsResults)
	if err != nil {
		return err
	}

	mov, err := strconv.Atoi(exp.ExperimentDetails.MovementsInExperiment)
	if err != nil {
		return err
	}

	res := ExpResult{
		Metadata: ExpResultsMeta{
			ExperimentId: experimentId,
			Type:         exp.ExperimentType,
			Strategy:     exp.ExperimentStrategy,
			Apps:         exp.ExperimentDetails.InitialAppsNumber,
			Movements:    mov,
		},
		Data: ExpResultsData{
			Erd:  ercResults,
			Topo: topoResults,
		},
	}

	c.AppendResult(res)

	return nil
}
