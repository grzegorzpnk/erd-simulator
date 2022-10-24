#!/bin/bash

emcoctl --config emco-cfg.yaml delete -f 4a-instantiate-edge-app.yaml -v values.yaml

sleep 6

emcoctl --config emco-cfg.yaml delete -f 0a-tac-intent.yaml -v values.yaml
emcoctl --config emco-cfg.yaml delete -f 3a-deploy-edge-app.yaml -v values.yaml

sleep 2

emcoctl --config emco-cfg.yaml delete -f 2-instantiate-logical-cloud.yaml -v values.yaml

sleep 2

emcoctl --config emco-cfg.yaml delete -f 1-prerequisites.yaml -v values.yaml
