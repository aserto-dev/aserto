package user

import (
	"fmt"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/pkg/errors"
)

type GetCmd struct {
	AccessToken      bool `xor:"group" help:"access token" group:"properties"`
	TenantID         bool `xor:"group" help:"tenant ID" group:"properties"`
	AuthorizerAPIKey bool `xor:"group" help:"authorizer API key" group:"properties"`
	DecisionLogsKey  bool `xor:"group" help:"decision logs key" group:"properties"`
	Token            bool `xor:"group" help:"token" hidden:"" group:"properties"`
}

func (cmd *GetCmd) Run(c *cc.CommonCtx) error {
	if !cmd.AccessToken && !cmd.TenantID && !cmd.AuthorizerAPIKey && !cmd.Token && !cmd.DecisionLogsKey {
		return errors.Errorf("no property flag provided")
	}

	var (
		propValue string
		err       error
	)

	switch {
	case cmd.AccessToken:
		propValue, err = c.AccessToken()
	case cmd.TenantID:
		propValue = c.TenantID()
	case cmd.AuthorizerAPIKey:
		propValue, err = c.AuthorizerAPIKey()
	case cmd.DecisionLogsKey:
		propValue, err = c.DecisionLogsKey()
	case cmd.Token:
		token, tokenErr := c.Token()
		if tokenErr != nil {
			return tokenErr
		}
		return jsonx.OutputJSON(c.UI.Output(), token)
	}

	if err != nil {
		return err
	}

	fmt.Fprintf(c.UI.Output(), "%s\n", propValue)

	return nil
}
