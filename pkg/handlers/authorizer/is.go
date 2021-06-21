package authorizer

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/grpcc"
	"github.com/aserto-dev/aserto/pkg/grpcc/authorizer"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/proto/aserto/api"
	authz "github.com/aserto-dev/proto/aserto/authorizer/authorizer"
)

type EvalDecisionCmd struct {
	PolicyID  string
	Path      string
	Decisions []string
}

func (cmd *EvalDecisionCmd) Run(c *cc.CommonCtx) error {
	conn, err := authorizer.Connection(
		c.Context,
		c.AuthorizerService(),
		grpcc.NewTokenAuth(c.AccessToken()),
	)
	if err != nil {
		return err
	}

	ctx := grpcc.SetTenantContext(c.Context, c.TenantID())
	ctx = grpcc.SetAsertoAPIKey(ctx, c.AuthorizerAPIKey())

	authzClient := conn.AuthorizerClient()
	resp, err := authzClient.Is(ctx, &authz.IsRequest{
		PolicyContext: &api.PolicyContext{
			Id:        cmd.PolicyID,
			Path:      cmd.Path,
			Decisions: cmd.Decisions,
		},
		IdentityContext: &api.IdentityContext{
			Mode:     api.IdentityMode_ANONYMOUS,
			Identity: "",
		},
		ResourceContext: map[string]string{},
	})

	if err != nil {
		return err
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}
