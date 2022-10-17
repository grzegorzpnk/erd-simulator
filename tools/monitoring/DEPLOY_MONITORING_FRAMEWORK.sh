#!/bin/bash

NS="monitoring"
NS_S="kube-system"
AC="apply"
AC_H="install"
emcoTag="22.06"
export MIMIR_ENDPOINT=10.254.185.27

declare -a clusters=("mec1" "mec2" "mec3" "mec4" "mec5" "mec6" "mec7" "mec11" "mec12" "mec13" "mec14" "mec15" "mec16" "mec17" "mec18" "mec19" "mec20" "mec21" "mec22" "mec23" "mec24" "mec25" "mec26")

for cluster in "${clusters[@]}"; do
  export MEC_NAME=$cluster
  KC="$HOME/workshop/orange/erd/tools/create-mec-hosts/configs/$MEC_NAME.config/master/etc/kubernetes/admin.conf"
  chmod go-r "$KC"

  echo "Creating NS[$NS] in CLUSTER[$MEC_NAME]"
  kubectl --kubeconfig "$KC" create ns $NS

  # shellcheck disable=SC2164
  cd ./kube-state-metrics-configs
  echo "Installing metrics-server: NS[$NS_S], CLUSTER[$MEC_NAME]"
  ./deploy-metrics-server.sh $AC "$KC" $NS_S
  echo "Installing kube-state-metrics: NS[$NS], CLUSTER[$MEC_NAME]"
  ./deploy-kube-state-metrics.sh $AC_H "$KC" $NS
  cd ../

  # shellcheck disable=SC2164
  cd ./kubernetes-prometheus
  echo "Installing Prometheus: NS[$NS], CLUSTER[$MEC_NAME]"
  ./deploy-prometheus.sh $AC "$KC" $NS
  cd ../

  # shellcheck disable=SC2164
  cd ./kubernetes-grafana
  echo "Installing Grafana: NS[$NS], CLUSTER[$MEC_NAME]"
  ./deploy-grafana.sh $AC "$KC" $NS
  cd ../

  echo "Installing emco monitor: NS[$NS], CLUSTER[$MEC_NAME]"
  helm --kubeconfig "$KC" "$AC_H" emco-agent -n $NS emco/monitor --set emcoTag=$emcoTag
done

echo "DONE!"

