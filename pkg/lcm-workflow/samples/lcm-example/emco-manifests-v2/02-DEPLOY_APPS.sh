#!/bin/bash

emcoctl --config emco-cfg.yaml apply -f 3a-deploy-edge-app.yaml -v values.yaml
emcoctl --config emco-cfg.yaml apply -f 0a-tac-intent.yaml -v values.yaml
emcoctl --config emco-cfg.yaml apply -f 4a-instantiate-edge-app.yaml -v values.yaml
