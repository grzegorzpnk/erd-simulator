package api

import (
	"github.com/gorilla/mux"
	"nmt/src/package/mec-topology"
)

var r *mux.Router

func NewRouter(graphClient *mec_topology.Graph) *mux.Router {

	var handler apiHandler
	handler.SetClients(*graphClient)

	r := mux.NewRouter().PathPrefix("/v1").Subrouter()

	baseUrl := "/nmt"
	//refactored:
	r.HandleFunc(baseUrl+"/graph/vertex", handler.createMecHostHandler).Methods("POST")
	r.HandleFunc(baseUrl+"/graph/vertex/{Id}", handler.getMecHostHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/graph/edge", handler.createEdgeHandler).Methods("POST")
	r.HandleFunc(baseUrl+"/graph/vertex", handler.getAllMecHostsHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/graph/edge", handler.getEdgesHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/graph/vertex/{Id}/metrics", handler.getClusterCPUResources).Methods("GET")
	r.HandleFunc(baseUrl+"/graph/vertex/{Id}/metrics", handler.updateClusterCPUResources).Methods("PUT")

	//to refactor:

	r.HandleFunc(baseUrl+"/graph/edge/{IdSource}/{IdTarget}/metrics", updateEdgeMetrics).Methods("POST")

	return r

}
