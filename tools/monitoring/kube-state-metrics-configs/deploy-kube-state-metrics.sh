if [[ $# -eq 3 ]];
then
    AC=$1 # install|uninstall
    KC=$2
    NS=$3 # suppose to be kube-system
    
    if [[ $AC != "install" && $AC != "uninstall" ]];
    then
        echo "Usage: ./deploy-kube-state-metrics.sh <install|uninstall> <kubeconfig-path> <namespace>"
	exit 0
    fi
    
    if [[ $AC == "install"  ]];
    then
        helm --kubeconfig $KC repo add prometheus-community https://prometheus-community.github.io/helm-charts
        helm --kubeconfig $KC repo update
        helm --kubeconfig $KC --namespace $NS $AC ksm prometheus-community/kube-state-metrics
    else
        helm --kubeconfig $KC --namespace $NS $AC ksm
    fi

else 
    echo "Usage: ./deploy-kube-state-metrics.sh <install|uninstall> <kubeconfig-path> <namespace>"
fi
