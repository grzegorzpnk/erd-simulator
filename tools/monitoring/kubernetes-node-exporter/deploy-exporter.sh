if [[ $# -eq 3 ]];
then
    AC=$1 # create|delete
    KC=$2
    NS=$3
    if [[ $AC != "create" && $AC != "delete" ]];
    then
        echo "Usage: ./deploy-exporter.sh <create|delete> <kubeconfig-path> <namespace>"
	exit 0
    fi
    AC=$1 # create|delete
    KC=$2
    NS=$3
    
    kubectl --kubeconfig $KC --namespace $NS $AC -f daemonset.yaml

    kubectl --kubeconfig $KC --namespace $NS $AC -f service.yaml
    
    sleep 3
    
    if [ $AC == "create" ];
    then
        echo "--- Verifying Exported daemonset..."
        kubectl --kubeconfig $KC --namespace $NS get daemonset
        
        echo "--- Verifying Exporter service..."
        kubectl --kubeconfig $KC --namespace $NS get services

	echo "--- Verifying Exporter endpoints..."
	kubectl --kubeconfig $KC --namespace $NS get endpoints
    fi
else 
    echo "Usage: ./deploy-exorter.sh <create|delete> <kubeconfig-path> <namespace>"
fi
