for kc in $(ls); do if [[ $kc == *.config ]]; then echo "-------------------- $kc -------------------" && kubectl --kubeconfig $kc get po; fi; done
