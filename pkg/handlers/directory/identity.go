package directory

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/go-aserto/client/authorizer"

	dir "github.com/aserto-dev/go-grpc/aserto/authorizer/directory/v1"

	"github.com/pkg/errors"
)

func NewClientWithIdentity(c *cc.CommonCtx, id string) (*authorizer.Client, *dir.GetIdentityResponse, error) {
	return nil, nil, errors.Errorf("NOT IMPLEMENTED")

	// client, err := c.AuthorizerClient()
	// if err != nil {
	// 	return nil, nil, err
	// }

	// identity, err := client.Directory.GetIdentity(c.Context, &dir.GetIdentityRequest{Identity: id})
	// if err != nil {
	// 	return nil, nil, errors.Wrapf(err, "resolve identity")
	// }

	// return client, identity, nil
}
