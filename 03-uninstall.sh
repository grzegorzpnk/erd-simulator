helm --kubeconfig ~/.kube/control.config uninstall relocate-worker
helm --kubeconfig ~/.kube/control.config uninstall relocate-client
helm --kubeconfig ~/.kube/control.config uninstall lcm-worker
helm --kubeconfig ~/.kube/control.config uninstall lcm-client
helm --kubeconfig ~/.kube/control.config uninstall erc
helm --kubeconfig ~/.kube/control.config uninstall obs
helm --kubeconfig ~/.kube/control.config uninstall nmt
helm --kubeconfig ~/.kube/core.config uninstall innot
