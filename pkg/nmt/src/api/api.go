package api

import (
	"github.com/gorilla/mux"
	"nmt/src/package/topology"
)

var r *mux.Router

func NewRouter(graphClient *topology.Graph) *mux.Router {

	var handler apiHandler
	handler.SetClients(*graphClient)

	r := mux.NewRouter().PathPrefix("/v1").Subrouter()

	baseUrl := "/nmt"

	r.HandleFunc(baseUrl+"/graph/vertex", handler.getAllVertexesHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/graph/vertex/{Id}", handler.getVertexHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/graph/vertex/{Id}/metrics", handler.updateClusterMetrics).Methods("PUT")
	r.HandleFunc(baseUrl+"/graph/vertex/{Id}/metrics", handler.getClusterMetrics).Methods("GET")
	r.HandleFunc(baseUrl+"/graph/vertex", handler.createVertex).Methods("POST")

	r.HandleFunc(baseUrl+"/graph/edge", handler.getEdgesHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/graph/edge", handler.createEdgeHandler).Methods("POST")
	r.HandleFunc(baseUrl+"/graph/edge/{IdSource}/{IdTarget}/metrics", updateEdgeMetrics).Methods("POST")

	return r

}
