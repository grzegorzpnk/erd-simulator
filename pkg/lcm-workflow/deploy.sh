#!/bin/bash

helm --kubeconfig ~/.kube/meh1.config uninstall lcm-worker lcm-client

sleep 1

helm --kubeconfig ~/.kube/meh1.config install lcm-client deployment/helm/workflowclient-0.1.0.tgz
helm --kubeconfig ~/.kube/meh1.config install lcm-worker deployment/helm/worker-0.1.0.tgz