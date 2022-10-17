#!/bin/bash

helm --kubeconfig ~/.kube/control.config uninstall relocate-client relocate-worker

sleep 1

cd deployments/helm && helm package worker/ workflowclient/ && cd ../..

helm --kubeconfig ~/.kube/control.config install relocate-client deployments/helm/workflowclient-0.1.0.tgz
helm --kubeconfig ~/.kube/control.config install relocate-worker deployments/helm/worker-0.1.0.tgz
