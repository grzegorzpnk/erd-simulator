#!/bin/bash

helm --kubeconfig ~/.kube/core.config uninstall erc

sleep 1

cd deployments/helm && helm package erc/ && cd ../..

helm --kubeconfig ~/.kube/core.config install erc deployments/helm/erc-0.1.0.tgz
