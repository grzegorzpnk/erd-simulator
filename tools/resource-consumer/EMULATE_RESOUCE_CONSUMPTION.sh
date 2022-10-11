#!/bin/bash

NS="res-consumer"
AC=$1 # apply | delete

declare -a clusters=("mec1" "mec2" "mec3" "mec4" "mec5" "mec6" "mec7" "mec8" "mec9")

for cluster in "${clusters[@]}"; do
    export MEC_NAME=$cluster
    KC="$HOME/workshop/orange/erd/tools/create-mec-hosts/configs/$MEC_NAME.config/master/etc/kubernetes/admin.conf"
    chmod go-r "$KC"

    if [ $AC == "apply" ];
    then
        echo "Creating NS[$NS] in CLUSTER[$MEC_NAME]"
        kubectl --kubeconfig "$KC" create ns $NS

        export REQ_MEMORY="$(expr $((1 + $RANDOM % 10)) \* 128)Mi"
        export REQ_CPU="$(expr $((1 + $RANDOM % 10)) \* 125)m"
    fi

    echo "Generating manifest for CLUSTER[$MEC_NAME]"
    envsubst < pod-template.yaml > manifest.yaml

    kubectl --kubeconfig $KC --namespace $NS $AC -f manifest.yaml

    if [ $AC == "delete" ];
    then
        echo "Deleting NS[$NS] in CLUSTER[$MEC_NAME]"
        kubectl --kubeconfig "$KC" delete ns $NS
    fi

done
echo "DONE!"