#!/bin/bash

emcoctl --config emco-cfg.yaml delete -f 4b-instantiate-ueransim.yaml -v values.yaml
emcoctl --config emco-cfg.yaml delete -f 4a-instantiate-free5gc.yaml -v values.yaml

sleep 6

emcoctl --config emco-cfg.yaml delete -f 3a-deploy-free5gc.yaml -v values.yaml
emcoctl --config emco-cfg.yaml delete -f 3b-deploy-ueransim.yaml -v values.yaml

sleep 2

emcoctl --config emco-cfg.yaml delete -f 2-instantiate-logical-cloud.yaml -v values.yaml

sleep 2

emcoctl --config emco-cfg.yaml delete -f 1-prerequisites.yaml -v values.yaml
