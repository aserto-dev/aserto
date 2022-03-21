package cmd_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto/pkg/auth0/api"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/cc/clients"
	"github.com/aserto-dev/aserto/pkg/cc/config"
	"github.com/aserto-dev/aserto/pkg/cc/iostream"
	"github.com/aserto-dev/aserto/pkg/cc/token"
	"github.com/aserto-dev/aserto/pkg/cmd"
	"github.com/aserto-dev/aserto/pkg/version"
	"github.com/aserto-dev/aserto/pkg/x"
	"github.com/stretchr/testify/require"
)

func TestVersionCmd(t *testing.T) {
	assert := require.New(t)
	cli := cmd.CLI{}
	parser, err := kong.New(&cli)
	assert.NoError(err)

	kongCtx, err := parser.Parse([]string{"version"})
	assert.NoError(err)

	ios := iostream.BytesIO()
	c, err := cc.BuildTestCtx(
		ios,
		bytes.NewReader([]byte{}),
		cli.ConfigOverrider,
		clients.NewServiceOptions(),
	)
	assert.NoError(err)
	assert.NoError(kongCtx.Run(c))

	assert.Equal(
		fmt.Sprintf("%s - %s (%s)\n", x.AppName, version.GetInfo().String(), x.AppVersionTag),
		ios.Out.String(),
	)
}

func TestTenantID(t *testing.T) {
	tests := []struct {
		Name     string
		TenantID string
		Token    *api.Token
		Expected string
	}{
		{"override with no cached token", "testID", nil, "testID"},
		{"override with cached token", "testID", newToken("cached", false), "testID"},
		{"cached token with no override", "", newToken("cached", false), "cached"},
		{"expired token and no override", "", newToken("cached", true), ""},
		{"no token and no override", "", nil, ""},
	}

	for _, test := range tests {
		t.Run(test.Name, func(tt *testing.T) {
			assert := require.New(tt)
			cli := &cmd.CLI{TenantOverride: test.TenantID}
			assert.Equal(test.Expected, cli.TenantID(token.New(test.Token)))
		})
	}
}

func TestConfigOverrider(t *testing.T) {
	tests := []struct {
		Name           string
		Config         io.Reader
		TenantOverride string
		Expected       string
	}{
		{"no config, no override", bytes.NewBufferString(""), "", ""},
		{"no config, with override", bytes.NewBufferString(""), "overrideID", "overrideID"},
		{"with config, no override", bytes.NewBufferString("tenant_id: configID"), "", "configID"},
		{"with config, and override", bytes.NewBufferString("tenant_id: configID"), "overrideID", "overrideID"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(tt *testing.T) {
			assert := require.New(tt)
			cli := &cmd.CLI{TenantOverride: test.TenantOverride}
			conf, err := config.NewTestConfig(test.Config, cli.ConfigOverrider)
			assert.NoError(err)
			assert.Equal(test.Expected, conf.TenantID)
		})
	}
}

func newToken(tenantID string, expired bool) *api.Token {
	expiresAt := time.Now().UTC()
	offset := 24 * time.Hour
	if expired {
		expiresAt = expiresAt.Add(-offset)
	} else {
		expiresAt = expiresAt.Add(offset)
	}

	return &api.Token{TenantID: tenantID, ExpiresAt: expiresAt}
}
