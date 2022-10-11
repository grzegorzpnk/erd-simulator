if [[ $# -eq 3 ]];
then
    AC=$1 # create|delete
    KC=$2
    NS=$3

    if [[ $AC != "create" && $AC != "delete" && $AC != "apply" ]];
    then
        echo "[GRAFANA] Usage: ./deploy-grafana.sh <apply|create|delete> <kubeconfig-path> <namespace>"
	      exit 0
    fi

    echo "--- [GRAFANA] Creating Grafana deployment..."
    kubectl --kubeconfig $KC --namespace $NS $AC -f .

else 
    echo "[GRAFANA] Usage: ./deploy-grafana.sh <apply|create|delete> <kubeconfig-path> <namespace>"
fi
