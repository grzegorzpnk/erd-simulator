#!/bin/bash

helm --kubeconfig ~/.kube/core.config uninstall lcm-worker lcm-client

sleep 1

cd deployments/helm && helm package worker/ workflowclient/ && cd ../..

helm --kubeconfig ~/.kube/core.config install lcm-client deployments/helm/workflowclient-0.1.0.tgz
helm --kubeconfig ~/.kube/core.config install lcm-worker deployments/helm/worker-0.1.0.tgz
