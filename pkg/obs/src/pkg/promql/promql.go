package promql

import (
	log "10.254.188.33/matyspi5/erd/pkg/obs/src/logger"
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

func (pql *PromQL) GetCurrentRequests(resType Resource, targetCluster string) (float64, error) {
	var fVal float64

	query := fmt.Sprintf("sum(kube_pod_container_resource_requests{cluster=~\"%s\",resource=\"%v\",node!~\"^(.*master).*$\"})", targetCluster, resType)

	val, err := pql.Query(query)
	if err != nil {
		log.Errorf("[PromQL] Error: %v", err)
		return -1, err
	}
	if len(val) > 0 {
		fVal, err = strconv.ParseFloat(val[0].Value.String(), 64)
		if err != nil {
			log.Errorf("[PromQL] Could not ParseFloat. Reason: %v", err)
			return -1, err
		}
	} else {
		log.Errorf("Error while fetching %v requests for cluster: %v. reason: %v.", resType, targetCluster, err)
	}

	return fVal, nil
}

func (pql *PromQL) GetCurrentLimits(resType Resource, targetCluster string) (float64, error) {
	var fVal float64

	query := fmt.Sprintf("sum(kube_pod_container_resource_limits{cluster=~\"%s\",resource=\"%v\",node!~\"^(.*master).*$\"})", targetCluster, resType)

	val, err := pql.Query(query)
	if err != nil {
		log.Errorf("[PromQL] Error: %v", err)
		return -1, err
	}
	if len(val) > 0 {
		fVal, err = strconv.ParseFloat(val[0].Value.String(), 64)
		if err != nil {
			log.Errorf("[PromQL] Could not ParseFloat. Reason: %v", err)
			return -1, err
		}
	} else {
		log.Errorf("Error while fetching %v limits for cluster: %v. reason: %v.", resType, targetCluster, err)
	}

	return fVal, nil
}

func (pql *PromQL) GetAllocatable(resType Resource, targetCluster string) (float64, error) {
	var fVal float64

	query := fmt.Sprintf("sum(kube_node_status_allocatable{cluster=~\"%s\",resource=\"%v\",node!~\"^(.*master).*$\"})", targetCluster, resType)

	val, err := pql.Query(query)
	if err != nil {
		log.Errorf("[PromQL] Error: %v", err)
		return -1, err
	}
	if len(val) > 0 {
		fVal, err = strconv.ParseFloat(val[0].Value.String(), 64)
		if err != nil {
			log.Errorf("[PromQL] Could not ParseFloat. Reason: %v", err)
			return -1, err
		}
	} else {
		log.Errorf("Error while fetching %v allocatable for cluster: %v. reason: %v.", resType, targetCluster, err)
	}

	return fVal, nil
}

// GetCpuRequestsLimits returns percentage of utilized CPU requests and CPU limits at cluster targetCluster
// Query is based on `kube-state-metrics` data-source. Skips nodes with master in the hostname
func (pql *PromQL) GetCpuRequestsLimits(targetCluster string) (float64, float64, error) {

	limits := fmt.Sprintf("sum(kube_pod_container_resource_limits{cluster=~\"%s\",resource=\"cpu\",node!~\"^(.*master).*$\"})", targetCluster)
	requests := fmt.Sprintf("sum(kube_pod_container_resource_requests{cluster=~\"%s\",resource=\"cpu\",node!~\"^(.*master).*$\"})", targetCluster)
	allocatable := fmt.Sprintf("sum(kube_node_status_allocatable{cluster=~\"%s\",resource=\"cpu\",node!~\"^(.*master).*$\"})", targetCluster)

	var lim, req, alloc float64

	val, err := pql.Query(limits)
	if err != nil {
		log.Errorf("[PromQL] Error: %v", err)
		return -1, -1, err
	}
	if len(val) > 0 {
		lim, err = strconv.ParseFloat(val[0].Value.String(), 64)
		if err != nil {
			log.Errorf("[PromQL] Could not ParseFloat. Reason: %v", err)
			return -1, -1, err
		}
	} else {
		log.Errorf("Error while fetching CPU limits: %v", err)
	}

	val, err = pql.Query(requests)
	if err != nil {
		log.Errorf("[PromQL] Error: %v", err)
		return -1, -1, err
	}
	if len(val) > 0 {
		req, err = strconv.ParseFloat(val[0].Value.String(), 64)
		if err != nil {
			log.Errorf("[PromQL] Could not ParseFloat. Reason: %v", err)
			return -1, -1, err
		}
	} else {
		log.Errorf("Error while fetching CPU requests: %v", err)
	}

	val, err = pql.Query(allocatable)
	if err != nil {
		log.Errorf("[PromQL] Error: %v", err)
		return -1, -1, err
	}
	if len(val) > 0 {
		alloc, err = strconv.ParseFloat(val[0].Value.String(), 64)
		if err != nil {
			log.Errorf("[PromQL] Could not ParseFloat. Reason: %v", err)
			return -1, -1, err
		}
	} else {
		log.Errorf("Error while fetching CPU allocatable: %v", err)
	}

	return 100 * (req / alloc), 100 * (lim / alloc), nil
}

// GetMemoryRequestsLimits returns percentage of utilized CPU requests and CPU limits at cluster targetCluster
// Query is based on `kube-state-metrics` data-source. Skips nodes with master in the hostname
func (pql *PromQL) GetMemoryRequestsLimits(targetCluster string) (float64, float64, error) {
	//log.Infof("Getting MEMORY (requests, limits) for cluster %s", targetCluster)
	limits := fmt.Sprintf("sum(kube_pod_container_resource_limits{cluster=~\"%s\",resource=\"memory\",node!~\"^(.*master).*$\"})", targetCluster)
	requests := fmt.Sprintf("sum(kube_pod_container_resource_requests{cluster=~\"%s\",resource=\"memory\",node!~\"^(.*master).*$\"})", targetCluster)
	allocatable := fmt.Sprintf("sum(kube_node_status_allocatable{cluster=~\"%s\",resource=\"memory\",node!~\"^(.*master).*$\"})", targetCluster)

	var lim, req, alloc float64

	val, err := pql.Query(limits)
	if err != nil {
		log.Errorf("[PromQL] Error: %v", err)
		return -1, -1, err
	}
	if len(val) > 0 {
		lim, err = strconv.ParseFloat(val[0].Value.String(), 64)
		if err != nil {
			log.Errorf("[PromQL] Could not ParseFloat. Reason: %v", err)
			return -1, -1, err
		}
	} else {
		log.Errorf("Error while fetching MEMORY limits: %v", err)
	}

	val, err = pql.Query(requests)
	if err != nil {
		log.Errorf("[PromQL] Error: %v", err)
		return -1, -1, err
	}
	if len(val) > 0 {
		req, err = strconv.ParseFloat(val[0].Value.String(), 64)
		if err != nil {
			log.Errorf("[PromQL] Could not ParseFloat. Reason: %v", err)
			return -1, -1, err
		}
	} else {
		log.Errorf("Error while fetching MEMORY requests: %v", err)
	}

	val, err = pql.Query(allocatable)
	if err != nil {
		log.Errorf("[PromQL] Error: %v", err)
		return -1, -1, err
	}
	if len(val) > 0 {
		alloc, err = strconv.ParseFloat(val[0].Value.String(), 64)
		if err != nil {
			log.Errorf("[PromQL] Could not ParseFloat. Reason: %v", err)
			return -1, -1, err
		}
	} else {
		log.Errorf("Error while fetching MEMORY allocatable: %v", err)
	}

	return 100 * (req / alloc), 100 * (lim / alloc), nil
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
