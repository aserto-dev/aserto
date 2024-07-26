package cmd_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/aserto-dev/aserto/pkg/cc/config"
	"github.com/aserto-dev/aserto/pkg/cmd"
	"github.com/stretchr/testify/require"
)

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
