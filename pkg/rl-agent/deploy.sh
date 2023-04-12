#!/bin/bash

helm --kubeconfig ~/.kube/config uninstall rl-agent

sleep 1

cd deployments/helm && helm package rl-agent/ && cd ../..

helm --kubeconfig ~/.kube/config install rl-agent deployments/helm/rl-agent-0.1.0.tgz
