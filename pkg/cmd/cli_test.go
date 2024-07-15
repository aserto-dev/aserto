package cmd_test

// import (
// 	"bytes"
// 	"fmt"
// 	"io"
// 	"strconv"
// 	"testing"

// 	"github.com/alecthomas/kong"
// 	"github.com/aserto-dev/aserto/pkg/cc"
// 	"github.com/aserto-dev/aserto/pkg/cc/clients"
// 	"github.com/aserto-dev/aserto/pkg/cc/config"
// 	"github.com/aserto-dev/aserto/pkg/cc/iostream"
// 	"github.com/aserto-dev/aserto/pkg/cmd"
// 	"github.com/aserto-dev/aserto/pkg/version"
// 	"github.com/aserto-dev/aserto/pkg/x"
// 	"github.com/stretchr/testify/require"

// 	topazCC "github.com/aserto-dev/topaz/pkg/cli/cc"
// )

// func TestVersionCmd(t *testing.T) {
// 	assert := require.New(t)
// 	cli := cmd.CLI{}
// 	vars := kong.Vars{
// 		"topaz_dir":          topazCC.GetTopazDir(),
// 		"topaz_certs_dir":    topazCC.GetTopazCertsDir(),
// 		"topaz_cfg_dir":      topazCC.GetTopazCfgDir(),
// 		"topaz_db_dir":       topazCC.GetTopazDataDir(),
// 		"container_registry": topazCC.ContainerRegistry(),
// 		"container_image":    topazCC.ContainerImage(),
// 		"container_tag":      topazCC.ContainerTag(),
// 		"container_platform": topazCC.ContainerPlatform(),
// 		"container_name":     topazCC.ContainerName(""),
// 		"directory_svc":      topazCC.DirectorySvc(),
// 		"directory_key":      topazCC.DirectoryKey(),
// 		"directory_token":    topazCC.DirectoryToken(),
// 		"authorizer_svc":     topazCC.AuthorizerSvc(),
// 		"authorizer_key":     topazCC.AuthorizerKey(),
// 		"authorizer_token":   topazCC.AuthorizerToken(),
// 		"tenant_id":          topazCC.TenantID(),
// 		"insecure":           strconv.FormatBool(topazCC.Insecure()),
// 		"no_check":           strconv.FormatBool(topazCC.NoCheck()),
// 	}
// 	parser, err := kong.New(&cli, vars)
// 	assert.NoError(err)

// 	kongCtx, err := parser.Parse([]string{"version"})
// 	assert.NoError(err)

// 	ios := iostream.BytesIO()
// 	c, err := cc.BuildTestCtx(
// 		ios,
// 		bytes.NewReader([]byte{}),
// 		cli.ConfigOverrider,
// 		clients.NewServiceOptions().ConfigOverrider,
// 	)
// 	assert.NoError(err)
// 	assert.NoError(kongCtx.Run(c))

// 	assert.Equal(
// 		fmt.Sprintf("%s - %s (%s)\n", x.AppName, version.GetInfo().String(), x.AppVersionTag),
// 		ios.Out.String(),
// 	)
// }

// func TestConfigOverrider(t *testing.T) {
// 	tests := []struct {
// 		Name           string
// 		Config         io.Reader
// 		TenantOverride string
// 		Expected       string
// 	}{
// 		{"no config, no override", bytes.NewBufferString(""), "", ""},
// 		{"no config, with override", bytes.NewBufferString(""), "overrideID", "overrideID"},
// 		{"with config, no override", bytes.NewBufferString("tenant_id: configID"), "", "configID"},
// 		{"with config, and override", bytes.NewBufferString("tenant_id: configID"), "overrideID", "overrideID"},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.Name, func(tt *testing.T) {
// 			assert := require.New(tt)
// 			cli := &cmd.CLI{TenantOverride: test.TenantOverride}
// 			conf, err := config.NewTestConfig(test.Config, cli.ConfigOverrider)
// 			assert.NoError(err)
// 			assert.Equal(test.Expected, conf.TenantID)
// 		})
// 	}
// }
