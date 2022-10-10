if [[ $# -eq 3 ]];
then
    AC=$1 # create|delete
    KC=$2
    NS=$3
    if [[ $AC != "create" && $AC != "delete" && $AC != "apply" ]];
    then
        echo "Usage: ./deploy-grafana.sh <apply|create|delete> <kubeconfig-path> <namespace>"
	      exit 0
    fi

    echo "--- Creating Grafana deployment..."
    kubectl --kubeconfig $KC --namespace $NS $AC -f .
    
    sleep 5
    
    if [ $AC == "create" ] || [ $AC == "apply" ];
    then
        echo "--- Verifying Grafana deployment..."
        kubectl --kubeconfig $KC --namespace $NS get deployments | grep grafana
        
        echo "--- Verifying Grafana service..."
        kubectl --kubeconfig $KC --namespace $NS get services | grep grafana
    fi
else 
    echo "Usage: ./deploy-grafana.sh <apply|create|delete> <kubeconfig-path> <namespace>"
fi
