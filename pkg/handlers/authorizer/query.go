package authorizer

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/grpcc"
	"github.com/aserto-dev/aserto/pkg/grpcc/authorizer"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/proto/aserto/api"
	authz "github.com/aserto-dev/proto/aserto/authorizer/authorizer"
)

type ExecQueryCmd struct {
	Statement string `arg:"stmt" name:"stmt" required:"" help:"query statement"`
	Input     string `name:"input" optional:"" help:"query input context"`
}

func (cmd *ExecQueryCmd) Run(c *cc.CommonCtx) error {

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
	resp, err := authzClient.Query(ctx, &authz.QueryRequest{
		Query: cmd.Statement,
		Input: cmd.Input,
		IdentityContext: &api.IdentityContext{
			Mode:     api.IdentityMode_ANONYMOUS,
			Identity: "",
		},
		Options: &authz.QueryOptions{
			Metrics:      false,
			Instrument:   false,
			Trace:        authz.TraceLevel_OFF,
			TraceSummary: false,
		},
	})

	if err != nil {
		return err
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}
