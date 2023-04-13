module 10.254.188.33/matyspi5/erd/pkg/lcm-workflow

go 1.18

require (
	github.com/deepmap/oapi-codegen v1.11.0
	github.com/getkin/kin-openapi v0.94.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/labstack/echo/v4 v4.8.0
	gitlab.com/project-emco/core/emco-base/src/workflowmgr v0.0.0-00010101000000-000000000000
	go.temporal.io/sdk v1.13.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/gogo/status v1.1.0 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.15.0 // indirect
	github.com/labstack/gommon v0.3.1 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/robfig/cron v1.2.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/stretchr/testify v1.7.1 // indirect
	github.com/tidwall/gjson v1.14.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	gitlab.com/project-emco/core/emco-base/src/orchestrator v0.0.0-20220308165843-684035a6c4f0 // indirect
	go.mongodb.org/mongo-driver v1.8.4 // indirect
	go.temporal.io/api v1.7.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	golang.org/x/crypto v0.0.0-20220513210258-46612604a0f9 // indirect
	golang.org/x/net v0.0.0-20220513224357-95641704303c // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220513210249-45d2b4557a2a // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20220411224347-583f2d630306 // indirect
	google.golang.org/genproto v0.0.0-20220308174144-ae0e22291548 // indirect
	google.golang.org/grpc v1.44.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace (
	gitlab.com/project-emco/core/emco-base => ../../../emco-base
	gitlab.com/project-emco/core/emco-base/src/clm => ../../../emco-base/src/clm
	gitlab.com/project-emco/core/emco-base/src/monitor => ../../../emco-base/src/monitor
	gitlab.com/project-emco/core/emco-base/src/orchestrator => ../../../emco-base/src/orchestrator
	gitlab.com/project-emco/core/emco-base/src/rsync => ../../../emco-base/src/rsync
	gitlab.com/project-emco/core/emco-base/src/workflowmgr => ../../../emco-base/src/workflowmgr
	gitlab.com/project-emco/core/emco-base/src/workflowmgr/pkg/emcotemporalapi => ../../../emco-base/src/workflowmgr/pkg/emcotemporalapi
)
