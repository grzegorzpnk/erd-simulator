Prerequisites
---

- Install EMCO - use `./emco/deployments`
- Install EMCO MONITOR on each cluster - use `./monitoring/monitor`

`Note: Use provided emco and monitor packages, because the newest packages are not correctly notifying about deployment rediness status - TODO resolve problem.`

- PowerDNS by default is running on machine `10.254.185.50`
- Temporal Server by default is running on cluster meh1 (`10.254.185.48`) in the `temporal-server` namespace
- Temporal Worker and Client by default are running on cluster meh1 (`10.254.185.48`) in the demo namespace
- Img-server is deployed by default on the cluster meh1 (`10.254.185.48`)
- Img-server is relocated by default to the cluster meh2 (`10.254.185.27`
)
Demo
---

1. Install and configure MariaDB and PowerDNS. Make sure that you can access PowerDNS from outside - via external IP address (remember to register SOA).
2. Install and configure externalDNS on each MEC Host
	2.1 adjust interval (interval=1s)
	2.2 adjust policy (policy=sync)
	2.3 adjust provider (provider=pdns)
	2.4 set pdns-server
	2.5 set pdns-api-key
	2.6 for each MEC Host set-up different txt-owner-id (that MEC Hosts can't overwrite existing entries)
	2.7 ...
3. Use `watch -n 0.5 sudo pdnsutil list-zone dtf-demo.com` command on the pdns node to watch existing entries
4. Install temporal server (and THEN exec into admintools pod and register namespace `default` using `tctl namespace register default` also expose web service as NodePort to access GUI.)

```
kubectl --kubeconfig ~/.kube/meh1.config create ns temporal-server
helm -n temporal-server install     --set server.replicaCount=1     --set cassandra.config.cluster_size=1     --set prometheus.enabled=false     --set grafana.enabled=false     --set elasticsearch.enabled=false     temporal-server ./temporal-server/ --timeout 15m --kubeconfig ~/.kube/meh1.config
```
5. Install `client` and `worker`

```
kubectl --kubeconfig ~/.kube/meh1.config create ns demo
helm --kubeconfig ~/.kube/meh1.config --namespace demo install demo  ./temporal-migrate-workflow/deployment/helm/workflowclient
helm --kubeconfig ~/.kube/meh1.config --namespace demo install demo1 ./temporal-migrate-workflow/deployment/helm/worker/

```

6. Deploy img-server

```bash
temporal-migrate-workflow/samples/relocate-example/emco-manifests/DEPLOY.sh
```

7. Turn on img-client

```bash
./demo-apps/RUN_CLIENT.sh
```

8. Perform relocation

```bash
./temporal-migrate-workflow/samples/relocate-example/emco-manifests/RELOCATE_TO_CLUSTER2.sh
./temporal-migrate-workflow/samples/relocate-example/emco-manifests/RELOCATE_TO_CLUSTER1.sh
```

