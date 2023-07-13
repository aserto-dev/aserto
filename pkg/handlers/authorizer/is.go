package authorizer

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/go-authorizer/aserto/authorizer/v2"
	"github.com/aserto-dev/go-authorizer/aserto/authorizer/v2/api"
)

type EvalDecisionCmd struct {
	AuthParams `embed:""`
	Path       string   `name:"path" required:"" help:"policy package to evaluate"`
	Decisions  []string `name:"decisions" required:"" help:"policy decisions to return"`
}

func (cmd *EvalDecisionCmd) Run(c *cc.CommonCtx) error {
	return errors.Errorf("NOT IMPLEMENTED")

	// client, err := c.AuthorizerClient()
	// if err != nil {
	// 	return err
	// }

<<<<<<< HEAD
	resp, err := client.Authorizer.Is(c.Context, &authorizer.IsRequest{
		PolicyContext: &api.PolicyContext{
			Path:      cmd.Path,
			Decisions: cmd.Decisions,
		},
		IdentityContext: cmd.IdentityContext(),
		ResourceContext: resource,
	})
	if err != nil {
		return err
	}
=======
	// resource, err := cmd.ResourceContext()
	// if err != nil {
	// 	return err
	// }
>>>>>>> 95d083a (compile & login)

	// resp, err := client.Authorizer.Is(c.Context, &authz.IsRequest{
	// 	PolicyContext: &api.PolicyContext{
	// 		Id:        cmd.PolicyID,
	// 		Path:      cmd.Path,
	// 		Decisions: cmd.Decisions,
	// 	},
	// 	IdentityContext: cmd.IdentityContext(),
	// 	ResourceContext: resource,
	// })
	// if err != nil {
	// 	return err
	// }

	// return jsonx.OutputJSONPB(c.UI.Output(), resp)
}
