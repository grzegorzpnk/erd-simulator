package results

import (
	"encoding/json"
	"simu/src/config"
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

// AppendResult TODO: check if experiment-id is unique
func (c *Client) AppendResult(result ExpResult) {
	c.expResults = append(c.expResults, result)
}

func (c *Client) IncExpId() {
	experimentId++
}

func (c *Client) CollectExperimentStats(iterId int, expType ExperimentType, appNum, expNum int) error {
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
			Timestamp:    iterId,
			Type:         expType,
			Apps:         appNum,
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
