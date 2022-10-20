#!/bin/bash

# Supported mec names in range [mec1, mec36]
# mec name is associated with NETWORK CIDR
#declare -a clusters=(
#"10.254.185.88,10.11.0.0/16,10.12.0.0/16,mec1"
#"10.254.185.99,10.15.0.0/16,10.16.0.0/16,mec3"
#"10.254.185.96,10.17.0.0/16,10.18.0.0/16,mec4"
#"10.254.185.89,10.19.0.0/16,10.20.0.0/16,mec5"
#"10.254.185.91,10.21.0.0/16,10.22.0.0/16,mec6"
#"10.254.185.60,10.23.0.0/16,10.24.0.0/16,mec7"
#"10.254.185.93,10.31.0.0/16,10.32.0.0/16,mec11"
#"10.254.185.31,10.33.0.0/16,10.34.0.0/16,mec12"
#"10.254.185.9./02,10.35.0.0/16,10.36.0.0/16,mec13"
#"10.254.185.61,10.37.0.0/16,10.38.0.0/16,mec14"
#"10.254.185.35,10.39.0.0/16,10.40.0.0/16,mec15"
#"10.254.185.75,10.41.0.0/16,10.42.0.0/16,mec16"
#"10.254.185.49,10.43.0.0/16,10.44.0.0/16,mec17"
#"10.254.185.85,10.45.0.0/16,10.46.0.0/16,mec18"
#"10.254.185.17,10.47.0.0/16,10.48.0.0/16,mec19"
#"10.254.185.70,10.49.0.0/16,10.50.0.0/16,mec20"
#"10.254.185.19,10.51.0.0/16,10.52.0.0/16,mec21"
#"10.254.185.48,10.53.0.0/16,10.54.0.0/16,mec22"
#"10.254.185.69,10.55.0.0/16,10.56.0.0/16,mec23"
#"10.254.185.46,10.57.0.0/16,10.58.0.0/16,mec24"
#"10.254.185.57,10.59.0.0/16,10.60.0.0/16,mec25"
#"10.254.185.84,10.61.0.0/16,10.62.0.0/16,mec26"
#)

declare -a clusters=(
"10.254.185.88,10.11.0.0/16,10.12.0.0/16,mec1"
"10.254.185.99,10.15.0.0/16,10.16.0.0/16,mec3"
"10.254.185.96,10.17.0.0/16,10.18.0.0/16,mec4"
"10.254.185.89,10.19.0.0/16,10.20.0.0/16,mec5"
"10.254.185.91,10.21.0.0/16,10.22.0.0/16,mec6"
"10.254.185.60,10.23.0.0/16,10.24.0.0/16,mec7"
)


for cluster in "${clusters[@]}"; do
    IFS=',' read -ra PARAMS <<< "$cluster"
    node_ip=${PARAMS[0]}  && export NODE_IP=$node_ip   
    network_cidr=${PARAMS[1]} && export NETWORK_CIDR=$network_cidr
    service_cidr=${PARAMS[2]} && export SERVICE_CIDR=$service_cidr
    node_name=${PARAMS[3]} && export MEC_NAME=$node_name

    echo "Generating files from template: hosts, 03-master.yaml & 04-configure.yaml files for cluster $node_ip"
    envsubst < ./playbooks/03-init-template.yaml >playbooks/03-init.yaml
    envsubst < ./playbooks/04-configure-template.yaml >playbooks/04-configure.yaml
    envsubst < ./hosts-template >hosts

    echo "Run Ansible playbooks"
    ansible-playbook -i ./hosts ./playbooks/01-initial.yaml
    ansible-playbook -i ./hosts ./playbooks/02-kube-dependencies.yaml
    ansible-playbook -i ./hosts ./playbooks/03-init.yaml

    # wait until new cluster is healthy to apply configuration via kubectl
    sleep 10
    ansible-playbook -i ./hosts ./playbooks/04-configure.yaml
done

echo "DONE! Clear up hosts, 03-init.yaml & 04-configure.yaml"
rm ./hosts ./playbooks/03-init.yaml ./playbooks/04-configure.yaml
