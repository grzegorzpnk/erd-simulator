#!/bin/bash

helm --kubeconfig ~/.kube/ran.config uninstall obs

sleep 1

helm --kubeconfig ~/.kube/ran.config install obs deployments/helm/obs-0.1.0.tgz
