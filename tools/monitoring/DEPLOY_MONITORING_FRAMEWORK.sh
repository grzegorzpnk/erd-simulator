#!/bin/bash

NS="monitoring"
S_NS="kube-system"
AC="apply"
AC_H="install"

declare -a clusters=("mec1" "mec2" "mec3" "mec4" "mec5" "mec6" "mec7" "mec8" "mec9")

for cluster in "${clusters[@]}"; do
  export MEC_NAME=$cluster
  KC="../../create-mec-hosts/configs/$MEC_NAME.config/master/etc/kubernetes/admin.conf"
  localKC="../create-mec-hosts/configs/$MEC_NAME.config/master/etc/kubernetes/admin.conf"
  chmod go-r "$localKC"

  echo "Creating NS[$NS] in CLUSTER[$MEC_NAME]"
  kubectl --kubeconfig "$localKC" create ns $NS

  # shellcheck disable=SC2164
  cd ./kube-state-metrics-configs
  echo "Creating metrics-server: NS[$S_NS], CLUSTER[$MEC_NAME]"
  ./deploy-metrics-server.sh $AC "$KC" $S_NS
  echo "Creating kube-state-metrics: NS[$NS], CLUSTER[$MEC_NAME]"
  ./deploy-kube-state-metrics.sh $AC_H "$KC" $NS
  cd ../

  # shellcheck disable=SC2164
  cd ./kubernetes-prometheus
  echo "Creating Prometheus deployment: NS[$NS], CLUSTER[$MEC_NAME]"
  ./deploy-prometheus.sh $AC "$KC" $NS
  cd ../

  # shellcheck disable=SC2164
  cd ./kubernetes-grafana
  echo "Creating Grafana deployment: NS[$NS], CLUSTER[$MEC_NAME]"
  ./deploy-grafana.sh $AC "$KC" $NS
  cd ../

  echo "Installing emco monitor: NS[$NS], CLUSTER[$MEC_NAME]"
  helm --kubeconfig "$localKC" "$AC_H" emco-agent -n $NS emco/monitor --set emcoTag=22.06
done

echo "DONE!"

