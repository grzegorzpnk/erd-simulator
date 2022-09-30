#!/bin/bash

emcoctl --config emco-cfg.yaml apply -f 00.define-clusters-proj.yaml -v values.yaml
emcoctl --config emco-cfg.yaml apply -f 01.instantiate-lc.yaml  -v values.yaml
