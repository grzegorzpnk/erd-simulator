#!/bin/bash

declare -a SERVICES=("nmt" "erc" "simu")

for svc in "${SERVICES[@]}";
do
  helm --kubeconfig ~/.kube/config uninstall $svc
done
