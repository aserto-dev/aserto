package authorizer

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/grpcc"
	"github.com/aserto-dev/aserto/pkg/grpcc/authorizer"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	authz "github.com/aserto-dev/go-grpc-authz/aserto/authorizer/authorizer/v1"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
)

type DecisionTreeCmd struct {
	PolicyID  string
	Path      string
	Decisions []string
}

func (cmd *DecisionTreeCmd) Run(c *cc.CommonCtx) error {
	conn, err := authorizer.Connection(
		c.Context,
		c.AuthorizerService(),
		grpcc.NewAPIKeyAuth(c.AuthorizerAPIKey()),
	)
	if err != nil {
		return err
	}

	ctx := grpcc.SetTenantContext(c.Context, c.TenantID())

	authzClient := conn.AuthorizerClient()
	resp, err := authzClient.DecisionTree(ctx, &authz.DecisionTreeRequest{
		PolicyContext: &api.PolicyContext{
			Id:        cmd.PolicyID,
			Path:      cmd.Path,
			Decisions: cmd.Decisions,
		},
		IdentityContext: &api.IdentityContext{
			Identity: "",
			Type:     api.IdentityType_IDENTITY_TYPE_NONE,
		},
		Options: &authz.DecisionTreeOptions{
			PathSeparator: authz.PathSeparator_PATH_SEPARATOR_DOT,
		},
	})

	if err != nil {
		return err
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}
