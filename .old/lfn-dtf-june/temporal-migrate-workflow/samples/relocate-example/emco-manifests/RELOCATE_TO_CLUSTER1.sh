emcoctl --config emco-cfg.yaml -v values.yaml delete -f 04.define-workflow-1.yaml

emcoctl --config emco-cfg.yaml -v values.yaml apply -f 08.define-workflow-2.yaml

emcoctl --config emco-cfg.yaml -v values.yaml apply -f 09.start-workflow-2.yaml

