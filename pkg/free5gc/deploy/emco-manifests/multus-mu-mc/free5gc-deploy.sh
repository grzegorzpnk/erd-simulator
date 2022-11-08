#!/bin/bash

emcoctl --config emco-cfg.yaml apply -f 3a-deploy-free5gc.yaml -v values.yaml
emcoctl --config emco-cfg.yaml apply -f 4a-instantiate-free5gc.yaml -v values.yaml
