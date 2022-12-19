#!/bin/bash

helm --kubeconfig ~/.kube/core.config uninstall innot

sleep 1

cd deployments/helm && helm package innot/ && cd ../..

helm --kubeconfig ~/.kube/core.config install innot deployments/helm/innot-0.1.0.tgz