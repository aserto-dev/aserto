package authorizer

import (
	"github.com/aserto-dev/aserto-go/client/grpc/authorizer"
	"github.com/aserto-dev/aserto/pkg/cc"
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
	client, err := authorizer.New(c.Context, c.AuthorizerSvcConnectionOptions()...)
	if err != nil {
		return err
	}

	resp, err := client.DecisionTree(c.Context, &authz.DecisionTreeRequest{
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
