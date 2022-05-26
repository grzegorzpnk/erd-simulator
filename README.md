Towards5gs DEMO - multicluster deployment
---

This demo allows to deploy free5gc CNF's and test connectivity with UERANSIM in multi-cluster scenario.

## Requirements:
 * Make sure that your environment meets configuration described in section [Clusters Configuration](#clusters-configuration)
 * Make sure that all of the clusters are in the same subnet.
    * For now we use Multus with MACVLAN CNI - in that case we have to make sure that all `worker nodes` can directly communicate with each other, because SMF <-> UPF and AMF <-> gNB communication is assured via virtual MACVLAN networks. This virtual networks (virtual ports - which can be named slave ports) send traffic via master ports (master port is "physical" port on the VM - in our example ens3 is a master port)
    * In the futere we would like to change this solution to be more flexible.


## Notes

* EMCO manifests are created and working properly
* Configuration should be applied via `profile` files
* Before deploying gNB and UE network functions, we should register UE via webui service
* All k8s clusters must be on the same network for the demo to work properly
* By default we assume that all control plane CNF's are deployed on one cluster. UPF, gNB and UE are deployed on the second cluster
* Provided manifests allow you to add a third cluster by editing the `values.yaml` file and configuring proper clusters in the `3-deployment.yaml` file
* Communication between AMF <-> gNB is allowed via NGAP NodePort. This endopoints can communicatie even if k8s clusters are in the different subnets. In present Towards5gs implementation communication between UPF <-> SMF are restricted to single network.

## Scripts

*All scripts should be called inside `emco/exemples/towards5gc-demo/root` directory*

* `automate-testing.sh` - run migration test scenario 20 times. Append results to `result.log` file.

```
./automate-testing.sh $kube-context1 $kube-context2 $kube-context3
```

* `test-timing.sh` - run migration test scenario once. Append results to `result.log` file.

```
./test-timing.sh $kube-context1 $kube-context2 $kube-context3
```

* `clear-up.sh` - delete all deployed resources

```
./clear-up.sh
```

* `free5gc-deploy.sh` - deploy all emco resources but instantiate only free5gc cnf's

```
./free5gc-deploy.sh
```

* `ueransim-deploy.sh` - if ueransim resources are already deployed, instantiate ueransim cnf's

```
./ueransim-deploy.sh
```

* `dig-update.sh` - update deployment intent group

```
./dig-update.sh
```

* `first-pi-update.sh` - update placement intent that UPF should be created in the new location (don't delete existing instance)

```
./first-pi-update.sh
```

* `second-pi-update.sh` - update placement intent that old instance of UPF should be deleted

```
./second-pi-update.sh
```

## Clusters configuration

To provide connectivity between all CNF's we have to perform additional configuration, which can depend on your environment. Here we present configuration steps which were performed on the following infrastructure:

* Ubuntu 20.04 LTS, kernel 5.4.0-90-generic VMs running on the OpenStack
* Kubernetes v1.22.2, with Flannel CNI. We have used ansible & kubeadm to create clusters manually.
* Multus will be used to provide connectivity between CNF's


`Note1:  `*There shouldn't be any problems with Linux kernel version 5.0.0-x. Requirements for specific kernel are due to some changes in kernel and compatibility with gtp5g module. If you're able to successfully install gtp5g module, then everything should be alright*

`Note2: `*Kubernetes should support SCTP . This is only requirement on K8s version that we are aware of*

### Prepare virtual machines

```bash
# Update & Upgrade
sudo apt -y update
sudo apt -y upgrade
```

```bash
# Need to be verified if this step is required
sudo apt install linux-image-extra-virtual
```

```bash
# Install golang
wget https://dl.google.com/go/go1.14.4.linux-amd64.tar.gz
sudo tar -C /usr/local -zxvf go1.14.4.linux-amd64.tar.gz
mkdir -p ~/go/{bin,pkg,src}
echo "The following assume that your shell is bash"
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export GOROOT=/usr/local/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin:$GOROOT/bin' >> ~/.bashrc
echo 'export GO111MODULE=auto' >> ~/.bashrc
source ~/.bashrc
```

```bash
# Install helm
curl https://baltocdn.com/helm/signing.asc | sudo apt-key add -
sudo apt-get install apt-transport-https --yes
echo "deb https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list
sudo apt-get update
sudo apt-get install helm
```

```bash
# Install Control Plane packages
sudo apt -y update
sudo apt -y install mongodb wget git
sudo systemctl start mongodb
```

```bash
# Install User Plane packages
sudo apt -y update
sudo apt -y install git gcc g++ cmake autoconf libtool pkg-config libmnl-dev libyaml-dev
go get -u github.com/sirupsen/logrus
```

```bash
# Install gtp5g kernel module
git clone -b v0.4.0 https://github.com/free5gc/gtp5g.git
cd gtp5g
make clean & make
sudo make install
```

```bash
# Install additional tools. Need to be verified if this step is required
sudo apt -y install lksctp-tools
sudo apt -y install ntp
```

Add below commands to the /etc/sysctl.conf file

```bash
sudo vim /etc/sysctl.conf
```

```bash
# Adjust settings via sysctl. Need to be verified which are required and which are optional. (Possible modifications: conf.all / conf.default / conf.<specific_interface>)
# Enable ip forwarding
net.ipv4.ip_forward=1
net.ipv6.conf.all.forwarding=1

# Define different modes for sending replies in response to ARP requests that resolve local target IP addr. 0 -> reply for any local target IP address, configured on any int. Values up to 8, where 8 -> do not reply for all local addresses.
net.ipv4.conf.all.arp_ignore=1

# Define restriction levels for announcing the local source IP address from IP packets in ARP requests. 0 - use any local address, configured on any interface; 1 -> Try to avoid local addresses that are not in the target's subnet for this interface. 2 -> always use the best local address for this target
net.ipv4.conf.all.arp_announce=2

# rp_filter description: 0 -> no source validation; 1 -> strict mode, each incoming packet is tested against the FIB and if the interface is not the best reverse path the packet check will fail. By default failed packets are dicarded; 2 -> ...
net.ipv4.conf.all.rp_filter=0
net.ipv4.conf.ens3.rp_filter=0
```

```bash
# And then apply changes (Please Note that there is no easy way to revert changes made via sysctl command. You can make a copy of current sysctl configuration and then manually set values to default values later.)
sudo sysctl -p
```

```bash
# Configure IP tables - minimal working configuration (Please Note that if you don't force iptables persistent, then this changes will be aborted after VM reboot)
# Note that ens3 is a master interface for macvlan interfaces (n2, n3, n4, n6)
sudo iptables -t nat -A POSTROUTING -o ens3 -j MASQUERADE
sudo iptables -I FORWARD 1 -j ACCEPT
```

```bash
# Make sure that firewall doesn't affect your network traffic
sudo systemctl stop ufw
sudo systemctl disable ufw
```

```bash
# Disable port security. (Please Note that ALLOW ALL IN/OUT traffic via ANY PROTOCOL don't necessary work. Make sure that port security (for now) is completely disabled.
# On OpenStack it can be performed via openstack CLI client. First off all install openstackclient & export appriopriate env variables:
sudo apt install python3-openstackclient

export OS_USERNAME=<username>
export OS_PASSWORD=<password>
export OS_PROJECT_NAME=<project_name>
export OS_USER_DOMAIN_NAME=Default
export OS_AUTH_URL=<auth_url>
export OS_IDENTITY_API_VERSION=3
```

```bash
# Disable security groups in the entire network (interfaces which are already created won't be affected)
openstack network set --disable-port-security <network-name>
```

```bash
# Delete (if assigned) security groups from port - required to disable port security later.
openstack port set --no-security-group <port-name>

# Disable port security
openstack port set --disable-port-security <port-name>

``` 

If VM Nodes are ready, you can create Kubernetes cluster, for exeple using kubeadm. When the cluster is READY there are additional actions which need to be performed.

```bash
# Deploy Multus CNI
git clone https://github.com/k8snetworkplumbingwg/multus-cni.git
cat ./multus-cni/deployments/multus-daemonset.yml | kubectl apply -f -
```

```bash
# On the worker Node create directory - it's required for mongo to start
/bitnami/mongodb/data/db

# And then create PersistentVolume specifying path created aboe and worker node name in `nodeSelectorTerms` field
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv-towards5gs-2
  labels:
    project: free5gc
spec:
  capacity:
    storage: 8Gi
  accessModes:
  - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  local:
    path: /bitnami/mongodb/data/db # Create this path (on worker) before executing this code
  nodeAffinity:
    required:
      nodeSelectorTerms:
      - matchExpressions:
        - key: kubernetes.io/hostname
          operator: In
          values:
          - free5gc-worker-2 # adjust worker name
EOF
```

If you performed all above actions you can move to the actual demo`

## Troubleshooting

### ISSUE: Can't delete parent without deleting child references first

To solve this problem you can either debug manually and search for resource which needs to be removed or you can wipe the database. To wipe the mongo database do as follows:

```bash
# On the cluster where EMCO is deployed exec to mongoDB POD

kubectl -n emco exec -it emco-emco-mongo-0 -- bash

# Enter mongo

mongo

# Then delete resources from appropriate collections

show dbs
use emco    # (or `mco` in old emco releases)
show collections
db.<collection_name>.remove({})

# Then exit mongo and pod bash.
```

### ISSUE: gtp5g device named upfgtp created fail

If you see below logs in the UPF pod you have to reinstall `gtp5g module`

```bash
2021-11-18T08:53:47Z [ERRO][UPF][Util] gtp5g device named upfgtp created fail
2021-11-18T08:53:47Z [ERRO][UPF][Util] Gtp5gDeviceAdd failed
```

```bash
# Reinstall gtp5g module
git clone -b v0.4.0 https://github.com/free5gc/gtp5g.git
cd gtp5g
make clean & make
sudo make install
```

### ISSUE: Pods in Init:0/1 state

MongoDB is a sub-chart of NRF network function. All other (control plane) functions waiting untill NRF is created, and NRF waits until MongoDB is created.
To resolve that issue you have to `Rollout` NRF and optionally WEBUI deployments.

```bash
kubectl rollout restart deployment r1-free5gc-nrf-nrf
kubectl rollout restart deployment r1-free5gc-webui-webui
```
