Intro
---

Submariner allows to use service mesh based communication between free5gc microservices (if needed).
* Communication between SMF <-> UPFB/UPF1/UPF2 (N4)
* Communication between gNB <-> UPFB (N3)
* Communication between UPFs (N9)
* ...

Submariner installation
---

1. Install subctl on the management cluster

```
curl -Ls https://get.submariner.io | bash
export PATH=$PATH:~/.local/bin
echo export PATH=\$PATH:~/.local/bin >> ~/.profile
```

2. Deploy broker into management cluster 

```
subctl deploy-broker --kubeconfig <PATH-TO-KUBECONFIG-BROKER>
```

3. Join clusters to the Service Mesh (mgnt, up, cp clusters)

```
subctl join --kubeconfig <PATH-TO-JOINING-CLUSTER> broker-info.subm --clusterid <ID> --cable-driver <CABLE-DRIVER>
```

`Note: libreswan is used by default (working but not all the time), the best results were achieved using: wireguard cable-drive. It needs further testing.`

Export services
---

```
subctl --kubeconfig <PATH-TO-KUBECONFIG-WHERE-SVC-IS-PRESENT> export service --namespace <NS-NAME> <SERVICE-NAME>
```

Use exported service
---

Service is avaliabla as `<svc-name>.<namespace>.svc.clusterset.local`


---

Uninstall Submariner
---

Perform [these steps](https://submariner.io/operations/cleanup/) to uninstall Submariner

---

For more information check out: [User Guide](https://submariner.io/operations/usage/), [Deployment](https://submariner.io/operations/deployment/)


Example commands
---

### Install broker

```
subctl deploy-broker --kubeconfig ../../cfg/meh1.config
```

### Join clusters

```
subctl join --kubeconfig ../../cfg/meh1.config broker-info.subm --clusterid 148 --cable-driver wireguard && \
subctl join --kubeconfig ../../cfg/meh3.config broker-info.subm --clusterid 127 --cable-driver wireguard && \
subctl join --kubeconfig ../../cfg/ran.config broker-info.subm --clusterid 150 --cable-driver wireguard && \
subctl join --kubeconfig ../../cfg/core.config broker-info.subm --clusterid 144 --cable-driver wireguard


# subctl join --kubeconfig ../../cfg/mgnt.config broker-info.subm --clusterid 170 --cable-driver wireguard && \
# subctl join --kubeconfig ../../cfg/meh2.config broker-info.subm --clusterid 117 --cable-driver wireguard && \

```

### Export

```
subctl --kubeconfig ../../cfg/meh1.config export service upf1-gtpu
subctl --kubeconfig ../../cfg/meh3.config export service upf2-gtpu
subctl --kubeconfig ../../cfg/ran.config export service upfb-gtpu

subctl --kubeconfig ../../cfg/meh1.config export service upf1-pfcp
subctl --kubeconfig ../../cfg/meh3.config export service upf2-pfcp
subctl --kubeconfig ../../cfg/ran.config export service upfb-pfcp
```

### Verify

```
kubectl get serviceexports.multicluster.x-k8s.io -A
kubectl get serviceimports.multicluster.x-k8s.io -A
```

### Delete

```
kubectl delete serviceexports.multicluster.x-k8s.io --all
```

### Clear-Up

```
kubectl delete namespace submariner-operator submariner-k8s-broker
for CRD in `kubectl get crds | grep -iE 'submariner|multicluster.x-k8s.io'| awk '{print $1}'`; do kubectl delete crd $CRD; done
roles="submariner-operator submariner-operator-globalnet submariner-lighthouse submariner-networkplugin-syncer"
kubectl delete clusterrole,clusterrolebinding $roles --ignore-not-found
kubectl label --all node submariner.io/gateway-

sudo iptables --flush SUBMARINER-INPUT
sudo iptables -D INPUT $(sudo iptables -L INPUT --line-numbers | grep SUBMARINER-INPUT | awk '{print $1}')
sudo iptables --delete-chain SUBMARINER-INPUT

sudo iptables -t nat --flush SUBMARINER-POSTROUTING
sudo iptables -t nat -D POSTROUTING $(sudo iptables -t nat -L POSTROUTING --line-numbers | grep SUBMARINER-POSTROUTING | awk '{print $1}')
sudo iptables -t nat --delete-chain SUBMARINER-POSTROUTING

sudo iptables -t mangle --flush SUBMARINER-POSTROUTING
sudo iptables -t mangle -D POSTROUTING $(sudo iptables -t mangle -L POSTROUTING --line-numbers | grep SUBMARINER-POSTROUTING | awk '{print $1}')
sudo iptables -t mangle --delete-chain SUBMARINER-POSTROUTING

sudo ipset destroy SUBMARINER-LOCALCIDRS
sudo ipset destroy SUBMARINER-REMOTECIDRS

sudo ip link delete vx-submariner
```
