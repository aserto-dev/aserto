module github.com/aserto-dev/aserto

go 1.21

toolchain go1.22.1

// replace github.com/aserto-dev/go-grpc => ../go-grpc
// replace github.com/aserto-dev/go-grpc-authz => ../go-grpc-authz
// replace github.com/aserto-dev/aserto-go => ../aserto-go

require (
	github.com/alecthomas/kong v0.8.1
	github.com/aserto-dev/certs v0.0.3
	github.com/aserto-dev/clui v0.8.3
	github.com/aserto-dev/go-aserto v0.31.3
	github.com/aserto-dev/go-authorizer v0.20.5
	github.com/aserto-dev/go-decision-logs v0.0.4
	github.com/aserto-dev/go-grpc v0.8.59
	github.com/aserto-dev/logger v0.0.4
	github.com/cli/browser v1.3.0
	github.com/fatih/color v1.16.0
	github.com/getkin/kin-openapi v0.115.0
	github.com/google/wire v0.5.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/joho/godotenv v1.5.1
	github.com/magefile/mage v1.15.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.32.0
	github.com/spf13/viper v1.18.1
	github.com/stretchr/testify v1.9.0
	github.com/zalando/go-keyring v0.2.3
	github.com/zenizh/go-capturer v0.0.0-20211219060012-52ea6c8fed04
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240401170217-c3f982113cda
	google.golang.org/grpc v1.63.2
	google.golang.org/protobuf v1.33.0
	gopkg.in/auth0.v5 v5.21.1
)

require (
	github.com/PuerkitoBio/rehttp v1.0.0 // indirect
	github.com/alessio/shellescape v1.4.1 // indirect
	github.com/aserto-dev/header v0.0.7 // indirect
	github.com/benbjohnson/clock v1.1.0 // indirect
	github.com/danieljoos/wincred v1.2.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-http-utils/headers v0.0.0-20181008091004-fed159eddc2a // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/subcommands v1.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/invopop/yaml v0.2.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/kyokomi/emoji v2.2.4+incompatible // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pelletier/go-toml/v2 v2.1.0 // indirect
	github.com/perimeterx/marshmallow v1.1.5 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/exp v0.0.0-20240325151524-a685a6edb6d8 // indirect
	golang.org/x/mod v0.16.0 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/oauth2 v0.17.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/term v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.19.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240401170217-c3f982113cda // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
