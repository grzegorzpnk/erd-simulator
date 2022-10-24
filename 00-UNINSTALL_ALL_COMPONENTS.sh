#!/bin/bash

declare -a SERVICES=("relocate-worker" "relocate-client" "lcm-worker" "lcm-client" "erc" "obs" "nmt" "innot")

for svc in "${SERVICES[@]}";
do
  helm --kubeconfig ~/.kube/core.config uninstall relocate-worker $svc
done
