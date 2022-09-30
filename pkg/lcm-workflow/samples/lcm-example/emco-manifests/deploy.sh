#!/bin/bash

emcoctl --config emco-cfg.yaml apply -f 02.define-app-dig.yaml  -v values.yaml
emcoctl --config emco-cfg.yaml apply -f 08-tac-intent.yaml -v values.yaml
emcoctl --config emco-cfg.yaml apply -f 03.instantiate-dig.yaml -v values.yaml
