package user

import (
	"fmt"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/pkg/errors"
)

type GetCmd struct {
	AccessToken         bool `xor:"group" help:"access token" group:"properties"`
	TenantID            bool `xor:"group" help:"tenant ID" group:"properties"`
	AuthorizerAPIKey    bool `xor:"group" help:"authorizer API key" group:"properties"`
	RegistryDownloadKey bool `xor:"group" help:"registry download key" group:"properties"`
	RegistryUploadKey   bool `xor:"group" help:"registry upload key" group:"properties"`
	DecisionLogsKey     bool `xor:"group" help:"decision logs key" group:"properties"`
	Token               bool `xor:"group" help:"token" hidden:"" group:"properties"`
}

func (cmd *GetCmd) Run(c *cc.CommonCtx) error {
	if !cmd.AccessToken && !cmd.TenantID && !cmd.AuthorizerAPIKey && !cmd.RegistryDownloadKey && !cmd.RegistryUploadKey && !cmd.Token && !cmd.DecisionLogsKey {
		return errors.Errorf("no property flag provided")
	}

	var propValue string
	switch {
	case cmd.AccessToken:
		propValue = c.AccessToken()
	case cmd.TenantID:
		propValue = c.TenantID()
	case cmd.AuthorizerAPIKey:
		propValue = c.AuthorizerAPIKey()
	case cmd.RegistryDownloadKey:
		propValue = c.RegistryDownloadKey()
	case cmd.RegistryUploadKey:
		propValue = c.RegistryUploadKey()
	case cmd.DecisionLogsKey:
		propValue = c.DecisionLogsKey()
	case cmd.Token:
		return jsonx.OutputJSON(c.OutWriter, c.Token())
	}

	fmt.Fprintf(c.OutWriter, "%s\n", propValue)

	return nil
}
