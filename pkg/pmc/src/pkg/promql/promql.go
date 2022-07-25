package promql

import (
	log "10.254.188.33/matyspi5/pmc/src/logger"
	"context"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"strconv"
	"time"
)

// Client is Prometheus client
type Client interface {
	v1.API
}

// NewClient creates and returns new prometheus client
func NewClient(pEndpoint string) (v1.API, error) {
	pClient, err := api.NewClient(api.Config{
		Address: pEndpoint,
	})
	if err != nil {
		log.Errorf("[PROMQL] Could not create prometheus client! Error: %v", err)
		return *new(v1.API), err
	}
	return v1.NewAPI(pClient), nil
}

// PromQL represents promql object
type PromQL struct {
	Host    string
	Timeout time.Duration
	Time    time.Time
	Client  v1.API
}

// Query queries PromQL using Prometheus api and returns value
func (pql *PromQL) Query(q string) (model.Vector, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pql.Timeout)
	defer cancel()

	result, _, err := pql.Client.Query(ctx, q, pql.Time)
	if err != nil {
		log.Errorf("[PROMQL] error while querying. Query: %v. Err: %v", q, err)
		return model.Vector{}, err
	}
	if value, ok := result.(model.Vector); ok {
		return value, nil
	} else {
		err = errors.New("error converting model.Value to model.Vector")
		return model.Vector{}, err
	}
}

//// Query queries PromQL using Prometheus api and returns value
//func (pql *PromQL) Targets() (model.Vector, error) {
//	ctx, cancel := context.WithTimeout(context.Background(), pql.Timeout)
//	defer cancel()
//
//	targets, err := pql.Client.Targets(ctx)
//	if err != nil {
//		log.Errorf("Error while getting Targets. Err: %v", err)
//		return model.Vector{}, err
//	}
//	fmt.Println(targets)
//	return model.Vector{}, nil
//	//if value, ok := result.(model.Vector); ok {
//	//	return value, nil
//	//} else {
//	//	err = errors.New("error converting model.Value to model.Vector")
//	//	return model.Vector{}, err
//	//}
//}

// GetValidNodes is just for testing purposes. It returns all matching targetNodes info.
// TODO the question is how to make a mapping between cluster and instance while fetching targets.
func (pql *PromQL) GetValidNodes() {
	job := "node-exporter"

	nodesQuery := fmt.Sprintf("node_uname_info{job=\"%s\"}", job)
	nodes, _ := pql.Query(nodesQuery)

	fmt.Printf("Nodes: %s\n", nodes)
}

func (pql *PromQL) GetCpuUtilisationNatively(targetNodes string) (float64, float64) {
	limits := fmt.Sprintf("sum(kube_pod_container_resource_limits{cluster=~\"\",resource=\"cpu\",node=~\"(%s)\"})", targetNodes)
	requests := fmt.Sprintf("sum(kube_pod_container_resource_requests{cluster=~\"\",resource=\"cpu\",node=~\"(%s)\"})", targetNodes)
	allocatable := fmt.Sprintf("sum(kube_node_status_allocatable{node=~\"(%s)\",cluster=~\"\",resource=\"cpu\"})", targetNodes)

	var lim, req, alloc float64

	val, _ := pql.Query(limits)
	lim, _ = strconv.ParseFloat(val[0].Value.String(), 64)

	val, _ = pql.Query(requests)
	req, _ = strconv.ParseFloat(val[0].Value.String(), 64)

	val, _ = pql.Query(allocatable)
	alloc, _ = strconv.ParseFloat(val[0].Value.String(), 64)

	return 100 * (req / alloc), 100 * (lim / alloc)
}

// GetCpuUtilisation returns cpu utilisation on targetNode, as a percentage of used resources.
// targetNode suppose to be a node-exported pod endpoint, eg. 10.41.1.168:9100
func (pql *PromQL) GetCpuUtilisation(targetNode string) (string, error) {
	query := fmt.Sprintf("(((count(count(node_cpu_seconds_total{instance=\"%[1]v\",job=\"node-exporter\"}) by (cpu))) - avg(sum by (mode)(rate(node_cpu_seconds_total{mode='idle',instance=\"%[1]v\",job=\"node-exporter\"}[5m15s])))) * 100) / count(count(node_cpu_seconds_total{instance=\"%[1]v\",job=\"node-exporter\"}) by (cpu))", targetNode)
	val, err := pql.Query(query)
	if err != nil {
		log.Errorf("[PROMQL] Could not GetCpuUtilisation parameter for %v. Reason: %v", targetNode, err)
		return "", err
	}
	if len(val) == 0 {
		err = errors.New(fmt.Sprintf("Could not GetCpuUtilisation parameter for %v. Reason: value slice is empty.", targetNode))
		return "", err
	}
	value := val[0].Value.String()
	return value, nil
}

// GetRamUtilisation returns cpu utilisation on targetNode, as a percentage of used resources.
// targetNode suppose to be a node-exported pod endpoint, eg. 10.41.1.168:9100
func (pql *PromQL) GetRamUtilisation(targetNode string) (string, error) {
	query := fmt.Sprintf("100 - ((node_memory_MemAvailable_bytes{instance=\"%[1]v\",job=\"node-exporter\"} * 100) / node_memory_MemTotal_bytes{instance=\"%[1]v\",job=\"node-exporter\"})", targetNode)
	val, err := pql.Query(query)
	if err != nil {
		log.Errorf("[PROMQL] Could not GetRamUtilisation parameter for %v. Reason: %v", targetNode, err)
		return "", err
	}
	if len(val) == 0 {
		err = errors.New(fmt.Sprintf("Could not GetRamUtilisation parameter for %v. Reason: value slice is empty.", targetNode))
		return "", err
	}
	value := val[0].Value.String()
	return value, nil
}
