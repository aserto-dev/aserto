package directory

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/grpcc"
	"github.com/aserto-dev/aserto/pkg/grpcc/authorizer"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/proto/aserto/api"
	dir "github.com/aserto-dev/proto/aserto/authorizer/directory"
	"github.com/pkg/errors"
)

type ListUsersCmd struct {
}

// TODO : add mask
// TODO : add pagination (instead of -1)
// TODO : add filtering?
func (cmd *ListUsersCmd) Run(c *cc.CommonCtx) error {
	conn, err := authorizer.Connection(
		c.Context,
		c.AuthorizerService(),
		grpcc.NewTokenAuth(c.AccessToken()),
	)
	if err != nil {
		return err
	}

	ctx := grpcc.SetTenantContext(c.Context, c.TenantID())

	dirClient := conn.DirectoryClient()
	resp, err := dirClient.ListUsers(ctx, &dir.ListUsersRequest{
		Page: &api.PaginationRequest{
			Size: -1,
		},
	})

	if err != nil {
		return err
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}

type GetIdentityCmd struct {
	Identity string `arg:"" name:"identity" required:"" help:"User identity"`
}

func (cmd *GetIdentityCmd) Run(c *cc.CommonCtx) error {
	conn, err := authorizer.Connection(
		c.Context,
		c.AuthorizerService(),
		grpcc.NewTokenAuth(c.AccessToken()),
	)
	if err != nil {
		return err
	}

	ctx := grpcc.SetTenantContext(c.Context, c.TenantID())

	dirClient := conn.DirectoryClient()
	resp, err := dirClient.GetIdentity(ctx, &dir.GetIdentityRequest{
		Identity: cmd.Identity,
	})

	if err != nil {
		return err
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}

type GetUserCmd struct {
	ID string `arg:"id" name:"id" required:"" help:"user id or identity"`
}

func (cmd *GetUserCmd) Run(c *cc.CommonCtx) error {
	conn, err := authorizer.Connection(
		c.Context,
		c.AuthorizerService(),
		grpcc.NewTokenAuth(c.AccessToken()),
	)
	if err != nil {
		return err
	}

	ctx := grpcc.SetTenantContext(c.Context, c.TenantID())

	dirClient := conn.DirectoryClient()
	idResp, err := dirClient.GetIdentity(ctx, &dir.GetIdentityRequest{
		Identity: cmd.ID,
	})
	if err != nil {
		return errors.Wrapf(err, "resolve identity")
	}

	resp, err := dirClient.GetUser(ctx, &dir.GetUserRequest{
		Id: idResp.Id,
	})
	if err != nil {
		return errors.Wrapf(err, "get user")
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}
