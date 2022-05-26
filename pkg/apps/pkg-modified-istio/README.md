# ???

* Please note that there helm charts are modified
* Init container is not trying to curl http://nrf-nnrf before creating NF container. The reason is that init container prevents injection of istio envoy proxy sidecar container and in case of service mesh, curl will never succed. That's why here we change init container in a way, that it wait fixed time before creating proper NF container.
* Also nrf wait fixed time until mongodb is started.
