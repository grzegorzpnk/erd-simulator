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

// Client is a Topology Client, which is aware of Topology Endpoint
// It's created per request, so topology client is aware of Current UE Cell ID
type Client struct {
	TopologyEndpoint string
	CurrentCell      model.CellId
}

// NewTopologyClient creates and returns Topology Client
func NewTopologyClient() *Client {
	return &Client{
		TopologyEndpoint: config.GetConfiguration().NMTEndpoint,
		CurrentCell:      "",
	}
}

// GetMecHostsByCellId gets a list of Mec Hosts related with current UE location (Cell ID)
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
		log.Infof("Looking for first reference ( City - Level) cluster")
		for i := 0; i < len(mhsIdentity); i++ {
			log.Infof("For CELL_ID=[%v] got CLUSTERS[%+v]", id, mhsIdentity[i].Cluster)
		}
	}

	for _, mhi := range mhsIdentity {
		mh := model.MecHost{Identity: mhi}
		mhs = append(mhs, mh)
	}

	return mhs, nil
}

// GetMecHost gets a single MEC Host with provider+cluster identifier
func (c *Client) GetMecHost(provider, cluster string) (model.MecHost, error) {
	var mh model.MecHost
	//var mhIdentity model.MecIdentity

	url, err := c.buildGetMecHostUrl(provider, cluster)
	if err != nil {
		return model.MecHost{}, err
	}

	respBody, err := getHTTPRespBody(url)

	if err := json.Unmarshal(respBody, &mh); err != nil {
		log.Errorf("[Topology] Couldn't unmarshal response body: %v", err)
		return model.MecHost{}, err
	}

	//mh = model.MecHost{Identity: mhIdentity}

	return mh, nil
}

// GetMecHostsByRegion gets all MEC Hosts, filters the by region and returns such list
// TODO: Consider to implement filtering logic on the Topology side (new endpoint)
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

		log.Infof("Looking for all MEC Hosts for the same searching zone as a new user location")
		for i := 0; i < len(mhsIdentity); i++ {
			log.Infof("For REGION=[%v] got CLUSTERS[%+v]", region, mhsIdentity[i].Cluster)
		}

	}

	for _, mhi := range mhsIdentity {
		mh := model.MecHost{Identity: mhi}
		mhs = append(mhs, mh)
	}

	return mhs, nil
}

func (c *Client) GetMecHosts() ([]model.MecHost, error) {
	var mhs []model.MecHost
	var mhsIdentity []model.MecIdentity

	url := c.buildGetAllMecHostsUrl()
	respBody, err := getHTTPRespBody(url)

	if err := json.Unmarshal(respBody, &mhsIdentity); err != nil {
		log.Errorf("[Topology] Couldn't unmarshal response body: %v", err)
		return []model.MecHost{}, err
	}

	if len(mhsIdentity) <= 0 {
		err = errors.New(fmt.Sprintf("no mec hosts found"))
		return []model.MecHost{}, err
	}

	for _, mhi := range mhsIdentity {
		mh := model.MecHost{Identity: mhi}
		mhs = append(mhs, mh)
	}

	return mhs, nil
}

// GetMecNeighbours gets a neighbours list for given Mec Host (mec)
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

// CollectResourcesInfo gets information such as latency, cpu utilization and memory utilization.
// This information are related with given MEC Host (mec) and are saved locally. Returns MEC Host object with updated information.
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
	if mem.Utilization == 0 || cpu.Utilization == 0 {
		log.Errorf("Take MEM from nmt. Utilization %v: ", mem.Utilization)
		log.Errorf("Take CPU from nmt. Utilization %v: ", cpu.Utilization)
	}

	mec.SetMemInfo(mem)

	return mec, err
}

// GetMecResource gets and returns information about CPU / Memory resources for given MEC Host
// resType acceptable values are: "cpu" and "memory"
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

// GetShortestPath gets and returns the minimum possible e2e latency for given MEC Hosts.
// E2E latency is considered between Cell (e.g. gNB) and MEC Host (application)
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

// buildGetAllMecHostsUrl returns topology endpoint to get all avaliable MEC Host
func (c *Client) buildGetAllMecHostsUrl() string {
	var endpoint string

	endpoint += c.TopologyEndpoint
	endpoint += ApiBase
	endpoint += ApiCollectionMecHosts

	return endpoint
}

// buildGetMecHostUrl returns topology endpoint to get a single MEC Host, represented as provider+cluster
func (c *Client) buildGetMecHostUrl(provider, cluster string) (string, error) {
	var endpoint string

	if provider == "" || cluster == "" {
		return "", errors.New("could not buildGeMecHostUrl: provider or cluster name empty")
	}

	endpoint = c.buildGetAllMecHostsUrl()
	endpoint += "/provider/" + provider + "/cluster/" + cluster

	return endpoint, nil
}

// buildGetShortestPathLatencyBased returns topology endpoint to get MEC latency between MecHost and given CellId
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
