emcoctl --config emco-cfg.yaml -v values.yaml delete -f 04.*
emcoctl --config emco-cfg.yaml -v values.yaml delete -f 08.*

emcoctl --config emco-cfg.yaml -v values.yaml delete -f 03.* -w 5

emcoctl --config emco-cfg.yaml -v values.yaml delete -f 02.*

emcoctl --config emco-cfg.yaml -v values.yaml delete -f 01.* -w 3

emcoctl --config emco-cfg.yaml -v values.yaml delete -f 00.*

