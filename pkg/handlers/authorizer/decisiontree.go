package authorizer

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	authz "github.com/aserto-dev/go-grpc-authz/aserto/authorizer/authorizer/v1"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
)

type DecisionTreeCmd struct {
	AuthParams `embed:""`
	Path       string   `name:"path" help:"policy package to evaluate"`
	Decisions  []string `name:"decisions" required:"" help:"policy decisions to return"`
}

func (cmd *DecisionTreeCmd) Run(c *cc.CommonCtx) error {
	client, err := c.AuthorizerClient()
	if err != nil {
		return err
	}

	resource, err := cmd.ResourceContext()
	if err != nil {
		return err
	}

	resp, err := client.Authorizer.DecisionTree(c.Context, &authz.DecisionTreeRequest{
		PolicyContext: &api.PolicyContext{
			Id:        cmd.PolicyID,
			Path:      cmd.Path,
			Decisions: cmd.Decisions,
		},
		IdentityContext: cmd.IdentityContext(),
		ResourceContext: resource,
		Options: &authz.DecisionTreeOptions{
			PathSeparator: authz.PathSeparator_PATH_SEPARATOR_DOT,
		},
	})
	if err != nil {
		return err
	}

	return jsonx.OutputJSONPB(c.UI.Output(), resp)
}
