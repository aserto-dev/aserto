module github.com/aserto-dev/aserto

go 1.17

// replace github.com/aserto-dev/go-grpc => ../go-grpc
// replace github.com/aserto-dev/go-grpc-authz => ../go-grpc-authz
// replace github.com/aserto-dev/go-lib => ../go-lib
// replace github.com/aserto-dev/aserto-tenant => ../aserto-tenant
// replace github.com/aserto-dev/aserto-go => ../aserto-go

replace oras.land/oras-go => oras.land/oras-go v0.4.0

require (
	github.com/alecthomas/kong v0.2.22
	github.com/aserto-dev/aserto-go v0.7.0
	github.com/aserto-dev/aserto-tenant v0.7.0
	github.com/aserto-dev/go-grpc v0.7.0
	github.com/aserto-dev/go-grpc-authz v0.7.0
	github.com/aserto-dev/go-lib v0.7.0
	github.com/aserto-dev/go-utils v0.7.1
	github.com/aserto-dev/mage-loot v0.7.0
	github.com/cli/browser v1.1.0
	github.com/containerd/containerd v1.5.9
	github.com/fatih/color v1.13.0
	github.com/google/uuid v1.3.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/joho/godotenv v1.4.0
	github.com/magefile/mage v1.12.1
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.26.1
	github.com/zalando/go-keyring v0.1.1
	google.golang.org/genproto v0.0.0-20220118154757-00ab72f36ad5
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/auth0.v5 v5.21.1
	oras.land/oras-go v0.4.0
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/Microsoft/go-winio v0.5.1 // indirect
	github.com/PuerkitoBio/rehttp v1.0.0 // indirect
	github.com/aserto-dev/aserto-idp-plugin-okta v0.7.0 // indirect
	github.com/aserto-dev/clui v0.7.0 // indirect
	github.com/aserto-dev/go-grpc-clients v0.7.1 // indirect
	github.com/aserto-dev/go-grpc-internal v0.7.0 // indirect
	github.com/aserto-dev/idp-plugin-sdk v0.7.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/danieljoos/wincred v1.1.0 // indirect
	github.com/docker/cli v20.10.11+incompatible // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v20.10.11+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.6.4 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/godbus/dbus/v5 v5.0.4 // indirect
	github.com/gofrs/uuid v4.2.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-github/v33 v33.0.0 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-hclog v0.16.2 // indirect
	github.com/hashicorp/go-plugin v1.4.3 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/kyokomi/emoji v2.2.4+incompatible // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/mitchellh/go-testing-interface v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.4.2 // indirect
	github.com/moby/locker v1.0.1 // indirect
	github.com/moby/term v0.0.0-20200312100748-672ec06f55cd // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/oklog/run v1.0.0 // indirect
	github.com/okta/okta-sdk-golang/v2 v2.9.1 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/open-policy-agent/opa v0.35.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/patrickmn/go-cache v0.0.0-20180815053127-5633e0862627 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/prometheus/client_golang v1.11.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.29.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.9.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/crypto v0.0.0-20220112180741-5e0467b6c7ce // indirect
	golang.org/x/net v0.0.0-20211216030914-fe4d6282115f // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220111092808-5a964db01320 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/ini.v1 v1.63.2 // indirect
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	rsc.io/letsencrypt v0.0.3 // indirect
)
