package directory

import (
	"github.com/aserto-dev/aserto-go/client/authorizer"
	"github.com/aserto-dev/aserto/pkg/cc"

	dir "github.com/aserto-dev/go-grpc/aserto/authorizer/directory/v1"

	"github.com/pkg/errors"
)

func NewClientWithIdentity(c *cc.CommonCtx, id string) (*authorizer.Client, *dir.GetIdentityResponse, error) {
	client, err := authorizer.New(c.Context, c.AuthorizerSvcConnectionOptions()...)
	if err != nil {
		return nil, nil, err
	}

	identity, err := client.Directory.GetIdentity(c.Context, &dir.GetIdentityRequest{Identity: id})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "resolve identity")
	}

	return client, identity, nil
}
