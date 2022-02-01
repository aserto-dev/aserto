package authorizer

import (
	aserto "github.com/aserto-dev/aserto-go/client"
	"github.com/aserto-dev/aserto-go/client/authorizer"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	authz "github.com/aserto-dev/go-grpc-authz/aserto/authorizer/authorizer/v1"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
)

type ExecQueryCmd struct {
	Statement string `arg:"stmt" name:"stmt" required:"" help:"query statement"`
	Input     string `name:"input" optional:"" help:"query input context"`
}

func (cmd *ExecQueryCmd) Run(c *cc.CommonCtx) error {
	client, err := authorizer.New(
		c.Context,
		aserto.WithAddr(c.AuthorizerService()),
		aserto.WithAPIKeyAuth(c.AuthorizerAPIKey()),
		aserto.WithTenantID(c.TenantID()),
	)
	if err != nil {
		return err
	}

	resp, err := client.Query(c.Context, &authz.QueryRequest{
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
