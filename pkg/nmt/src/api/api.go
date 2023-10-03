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

	//Create infrastructure:
	r.HandleFunc(baseUrl+"/topology/mecHost", handler.createMecHostHandler).Methods("POST")
	r.HandleFunc(baseUrl+"/topology/edge", handler.createLinkHandler).Methods("POST")

	//Communication with Placement Controller:
	r.HandleFunc(baseUrl+"/topology/mecHosts", handler.getAllMecHostsHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/topology/mecHosts/provider/{provider}/cluster/{cluster}", handler.getMecHostHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/topology/cells/{cell-id}/mec-hosts", handler.getCellAssociatedMecHostsHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/topology/cells/{cell-id}/mecHosts/provider/{provider}/cluster/{cluster}/shortest-path", handler.shortestPathHandler).Methods("GET")

	r.HandleFunc(baseUrl+"/topology/mecHosts/metrics", handler.getAllMecHostsWithMetricsHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/topology/mecHosts/provider/{provider}/cluster/{cluster}/cpu", handler.getMECCpu).Methods("GET")
	r.HandleFunc(baseUrl+"/topology/mecHosts/provider/{provider}/cluster/{cluster}/memory", handler.getMECMemory).Methods("GET")
	r.HandleFunc(baseUrl+"/topology/mecHosts/provider/{provider}/cluster/{cluster}/neighbours", handler.getMECNeighbours).Methods("GET")

	// Communication with Orchestrator
	r.HandleFunc(baseUrl+"/topology/application", handler.getAllMecApps).Methods("GET")
	r.HandleFunc(baseUrl+"/topology/mecHosts/application", handler.getAllMecHostsWithApps).Methods("GET")
	r.HandleFunc(baseUrl+"/topology/mecHosts/provider/{provider}/cluster/{cluster}/application", handler.InstantiateApplication).Methods("POST")
	r.HandleFunc(baseUrl+"/topology/mecHosts/provider/{provider}/cluster/{cluster}/removeApplication", handler.DeleteApplication).Methods("POST")
	r.HandleFunc(baseUrl+"/orchestrator/relocate/old-cluster/{old-cluster}/new-cluster/{new-cluster}/application", handler.RelocateApplication2).Methods("POST")

	//r.HandleFunc(baseUrl+"/topology/declare-apps/{applications}", handler.DefineApplications).Methods("POST")
	r.HandleFunc(baseUrl+"/topology/initial-placement-generator", handler.GenerateInitialClusters).Methods("POST")
	r.HandleFunc(baseUrl+"/topology/instantiate-all", handler.InstantiateAllDefinedApps).Methods("POST")
	//this API defines one-click (DefineApplications & GenerateInitialClusters & InstantiateAllDefinedApps)
	r.HandleFunc(baseUrl+"/topology/prerequesties/{applications}", handler.Prerequisites).Methods("POST")
	r.HandleFunc(baseUrl+"/topology/prerequesties", handler.PrerequisitesTunning).Methods("POST")
	r.HandleFunc(baseUrl+"/topology/recreate-initial", handler.Recreate).Methods("POST")

	// interaction with ML training app (currently not used). Used only to create first config input in order to train on
	r.HandleFunc(baseUrl+"/topology/ml/InitialState/{applications}", handler.MLInitialState).Methods("GET")
	r.HandleFunc(baseUrl+"/topology/ml/InitialConfig", handler.MLInitialConfig).Methods("GET")
	r.HandleFunc(baseUrl+"/topology/ml/rans", handler.MLGetRANs).Methods("GET")

	//interaction with Simulator in order to preapre ML smart placement intent  for current state of MECs
	r.HandleFunc(baseUrl+"/topology/ml/get-state", handler.GetCurrentState).Methods("POST")
	r.HandleFunc(baseUrl+"/topology/ml/get-mask", handler.GetCurrentMask).Methods("POST")

	return r

}
