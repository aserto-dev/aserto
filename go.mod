module github.com/aserto-dev/aserto

go 1.16

// replace github.com/aserto-dev/proto => ../proto

require (
	github.com/99designs/keyring v1.1.6
	github.com/alecthomas/kong v0.2.17
	github.com/aserto-dev/aserto-tenant v0.0.149
	github.com/aserto-dev/go-lib v0.2.64
	github.com/aserto-dev/proto v0.0.44
	github.com/cli/browser v1.1.0
	github.com/containerd/containerd v1.5.2
	github.com/fatih/color v1.12.0
	github.com/google/uuid v1.3.0
	github.com/joho/godotenv v1.3.0
	github.com/magefile/mage v1.11.0
	github.com/pkg/errors v0.9.1
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
	google.golang.org/genproto v0.0.0-20210713002101-d411969a0d9a
	google.golang.org/grpc v1.39.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/auth0.v5 v5.19.1
	oras.land/oras-go v0.4.0
	rsc.io/letsencrypt v0.0.3 // indirect
)
