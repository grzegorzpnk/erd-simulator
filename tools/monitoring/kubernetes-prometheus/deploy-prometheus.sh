if [[ $# -eq 3 ]];
then
    AC=$1 # create|delete|apply
    KC=$2 # kubeconfig path
    NS=$3 # namespace

    if [[ $MEC_NAME == "" ]] || [[ $MIMIR_ENDPOINT == "" ]];
    then
      echo "[PROMETHEUS] Error: MEC_NAME or MIMIR_ENDPOINT env variable not set"
      exit 0
    else
      echo "--- [PROMETHEUS] Deploying Prometheus for Cluster[$MEC_NAME]"
      envsubst < config-map-template.yaml >config-map.yaml
    fi

    if [[ $AC != "create" && $AC != "delete" && $AC != "apply" ]];
    then
        echo "[PROMETHEUS] Usage: ./deploy-prometheus.sh <create|delete|apply> <kubeconfig-path> <namespace>"
	    exit 0
    fi

    kubectl --kubeconfig $KC --namespace $NS $AC -f cluster-role.yaml
    
    kubectl --kubeconfig $KC --namespace $NS $AC -f config-map.yaml
    
    kubectl --kubeconfig $KC --namespace $NS $AC -f prometheus-deployment.yaml
    
    kubectl --kubeconfig $KC --namespace $NS $AC -f prometheus-service.yaml

     echo "--- [PROMETHEUS] Clear up: delete config-map.yaml..."
     rm config-map.yaml
else 
    echo "[PROMETHEUS] Usage: ./deploy-prometheus.sh <create|delete> <kubeconfig-path> <namespace>"
fi
