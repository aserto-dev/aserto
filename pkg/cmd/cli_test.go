package cmd_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/cc/clients"
	"github.com/aserto-dev/aserto/pkg/cc/iostream"
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
		clients.NewServiceOptions().ConfigOverrider,
	)
	assert.NoError(err)
	assert.NoError(kongCtx.Run(c))

	assert.Equal(
		fmt.Sprintf("%s - %s (%s)\n", x.AppName, version.GetInfo().String(), x.AppVersionTag),
		ios.Out.String(),
	)
}
