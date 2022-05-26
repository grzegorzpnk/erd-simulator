#!/bin/bash
emcoctl --config emco-cfg.yaml delete -f 5-instantiate-ueransim.yaml -v values.yaml

emcoctl --config emco-cfg.yaml delete -f 4-instantiate-free5gc.yaml -v values.yaml

sleep 6

emcoctl --config emco-cfg.yaml delete -f 3-deployment.yaml -v values.yaml

sleep 2

emcoctl --config emco-cfg.yaml delete -f 2-instantiate-logical-cloud.yaml -v values.yaml

sleep 2

emcoctl --config emco-cfg.yaml delete -f 1-prerequisites.yaml -v values.yaml

