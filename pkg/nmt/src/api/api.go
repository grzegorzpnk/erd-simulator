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

	baseUrl := "/nmt"
	//refactored:
	r.HandleFunc(baseUrl+"/graph/vertex", handler.createMecHostHandler).Methods("POST")
	r.HandleFunc(baseUrl+"/graph/vertex/provider/{provider}/cluster/{cluster}", handler.getMecHostHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/graph/edge", handler.createEdgeHandler).Methods("POST")
	r.HandleFunc(baseUrl+"/graph/mecHosts", handler.getAllMecHostsHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/graph/mecHost/provider/{provider}/cluster/{cluster}/cpu", handler.getMECCpu).Methods("GET")
	r.HandleFunc(baseUrl+"/graph/mecHost/provider/{provider}/cluster/{cluster}/memory", handler.getMECCpu).Methods("GET")
	r.HandleFunc(baseUrl+"/graph/mecHost/provider/{provider}/cluster/{cluster}/neighbours", handler.getMECNeighbours).Methods("GET")

	r.HandleFunc(baseUrl+"/graph/edge", handler.getEdgesHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/graph/vertex/{Id}/metrics", handler.getClusterCPUResources).Methods("GET")
	r.HandleFunc(baseUrl+"/graph/vertex/{Id}/metrics", handler.updateClusterCPUResources).Methods("PUT")

	//to refactor:

	r.HandleFunc(baseUrl+"/graph/edge/{IdSource}/{IdTarget}/metrics", updateEdgeMetrics).Methods("POST")

	return r

}
