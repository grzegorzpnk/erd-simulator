#!/bin/bash

function generateMaxLatency() {
  echo "$(expr $((1 + $RANDOM % 20)))"
}

function generateMaxResUtil() {
   echo "$(expr $((50 + $RANDOM % 50)))"
}

function mebibytesToBytes() {
  echo "$(expr 1048576 \* $1)"
}

function miliCpuToCpu() {
  echo "scale=3; $1 / 1000" | bc -l
}

function generateAppChart() {
  cd ./apps || exit
  cp -r ./edge-app-template ./edge-app-"$app"

  export APP_NUMBER=$app
  envsubst < ./edge-app-"$app"/Chart.yaml > ./edge-app-"$app"/temp.yaml
  mv ./edge-app-"$app"/temp.yaml ./edge-app-"$app"/Chart.yaml

  setRequests $1 $2

  helm package ./edge-app-"$app"
  rm -rf ./edge-app-"$app"
  cd ../
}

function setRequests() {
  export REQUESTS_CPU=$1m
  export REQUESTS_MEM=$2Mi

  envsubst < ./edge-app-"$app"/values.yaml > ./edge-app-"$app"/temp.yaml
  mv ./edge-app-"$app"/temp.yaml ./edge-app-"$app"/values.yaml

  echo "App[edge-app-$app] requests:"
  cat ./edge-app-"$app"/values.yaml | grep "cpu"
  cat ./edge-app-"$app"/values.yaml | grep "memory"

}

function run() {

declare -a vars=("TARGET_CLUSTER_" "MAX_LATENCY_" "CPU_UTIL_MAX_" "MEM_UTIL_MAX_" "LTC_WEIGHT_" "RES_WEIGHT_" "CPU_WEIGHT_" "MEM_WEIGHT_" "APP_CPU_REQ_" "APP_MEM_REQ_")

for app in "${apps[@]}"; do
#    targetCluster=${clusters[$index]}
#    index=$(expr $(($RANDOM % (${#clusters[@]}))))
    cpu_req=$(expr $((125 + $RANDOM % 200)))
    mem_req=$(expr $((200 + $RANDOM % 200)))
    mem_req_bytes=$(mebibytesToBytes "$mem_req")
    cpu_req_unit="0"$(miliCpuToCpu "$cpu_req")

    generateAppChart "$cpu_req" "$mem_req"

    for var in "${vars[@]}"; do
      temp="$var$app"
      if [[ $var == "TARGET_CLUSTER_" ]]; then
        export "$temp"="\$TARGET_CLUSTER_$app"
      elif [[ $var == "MAX_LATENCY_" ]]; then
        export "$temp"="$(generateMaxLatency)"
      elif [[ $var == "CPU_UTIL_MAX_" ]]; then
        export "$temp"="$(generateMaxResUtil)"
      elif [[ $var == "MEM_UTIL_MAX_" ]]; then
        export "$temp"="$(generateMaxResUtil)"
      elif [[ $var == "LTC_WEIGHT_" ]]; then
        export "$temp"=0.0
      elif [[ $var == "RES_WEIGHT_" ]]; then
        export "$temp"=1.0
      elif [[ $var == "CPU_WEIGHT_" ]]; then
        export "$temp"=0.5
      elif [[ $var == "MEM_WEIGHT_" ]]; then
        export "$temp"=0.5
      elif [[ $var == "APP_CPU_REQ_" ]]; then
        export "$temp"=$cpu_req_unit
      elif [[ $var == "APP_MEM_REQ_" ]]; then
        export "$temp"=$mem_req_bytes
      fi
    done
  done
  envsubst < ./emco-manifests-v2/values-template.yaml >./emco-manifests-v2/values-no-clusters.yaml
}

declare -a apps=("1" "2" "3" "4" "5" "6" "7" "8" "9" "10"
"11" "12" "13" "14" "15" "16" "17" "18" "19" "20" "21" "22" "23" "24" "25" "26" "27" "28" "29" "30"
"31" "32" "33" "34" "35" "36" "37" "38" "39" "40" "41" "42" "43" "44" "45" "46" "47" "48" "49" "50"
"51" "52" "53" "54" "55" "56" "57" "58" "59" "60" "61" "62" "63" "64" "65" "66" "67" "68" "69" "70"
"71" "72" "73" "74" "75" "76" "77" "78" "79" "80" "81" "82" "83" "84" "85" "86" "87" "88" "89" "90"
"91" "92" "93" "94" "95" "96" "97" "98" "99" "100")

declare -a clusters=("mec1" "mec3" "mec4" "mec5" "mec6" "mec7" "mec11" "mec12" "mec13" "mec14" "mec15" "mec16" "mec17" "mec18" "mec19" "mec20" "mec21" "mec22" "mec23" "mec24" "mec25" "mec26")

run
