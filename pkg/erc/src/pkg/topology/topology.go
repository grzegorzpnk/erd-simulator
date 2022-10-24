package topology

import (
	"10.254.188.33/matyspi5/erd/pkg/erc/src/config"
	log "10.254.188.33/matyspi5/erd/pkg/erc/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/model"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	ApiBase               = "/topology"
	ApiCollectionMecHosts = "/mecHosts"
	ApiCollectionCells    = "/cells"
	ApiGetMecHosts        = "/mec-hosts"
	ApiGetMecNeighbours   = "/neighbours"
	ApiGetShortestPath    = "/shortest-path"
)

type Client struct {
	TopologyEndpoint string
	CurrentCell      model.CellId
}

func NewTopologyClient() *Client {
	return &Client{
		TopologyEndpoint: config.GetConfiguration().TopologyEndpoint,
		CurrentCell:      "",
	}
}

func (c *Client) GetMecHostsByCellId(id model.CellId) ([]model.MecHost, error) {
	var mhs []model.MecHost
	var mhsIdentity []model.MecIdentity

	url, err := c.buildGetMecHostsByCellIdUrl(id)
	if err != nil {
		return []model.MecHost{}, err
	}

	respBody, err := getHTTPRespBody(url)

	if err := json.Unmarshal(respBody, &mhsIdentity); err != nil {
		log.Errorf("[Topology] Couldn't unmarshal response body: %v", err)
		return []model.MecHost{}, err
	}

	if len(mhsIdentity) <= 0 {
		err = errors.New(fmt.Sprintf("no mec hosts found for given cell: %v\n.", id))
		return []model.MecHost{}, err
	} else {
		log.Infof("Got CLUSTERS[%+v] for CELL_ID=[%v] ", mhsIdentity, id)
	}

	for _, mhi := range mhsIdentity {
		mh := model.MecHost{Identity: mhi}
		mhs = append(mhs, mh)
	}

	return mhs, nil
}

func (c *Client) GetMecHostsByRegion(region string) ([]model.MecHost, error) {
	var mhs []model.MecHost
	var tempMhs, mhsIdentity []model.MecIdentity

	url := c.buildGetAllMecHostsUrl()

	respBody, err := getHTTPRespBody(url)

	if err := json.Unmarshal(respBody, &tempMhs); err != nil {
		log.Errorf("[Topology] Couldn't unmarshal response body: %v", err)
		return []model.MecHost{}, err
	}

	if len(tempMhs) <= 0 {
		err = errors.New(fmt.Sprintf("no mec hosts found for given region: %v\n.", region))
		return []model.MecHost{}, err
	} else {
		// filter by region
		for _, cl := range tempMhs {
			if cl.Location.Region == region {
				mhsIdentity = append(mhsIdentity, cl)
			}
		}
		log.Infof("Got CLUSTERS[%+v] for REGION=[%v] ", mhsIdentity, region)
	}

	for _, mhi := range mhsIdentity {
		mh := model.MecHost{Identity: mhi}
		mhs = append(mhs, mh)
	}

	return mhs, nil
}

// GetMecNeighbours calls topology server to get neighbours list for given MecHost
func (c *Client) GetMecNeighbours(mec model.MecHost) (model.MecHost, error) {
	var mhsIdentity []model.MecIdentity

	url, err := c.buildGetNeighboursUrl(mec)
	if err != nil {
		return model.MecHost{}, err
	}

	respBody, err := getHTTPRespBody(url)
	if err := json.Unmarshal(respBody, &mhsIdentity); err != nil {
		log.Errorf("[Topology] Couldn't unmarshal response body: %v", err)
		return model.MecHost{}, err
	}

	if len(mhsIdentity) <= 0 {
		err = errors.New(fmt.Sprintf("no neihbours found for given mec: %v\n.", mec.BuildClusterEmcoFQDN()))
		return mec, err
	}

	for _, mhi := range mhsIdentity {
		mh := model.MecHost{Identity: mhi}
		mec.Neighbours = append(mec.Neighbours, &mh)

	}

	return mec, nil
}

// CollectResourcesInfo calls topology server for MecHost information such as latency,
// cpu utilization and memory utilization.
func (c *Client) CollectResourcesInfo(mec model.MecHost) (model.MecHost, error) {

	cpu, err := c.GetMecResource(model.MecCpu, mec)
	if err != nil {
		log.Warnf("[Topology] Could not get cpu utilization: %v", err)
	}
	mec.SetCpuInfo(cpu)

	mem, err := c.GetMecResource(model.MecMem, mec)
	if err != nil {
		log.Warnf("[Topology]  Could not get memory utilization: %v", err)
	}

	mec.SetMemInfo(mem)

	return mec, err
}

