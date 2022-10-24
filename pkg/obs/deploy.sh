#!/bin/bash

helm --kubeconfig ~/.kube/core.config uninstall obs

sleep 1

cd deployments/helm && helm package obs/ && cd ../..

helm --kubeconfig ~/.kube/core.config install obs deployments/helm/obs-0.1.0.tgz
