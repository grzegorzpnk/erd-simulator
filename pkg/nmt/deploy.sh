#!/bin/bash

helm --kubeconfig ~/.kube/core.config uninstall nmt

sleep 1

cd deployments/helm && helm package nmt/ && cd ../..

helm --kubeconfig ~/.kube/core.config install nmt deployments/helm/nmt-0.1.0.tgz