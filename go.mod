module github.com/aserto-dev/aserto

go 1.16

// replace github.com/aserto-dev/go-grpc => ../go-grpc
// replace github.com/aserto-dev/go-grpc-authz => ../go-grpc-authz
// replace github.com/aserto-dev/go-lib => ../go-lib
// replace github.com/aserto-dev/aserto-tenant => ../aserto-tenant
// replace github.com/aserto-dev/aserto-go => ../aserto-go

require (
	github.com/alecthomas/kong v0.2.17
	github.com/aserto-dev/aserto-go v0.0.5
	github.com/aserto-dev/aserto-tenant v0.1.2
	github.com/aserto-dev/go-grpc v0.0.2
	github.com/aserto-dev/go-grpc-authz v0.0.2
	github.com/aserto-dev/go-lib v0.3.0
	github.com/cli/browser v1.1.0
	github.com/containerd/containerd v1.5.4
	github.com/fatih/color v1.12.0
	github.com/google/uuid v1.3.0
	github.com/joho/godotenv v1.3.0
	github.com/magefile/mage v1.11.0
	github.com/pkg/errors v0.9.1
	github.com/zalando/go-keyring v0.1.1
	golang.org/x/sys v0.0.0-20210819135213-f52c844e1c1c // indirect
	google.golang.org/genproto v0.0.0-20210819152826-bb2f42df81c9
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/auth0.v5 v5.19.2
	oras.land/oras-go v0.4.0
	rsc.io/letsencrypt v0.0.3 // indirect
)
