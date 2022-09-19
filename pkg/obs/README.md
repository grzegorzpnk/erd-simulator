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
SAMPLE URL: http://10.254.185.50:32138/v1/obs/ksm/provider/edge-provider/cluster/mec1/cpu/requests
```


```go
// Get provider+cluster CPU/Memory allocated requests

localhost:8282/v1/obs/ksm/provider/{provider}/cluster/{cluster}/cpu/requests
localhost:8282/v1/obs/ksm/provider/{provider}/cluster/{cluster}/memory/requests
```

```go
// Get provider+cluster CPU/Memory allocated limits

localhost:8282/v1/obs/ksm/provider/{provider}/cluster/{cluster}/cpu/limits
localhost:8282/v1/obs/ksm/provider/{provider}/cluster/{cluster}/memory/limits
```

```go
// Get provider+cluster CPU/Memory utilized requests (in percentage)

localhost:8282/v1/obs/ksm/provider/{provider}/cluster/{cluster}/cpu/utilization
localhost:8282/v1/obs/ksm/provider/{provider}/cluster/{cluster}/memory/utilization
```

### Latency Endpoints

```yaml
SAMPLE URL: http://localhost:8282/v1/obs/ltc/cell/1/mec/edge-provider+mec1/latency-ms
```


```go
// Get (mocked) latency between Cell: cell-id AND MEC Host: mec-id

localhost:8282/v1/obs/ltc/cell/{cell-id}/mec/{mec-id}/latency-ms
```

