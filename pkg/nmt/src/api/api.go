package api

import (
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/mec-topology"
	"github.com/gorilla/mux"
)

var r *mux.Router

func NewRouter(graphClient *mec_topology.Graph) *mux.Router {

	var handler apiHandler
	handler.SetClients(*graphClient)

	r := mux.NewRouter().PathPrefix("/v1").Subrouter()

	baseUrl := ""
	//refactored:
	r.HandleFunc(baseUrl+"/graph/vertex", handler.createMecHostHandler).Methods("POST")
	r.HandleFunc(baseUrl+"/graph/vertex/provider/{provider}/cluster/{cluster}", handler.getMecHostHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/graph/edge", handler.createLinkHandler).Methods("POST")
	//r.HandleFunc(baseUrl+"/graph/edge", handler.getEdgesHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/graph/vertex/{Id}/metrics", handler.getClusterCPUResources).Methods("GET")
	r.HandleFunc(baseUrl+"/graph/vertex/{Id}/metrics", handler.updateClusterCPUResources).Methods("PUT")

	//Communication with Placement Controller:

	r.HandleFunc(baseUrl+"/topology/cells/{cell-id}/mec-hosts", handler.getCellAssociatedMecHostsHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/topology/cells/{cell-id}/mecHosts/provider/{provider}/cluster/{cluster}/shortest-path", handler.shortestPathHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/topology/mecHosts", handler.getAllMecHostsHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/topology/mecHosts/metrics", handler.getAllMecHostsWithMetricsHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/topology/mecHosts/provider/{provider}/cluster/{cluster}", handler.getMecHostHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/topology/mecHosts/provider/{provider}/cluster/{cluster}/cpu", handler.getMECCpu).Methods("GET")
	r.HandleFunc(baseUrl+"/topology/mecHosts/provider/{provider}/cluster/{cluster}/memory", handler.getMECMemory).Methods("GET")
	r.HandleFunc(baseUrl+"/topology/mecHosts/provider/{provider}/cluster/{cluster}/neighbours", handler.getMECNeighbours).Methods("GET")

	//to refactor:

	r.HandleFunc(baseUrl+"/graph/edge/{IdSource}/{IdTarget}/metrics", updateEdgeMetrics).Methods("POST")

	return r

}
