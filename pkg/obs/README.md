Prometheus Monitoring Controller (PMC)
---

## Design

- `Observability` part of PMC makes use of data provided by `kube-state-metrics (ksm)` which is deployed on each k8s cluster. These data are collected by
`Grafana Mimir` and they are fetched from there.

- `Observability` is focused on Cluster `Memory/CPU` `Requests/Limits`, but there are also Node level info (which is not documented below). We think, that
we should rely on K8s level parameters, rather than Node level resource utilization.

- `Latency` for now is generated randomly.

### Observability Endpoints

```yaml
SAMPLE URL: http://localhost:8282/v1/obs/ksm/provider/edge-provider/cluster/meh01/memory-requests
```


```go
// Get provider+cluster CPU Requsts utilization (in percentage)

localhost:8282/v1/obs/ksm/provider/{provider}/cluster/{cluster}/cpu-requests
```

```go
// Get provider+cluster CPU Limits utilization (in percentage)

localhost:8282/v1/obs/ksm/provider/{provider}/cluster/{cluster}/cpu-limits
```

```go
// Get provider+cluster MEMORY Requsts utilization (in percentage)

localhost:8282/v1/obs/ksm/provider/{provider}/cluster/{cluster}/memory-requests
```

```go
// Get provider+cluster Memory Limits utilization (in percentage)

localhost:8282/v1/obs/ksm/provider/{provider}/cluster/{cluster}/memory-limits
```

### Latency Endpoints

```yaml
SAMPLE URL: http://localhost:8282/v1/obs/ltc/cell/1/meh/edge-provider+meh01/latency-ms
```


```go
// Get (mocked) latency between Cell: cell-id AND MEC Host: meh-id

localhost:8282/v1/obs/ltc/cell/{cell-id}/meh/{meh-id}/latency-ms
```

