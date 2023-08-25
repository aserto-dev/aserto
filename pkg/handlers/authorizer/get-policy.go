package authorizer

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/go-authorizer/aserto/authorizer/v2"
	"github.com/aserto-dev/go-authorizer/aserto/authorizer/v2/api"
)

type GetPolicyCmd struct {
	PolicyID      string `name:"policy-id" required:"" help:"policy id"`
	PolicyName    string `name:"policy-name" required:"" help:"policy name"`
	InstanceLabel string `name:"instance-label" required:"" help:"policy's instance label"`
}

func (cmd *GetPolicyCmd) Run(c *cc.CommonCtx) error {
	client, err := c.AuthorizerClient()
	if err != nil {
		return err
	}

	resp, err := client.Authorizer.GetPolicy(c.Context, &authorizer.GetPolicyRequest{
		Id:             cmd.PolicyID,
		PolicyInstance: &api.PolicyInstance{Name: cmd.PolicyName, InstanceLabel: cmd.InstanceLabel},
	})
	if err != nil {
		return err
	}
	return jsonx.OutputJSONPB(c.UI.Output(), resp)
}
