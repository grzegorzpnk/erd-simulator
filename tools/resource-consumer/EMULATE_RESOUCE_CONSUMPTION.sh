#!/bin/bash

NS="res-consumer"
AC=$1 # apply | delete

# Declare all clusters in a coverage zone, including all levels [0, 1, 2], to distinguish between resource utilization
declare -a clusters0=("mec11" "mec12" "mec13" "mec14" "mec15" "mec16" "mec17" "mec18" "mec19" "mec20" "mec21" "mec22" "mec23" "mec24" "mec25" "mec26")
declare -a clusters1=("mec2" "mec3" "mec4" "mec5" "mec6" "mec7")
#declare -a clusters2=("mec1")

for cluster in "${clusters0[@]}"; do
    export MEC_NAME=$cluster
    KC="$HOME/workshop/orange/erd/tools/create-mec-hosts/configs/$MEC_NAME.config/master/etc/kubernetes/admin.conf"
    chmod go-r "$KC"

    if [ "$AC" == "apply" ];
    then
        echo "Creating NS[$NS] in CLUSTER[$MEC_NAME]"
        kubectl --kubeconfig "$KC" create ns $NS

        export REQ_MEMORY="$(expr $((4 + $RANDOM % 6)) \* 128)Mi"
        export REQ_CPU="$(expr $((4 + $RANDOM % 6)) \* 125)m"
    fi

    echo "Generating manifest for CLUSTER[$MEC_NAME]"
    envsubst < pod-template.yaml > manifest.yaml

    kubectl --kubeconfig $KC --namespace $NS $AC -f manifest.yaml

    if [ "$AC" == "delete" ];
    then
        echo "Deleting NS[$NS] in CLUSTER[$MEC_NAME]"
        kubectl --kubeconfig "$KC" delete ns $NS
    fi

done

for cluster in "${clusters1[@]}"; do
    export MEC_NAME=$cluster
    KC="$HOME/workshop/orange/erd/tools/create-mec-hosts/configs/$MEC_NAME.config/master/etc/kubernetes/admin.conf"
    chmod go-r "$KC"

    if [ "$AC" == "apply" ];
    then
        echo "Creating NS[$NS] in CLUSTER[$MEC_NAME]"
        kubectl --kubeconfig "$KC" create ns $NS

        export REQ_MEMORY="$(expr $((1 + $RANDOM % 3)) \* 128)Mi"
        export REQ_CPU="$(expr $((1 + $RANDOM % 3)) \* 125)m"
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