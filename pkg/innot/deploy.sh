#!/bin/bash

helm --kubeconfig ~/.kube/core.config uninstall innot

sleep 1

helm --kubeconfig ~/.kube/core.config install innot deployments/helm/innot-0.1.0.tgz
