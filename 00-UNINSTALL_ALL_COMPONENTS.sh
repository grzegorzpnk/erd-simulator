#!/bin/bash

declare -a SERVICES=("nmt" "erc" "rl-agent")

for svc in "${SERVICES[@]}";
do
  helm --kubeconfig ~/.kube/config uninstall $svc
done
