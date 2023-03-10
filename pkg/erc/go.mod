module 10.254.188.33/matyspi5/erd/pkg/erc

require (
	github.com/antonfisher/nested-logrus-formatter v1.3.1
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	gitlab.com/project-emco/core/emco-base/src/orchestrator v0.0.0-00010101000000-000000000000
)

replace (
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
	gitlab.com/project-emco/core/emco-base/src/clm => ./../../../emco-base/src/clm //Please modify accordingly.
	gitlab.com/project-emco/core/emco-base/src/monitor => ./../../../emco-base/src/monitor //Please modify accordingly.
	gitlab.com/project-emco/core/emco-base/src/orchestrator => ./../../../emco-base/src/orchestrator //Please modify accordingly.
	gitlab.com/project-emco/core/emco-base/src/rsync => ./../../../emco-base/src/rsync //Please modify accordingly.
	go.etcd.io/etcd => go.etcd.io/etcd v0.5.0-alpha.5.0.20200819165624-17cef6e3e9d5 // 17cef6e3e9d5 is the SHA for git tag v3.4.12
	google.golang.org/grpc => google.golang.org/grpc v1.28.0
	helm.sh/helm/v3 => helm.sh/helm/v3 v3.5.3
)

go 1.16