// GetMecResource is single method to get MecInfo (cpu | memory) from topology server
func (c *Client) GetMecResource(resType model.MecInfo, mec model.MecHost) (model.MecResInfo, error) {
	var val model.MecResInfo

	url, err := c.buildGetResourcesUrl(resType, mec)
	if err != nil {
		log.Errorf("[Topology] Couldn't build get resources url: %v", err)
		return model.MecResInfo{}, err
	}

	respBody, err := getHTTPRespBody(url)
	if err != nil {
		log.Errorf("[Topology] Error while getting response body: %v", err)
	}
	if err := json.Unmarshal(respBody, &val); err != nil {
		log.Errorf("[Topology] Couldn't unmarshal response body: %v", err)
		return model.MecResInfo{}, err
	}

	return val, nil
}

// GetShortestPath returns minimum latency for given MEC Hosts.
func (c *Client) GetShortestPath(cell model.CellId, mec model.MecHost) (float64, error) {
	var val float64

	url, err := c.buildGetShortestPathLatencyBased(mec, cell)
	if err != nil {
		return -1, err
	}

	respBody, err := getHTTPRespBody(url)
	if err := json.Unmarshal(respBody, &val); err != nil {
		log.Errorf("[Topology] couldn't unmarshal response body: %v", err)
		return -1, err
	}

	return val, nil
}

// buildGetMecHostsByCellIdUrl returns topology endpoint to get MEC hosts associated with given CellId
// e.g. /topology/cell/{cell-id}/mec-hosts
func (c *Client) buildGetMecHostsByCellIdUrl(id model.CellId) (string, error) {
	var endpoint string

	if string(id) == "" {
		return "", errors.New("cell-id is empty")
	}

	endpoint += c.TopologyEndpoint
	endpoint += ApiBase
	endpoint += ApiCollectionCells
	endpoint += "/" + string(id)
	endpoint += ApiGetMecHosts

	return endpoint, nil
}

// buildGetResourcesUrl returns topology endpoint to get MecHost resources (cpu | memory)
// e.g. /topology/mec/{mec-name}{/cpu | /memory}
func (c *Client) buildGetResourcesUrl(resType model.MecInfo, mec model.MecHost) (string, error) {
	var endpoint string

	if mec.Identity.Provider == "" || mec.Identity.Cluster == "" {
		return "", errors.New("mec-name is empty")
	}

	endpoint += c.TopologyEndpoint
	endpoint += ApiBase
	endpoint += ApiCollectionMecHosts
	endpoint += "/provider/" + mec.Identity.Provider + "/cluster/" + mec.Identity.Cluster
	endpoint += "/" + string(resType)

	return endpoint, nil
}

func (c *Client) buildGetAllMecHostsUrl() string {
	var endpoint string

	endpoint += c.TopologyEndpoint
	endpoint += ApiBase
	endpoint += ApiCollectionMecHosts

	return endpoint
}

// buildGetShortestPathLatencyBased returns topology endpoint to get MEC latency between MecHost and given CellId
// e.g. /topology/cell/{cell-id}/mec/{mec-name}/latency
func (c *Client) buildGetShortestPathLatencyBased(mec model.MecHost, cell model.CellId) (string, error) {
	var endpoint string

	if mec.Identity.Provider == "" || mec.Identity.Cluster == "" {
		return "", errors.New("mec provider|cluster is zero value")
	}
	if cell == "" {
		return "", errors.New("provided cell is zero value")
	}

	endpoint += c.TopologyEndpoint
	endpoint += ApiBase
	endpoint += ApiCollectionCells
	endpoint += "/" + string(cell)
	endpoint += ApiCollectionMecHosts
	endpoint += "/provider/" + mec.Identity.Provider + "/cluster/" + mec.Identity.Cluster
	endpoint += ApiGetShortestPath

	return endpoint, nil
}

// buildGetNeighboursUrl returns topology endpoint to get given MecHost neighbours
// e.g. /topology/mec/{mec-name}/neighbours
func (c *Client) buildGetNeighboursUrl(mec model.MecHost) (string, error) {
	var endpoint string

	if mec.Identity.Provider == "" || mec.Identity.Cluster == "" {
		return "", errors.New("mec-name is empty")
	}

	endpoint += c.TopologyEndpoint
	endpoint += ApiBase
	endpoint += ApiCollectionMecHosts
	endpoint += "/provider/" + mec.Identity.Provider + "/cluster/" + mec.Identity.Cluster
	endpoint += ApiGetMecNeighbours

	return endpoint, nil
}
