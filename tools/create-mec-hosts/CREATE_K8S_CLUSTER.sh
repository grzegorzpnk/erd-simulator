#!/bin/bash

# Supported mec names in range [mec1, mec36]
# mec name is associated with NETWORK CIDR
declare -a clusters=(
"10.254.185.75,10.11.0.0/16,10.12.0.0/16,mec1"
"10.254.185.57,10.13.0.0/16,10.14.0.0/16,mec2"
"10.254.185.84,10.15.0.0/16,10.16.0.0/16,mec3"
"10.254.185.46,10.17.0.0/16,10.18.0.0/16,mec4"
"10.254.185.69,10.19.0.0/16,10.20.0.0/16,mec5"
"10.254.185.48,10.21.0.0/16,10.22.0.0/16,mec6"
"10.254.185.19,10.23.0.0/16,10.24.0.0/16,mec7"
"10.254.185.70,10.25.0.0/16,10.26.0.0/16,mec8"
"10.254.185.17,10.27.0.0/16,10.28.0.0/16,mec9")


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
