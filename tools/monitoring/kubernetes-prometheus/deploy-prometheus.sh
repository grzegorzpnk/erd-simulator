if [[ $# -eq 3 ]];
then
    AC=$1 # create|delete|apply
    KC=$2
    NS=$3

    if [[ $AC != "create" && $AC != "delete" && $AC != "apply" ]];
    then
        echo "Usage: ./deploy-prometheus.sh <create|delete|apply> <kubeconfig-path> <namespace>"
	exit 0
    fi
    AC=$1 # create|delete
    KC=$2
    NS=$3
    
    # Create ns monitoring
    # kubectl --kubeconfig $KC $AC ns $NS
    
    kubectl --kubeconfig $KC --namespace $NS $AC -f cluster-role.yaml
    
    kubectl --kubeconfig $KC --namespace $NS $AC -f config-map.yaml
    
    kubectl --kubeconfig $KC --namespace $NS $AC -f prometheus-deployment.yaml
    
    kubectl --kubeconfig $KC --namespace $NS $AC -f prometheus-service.yaml
    
    sleep 3
    
    if [ $AC == "create" ];
    then
        echo "--- Verifying Prometheus deployment..."
        kubectl --kubeconfig $KC --namespace $NS get deployments
        
        echo "--- Verifying Prometheus service..."
        kubectl --kubeconfig $KC --namespace $NS get services
    fi
else 
    echo "Usage: ./deploy-prometheus.sh <create|delete> <kubeconfig-path> <namespace>"
fi
