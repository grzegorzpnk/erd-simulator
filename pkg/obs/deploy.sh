#!/bin/bash

helm --kubeconfig ~/.kube/control.config uninstall obs

sleep 1

cd deployments/helm && helm package obs/ && cd ../..

helm --kubeconfig ~/.kube/control.config install obs deployments/helm/obs-0.1.0.tgz
