#!/bin/bash

helm --kubeconfig ~/.kube/config uninstall simu

sleep 1

cd deployments/helm && helm package simu/ && cd ../..

helm --kubeconfig ~/.kube/config install simu deployments/helm/simu-0.1.0.tgz
