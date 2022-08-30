if [[ $# -eq 3 ]];
then
    AC=$1 # create|delete|apply
    KC=$2
    NS=$3 # suppose to be kube-system
    if [[ $AC != "create" && $AC != "delete" && $AC != "apply" ]];
    then
        echo "Usage: ./deploy-metrics-server.sh  <create|delete|apply> <kubeconfig-path> <namespace>"
	exit 0
    fi
    AC=$1 # create|delete|apply
    KC=$2
    NS=$3 # suppose to be kube-system
    
    kubectl --kubeconfig $KC --namespace $NS $AC -f .

    if [ $AC == "create" || $AC == "apply" ];
    then
        echo "--- Created!"
    fi
else 
    echo "Usage: ./deploy-metrics-server.sh  <create|delete|apply> <kubeconfig-path> <namespace>"
fi
