#!/bin/bash

helm --kubeconfig ~/.kube/meh2.config uninstall relocate-client relocate-worker

sleep 1

helm --kubeconfig ~/.kube/meh2.config install relocate-client deployment/helm/workflowclient-0.1.0.tgz
helm --kubeconfig ~/.kube/meh2.config install relocate-worker deployment/helm/worker-0.1.0.tgz
