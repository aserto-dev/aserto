package authorizer

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	authz "github.com/aserto-dev/go-grpc-authz/aserto/authorizer/authorizer/v1"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

type ExecQueryCmd struct {
	PolicyID  string `name:"policy_id" required:"" help:"policy id"`
	Statement string `arg:"stmt" name:"stmt" required:"" help:"query statement"`
	Input     string `name:"input" optional:"" help:"query input context"`
}

func (cmd *ExecQueryCmd) Run(c *cc.CommonCtx) error {
	client, err := c.AuthorizerClient()
	if err != nil {
		return err
	}

	resp, err := client.Authorizer.Query(c.Context, &authz.QueryRequest{
		Query: cmd.Statement,
		Input: cmd.Input,
		IdentityContext: &api.IdentityContext{
			Identity: "",
			Type:     api.IdentityType_IDENTITY_TYPE_NONE,
		},
		PolicyContext: &api.PolicyContext{
			Id: cmd.PolicyID,
		},
		ResourceContext: &structpb.Struct{},
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

	return jsonx.OutputJSONPB(c.UI.Output(), resp)
}
