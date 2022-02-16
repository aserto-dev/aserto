package authorizer

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	authz "github.com/aserto-dev/go-grpc-authz/aserto/authorizer/authorizer/v1"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/go-utils/pb"
)

type EvalDecisionCmd struct {
	PolicyID  string
	Path      string
	Decisions []string
}

func (cmd *EvalDecisionCmd) Run(c *cc.CommonCtx) error {
	client, err := c.AuthorizerClient()
	if err != nil {
		return err
	}

	resp, err := client.Authorizer.Is(c.Context, &authz.IsRequest{
		PolicyContext: &api.PolicyContext{
			Id:        cmd.PolicyID,
			Path:      cmd.Path,
			Decisions: cmd.Decisions,
		},
		IdentityContext: &api.IdentityContext{
			Identity: "",
			Type:     api.IdentityType_IDENTITY_TYPE_NONE,
		},
		ResourceContext: pb.NewStruct(),
	})
	if err != nil {
		return err
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}
