if [[ $# -eq 3 ]];
then
    AC=$1 # create|delete|apply
    KC=$2
    NS=$3 # suppose to be kube-system

    if [[ $AC != "create" && $AC != "delete" && $AC != "apply" ]];
    then
        echo "[METRICS-SERVER] Usage: ./deploy-metrics-server.sh  <create|delete|apply> <kubeconfig-path> <namespace>"
	      exit 0
    fi

    echo "--- [METRICS-SERVER] Installing metrics-server..."
    kubectl --kubeconfig "$KC" --namespace "$NS" "$AC" -f .

else 
    echo "[METRICS-SERVER] Usage: ./deploy-metrics-server.sh  <create|delete|apply> <kubeconfig-path> <namespace>"
fi
