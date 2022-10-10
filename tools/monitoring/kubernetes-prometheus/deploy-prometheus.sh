if [[ $# -eq 3 ]];
then
    AC=$1 # create|delete|apply
    KC=$2 # kubeconfig path
    NS=$3 # namespace

    if [[ $MEC_NAME == "" ]];
    then
      echo "Error: MEC_NAME env variable not set"
      exit 0
    else
      echo "Deploying prometheus for Cluster[$MEC_NAME]"
      envsubst < config-map-template.yaml >config-map.yaml
    fi

    if [[ $AC != "create" && $AC != "delete" && $AC != "apply" ]];
    then
        echo "Usage: ./deploy-prometheus.sh <create|delete|apply> <kubeconfig-path> <namespace>"
	    exit 0
    fi

    # kubectl --kubeconfig $KC $AC ns $NS
    kubectl --kubeconfig $KC --namespace $NS $AC -f cluster-role.yaml
    
    kubectl --kubeconfig $KC --namespace $NS $AC -f config-map.yaml
    
    kubectl --kubeconfig $KC --namespace $NS $AC -f prometheus-deployment.yaml
    
    kubectl --kubeconfig $KC --namespace $NS $AC -f prometheus-service.yaml
    
    sleep 5
    
    if [ $AC == "create" ] || [ $AC == "apply" ];
    then
        echo "--- Verifying Prometheus deployment..."
        kubectl --kubeconfig $KC --namespace $NS get deployments
        
        echo "--- Verifying Prometheus service..."
        kubectl --kubeconfig $KC --namespace $NS get services
    fi

     echo "--- Clear up: delete config-map.yaml..."
     rm config-map.yaml
else 
    echo "Usage: ./deploy-prometheus.sh <create|delete> <kubeconfig-path> <namespace>"
fi
