#!/bin/bash

helm --kubeconfig ~/.kube/config uninstall nmt

sleep 1

cd deployments/helm && helm package nmt/ && cd ../..

helm --kubeconfig ~/.kube/config install nmt deployments/helm/nmt-0.1.0.tgz
