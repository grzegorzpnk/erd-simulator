#!/bin/bash

emcoctl --config emco-cfg.yaml delete -f 03.instantiate-dig.yaml -v values.yaml

sleep 4

emcoctl --config emco-cfg.yaml delete -f 08-tac-intent.yaml -v values.yaml
emcoctl --config emco-cfg.yaml delete -f 02.define-app-dig.yaml  -v values.yaml

sleep 2

emcoctl --config emco-cfg.yaml delete -f 01.instantiate-lc.yaml  -v values.yaml

sleep 4

emcoctl --config emco-cfg.yaml delete -f 00.define-clusters-proj.yaml -v values.yaml
