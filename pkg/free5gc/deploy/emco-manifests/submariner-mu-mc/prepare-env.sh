#!/bin/bash

emcoctl --config emco-cfg.yaml apply -f 1-prerequisites.yaml -v values.yaml
emcoctl --config emco-cfg.yaml apply -f 2-instantiate-logical-cloud.yaml -v values.yaml
