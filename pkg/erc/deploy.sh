#!/bin/bash

helm --kubeconfig "$ER_KUBECONFIG" uninstall erc

sleep 1

cd deployments/helm && helm package erc/ && cd ../..

helm --kubeconfig "$ER_KUBECONFIG" install erc deployments/helm/erc-0.1.0.tgz
