package authorizer

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/grpcc"
	"github.com/aserto-dev/aserto/pkg/grpcc/authorizer"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	authz "github.com/aserto-dev/go-grpc-authz/aserto/authorizer/authorizer/v1"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
)

type ExecQueryCmd struct {
	Statement string `arg:"stmt" name:"stmt" required:"" help:"query statement"`
	Input     string `name:"input" optional:"" help:"query input context"`
}

func (cmd *ExecQueryCmd) Run(c *cc.CommonCtx) error {

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
	resp, err := authzClient.Query(ctx, &authz.QueryRequest{
		Query: cmd.Statement,
		Input: cmd.Input,
		IdentityContext: &api.IdentityContext{
			Identity: "",
			Type:     api.IdentityType_IDENTITY_TYPE_NONE,
		},
		Options: &authz.QueryOptions{
			Metrics:      false,
			Instrument:   false,
			Trace:        authz.TraceLevel_TRACE_LEVEL_OFF,
			TraceSummary: false,
		},
	})

	if err != nil {
		return err
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}
