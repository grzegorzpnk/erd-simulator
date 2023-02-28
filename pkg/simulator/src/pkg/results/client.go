package results

import (
	"encoding/json"
	"fmt"
	"simu/src/config"
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
//		- et		ExperimentType
//		- strategy	string
//		- mt		MecType
func (c *Client) GetMecUtilizationAggregated(et ExperimentType, strategy string, mt MecType, resType string) float64 {
	var okResults []ExpResult

	for _, result := range c.GetResults() {
		if result.Metadata.Type == et && result.Metadata.Strategy == strategy {
			okResults = append(okResults, result)
		}
	}

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

func (c *Client) GetRateValue(expType ExperimentType, strategy string, eventType string, appType AppType) float64 {
	var okResults []ExpResult

	for _, result := range c.GetResults() {
		if result.Metadata.Type == expType && result.Metadata.Strategy == strategy {
			okResults = append(okResults, result)
		}
	}

	var pc percentageCounter

	for _, result := range okResults {
		switch strings.ToLower(eventType) {
		case "successful":
			switch appType {
			case CG:
				pc.total += float64(result.Data.Erd.Successful[CG])
				pc.total += float64(result.Data.Erd.Redundant[CG])
				pc.divisor += float64(result.Metadata.Movements)
			case V2X:
				pc.total += float64(result.Data.Erd.Successful[V2X])
				pc.total += float64(result.Data.Erd.Redundant[V2X])
				pc.divisor += float64(result.Metadata.Movements)
			case UAV:
				pc.total += float64(result.Data.Erd.Successful[UAV])
				pc.total += float64(result.Data.Erd.Redundant[UAV])
				pc.divisor += float64(result.Metadata.Movements)
			}
		case "triggering":
			switch appType {
			case CG:
				pc.total += float64(result.Data.Erd.Successful[CG])
				pc.divisor += float64(result.Metadata.Movements)
			case V2X:
				pc.total += float64(result.Data.Erd.Successful[V2X])
				pc.divisor += float64(result.Metadata.Movements)
			case UAV:
				pc.total += float64(result.Data.Erd.Successful[UAV])
				pc.divisor += float64(result.Metadata.Movements)
			}
		case "failed":
			switch appType {
			case CG:
				pc.total += float64(result.Data.Erd.Failed[CG])
				pc.divisor += float64(result.Metadata.Movements)
			case V2X:
				pc.total += float64(result.Data.Erd.Failed[V2X])
				pc.divisor += float64(result.Metadata.Movements)
			case UAV:
				pc.total += float64(result.Data.Erd.Failed[UAV])
				pc.divisor += float64(result.Metadata.Movements)
			}
		case "redundant":
			switch appType {
			case CG:
				pc.total += float64(result.Data.Erd.Redundant[CG])
				pc.divisor += float64(result.Metadata.Movements)
			case V2X:
				pc.total += float64(result.Data.Erd.Redundant[V2X])
				pc.divisor += float64(result.Metadata.Movements)
			case UAV:
				pc.total += float64(result.Data.Erd.Redundant[UAV])
				pc.divisor += float64(result.Metadata.Movements)
			}
		case "skipped":
			switch appType {
			case CG:
				pc.total += float64(result.Data.Erd.Skipped[CG])
				pc.divisor += float64(result.Metadata.Movements)
			case V2X:
				pc.total += float64(result.Data.Erd.Skipped[V2X])
				pc.divisor += float64(result.Metadata.Movements)
			case UAV:
				pc.total += float64(result.Data.Erd.Skipped[UAV])
				pc.divisor += float64(result.Metadata.Movements)
			}
		}
	}

	return pc.getPercentage()
}

func (c *Client) CollectExperimentStats(expType ExperimentType, strategy string, appsNum AppCounter, expNum int) error {
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

	res := ExpResult{
		Metadata: ExpResultsMeta{
			ExperimentId: experimentId,
			Type:         expType,
			Strategy:     strategy,
			Apps:         appsNum,
			Movements:    expNum,
		},
		Data: ExpResultsData{
			Erd:  ercResults,
			Topo: topoResults,
		},
	}

	c.AppendResult(res)

	return nil
}
