#!/bin/bash

emcoctl --config emco-cfg.yaml apply -f 3b-deploy-ueransim.yaml -v values.yaml
emcoctl --config emco-cfg.yaml apply -f 4b-instantiate-ueransim.yaml -v values.yaml
