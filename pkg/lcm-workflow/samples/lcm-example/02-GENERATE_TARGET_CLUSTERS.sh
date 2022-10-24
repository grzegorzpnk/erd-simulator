#!/bin/bash

function run() {

for app in "${apps[@]}"; do
    index=$(expr $(($RANDOM % (${#clusters[@]}))))
    targetCluster=${clusters[$index]}
    export "TARGET_CLUSTER_$app"="$targetCluster"
done
  envsubst < ./emco-manifests-v2/values-no-clusters.yaml >./emco-manifests-v2/values.yaml
}

declare -a apps=("1" "2" "3" "4" "5" "6" "7" "8" "9" "10"
"11" "12" "13" "14" "15" "16" "17" "18" "19" "20" "21" "22" "23" "24" "25" "26" "27" "28" "29" "30"
"31" "32" "33" "34" "35" "36" "37" "38" "39" "40" "41" "42" "43" "44" "45" "46" "47" "48" "49" "50"
"51" "52" "53" "54" "55" "56" "57" "58" "59" "60" "61" "62" "63" "64" "65" "66" "67" "68" "69" "70"
"71" "72" "73" "74" "75" "76" "77" "78" "79" "80" "81" "82" "83" "84" "85" "86" "87" "88" "89" "90"
"91" "92" "93" "94" "95" "96" "97" "98" "99" "100")

declare -a clusters=("mec1" "mec3" "mec4" "mec5" "mec6" "mec7" "mec11" "mec12" "mec13" "mec14" "mec15" "mec16" "mec17" "mec18" "mec19" "mec20" "mec21" "mec22" "mec23" "mec24" "mec25" "mec26")

run
