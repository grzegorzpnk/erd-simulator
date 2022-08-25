package api

import (
	"github.com/gorilla/mux"
	"nmt/src/package/topology"
)

var r *mux.Router

func NewRouter(graphClient *topology.Graph) *mux.Router {

	var handler apiHandler
	handler.SetClients(*graphClient)

	r = mux.NewRouter()

	r.HandleFunc("/graph/vertex", handler.getAllVertexesHandler).Methods("GET")
	r.HandleFunc("/graph/vertex/{Id}", handler.getVertexHandler).Methods("GET")
	r.HandleFunc("/graph/vertex/{Id}/metrics", handler.updateClusterMetrics).Methods("PUT")
	r.HandleFunc("/graph/vertex", handler.createVertex).Methods("POST")

	r.HandleFunc("/graph/edge", handler.getEdgesHandler).Methods("GET")
	r.HandleFunc("/graph/edge", handler.createEdgeHandler).Methods("POST")
	r.HandleFunc("/graph/edge/{IdSource}/{IdTarget}/metrics", updateEdgeMetrics).Methods("POST")

	return r

}
