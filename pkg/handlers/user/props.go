package user

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/topaz/pkg/cli/jsonx"
	"github.com/pkg/errors"
)

type GetCmd struct {
	Property string `kong:"-"`
}

func (cmd *GetCmd) BeforeApply(k *kong.Kong, c *kong.Context) error {
	p := c.Path[len(c.Path)-1]
	cmd.Property = p.Command.Name
	return nil
}

func (cmd *GetCmd) Run(c *cc.CommonCtx) error {
	var (
		propValue string
		err       error
	)

	switch cmd.Property {
	case "access-token":
		propValue, err = c.AccessToken()
	case "tenant-id":
		token, tokenErr := c.Token()
		if tokenErr != nil {
			return tokenErr
		}
		propValue = token.TenantID
	case "authorizer-key":
		propValue, err = c.AuthorizerAPIKey()
	case "directory-read-key":
		propValue, err = c.DirectoryReadKey()
	case "directory-write-key":
		propValue, err = c.DirectoryWriteKey()
	case "discovery-key":
		propValue, err = c.DiscoveryKey()
	case "decision-logs-key":
		propValue, err = c.DecisionLogsKey()
	case "registry-read-key":
		propValue, err = c.RegistryReadKey()
	case "registry-write-key":
		propValue, err = c.RegistryWriteKey()
	case "token":
		token, tokenErr := c.Token()
		if tokenErr != nil {
			return tokenErr
		}
		return jsonx.OutputJSON(c.StdOut(), token)

	default:
		return errors.Errorf("unknown property name %s", cmd.Property)
	}

	if err != nil {
		return err
	}

	fmt.Fprintf(c.StdOut(), "%s\n", propValue)

	return nil
}
