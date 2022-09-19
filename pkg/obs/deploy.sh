#!/bin/bash

helm --kubeconfig ~/.kube/ran.config uninstall obs

sleep 1

cd deployments/helm && helm package obs/ && cd ../..

helm --kubeconfig ~/.kube/ran.config install obs deployments/helm/obs-0.1.0.tgz
