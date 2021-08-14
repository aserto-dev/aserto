package directory

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/grpcc"
	"github.com/aserto-dev/aserto/pkg/grpcc/authorizer"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	dir "github.com/aserto-dev/proto/aserto/authorizer/directory"
	"github.com/pkg/errors"
)

type SetUserCmd struct {
	UserID  string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	Disable bool   `xor:"group" help:"disable user"`
	Enable  bool   `xor:"group" help:"enable user"`
	Output  bool   `optional:"" help:"output updated user object"`
}

func (cmd *SetUserCmd) Run(c *cc.CommonCtx) error {
	if !cmd.Disable && !cmd.Enable {
		return errors.Errorf("must provide either --disable or --enable flag")
	}

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
		Identity: cmd.UserID,
	})
	if err != nil {
		return errors.Wrapf(err, "resolve identity")
	}

	resp, err := dirClient.GetUser(ctx, &dir.GetUserRequest{Id: idResp.Id})
	if err != nil {
		return err
	}

	user := resp.Result
	user.Enabled = &cmd.Enable

	updResp, err := dirClient.UpdateUser(ctx, &dir.UpdateUserRequest{
		Id:   idResp.Id,
		User: user,
	})
	if err != nil {
		return err
	}

	if cmd.Output {
		return jsonx.OutputJSONPB(c.OutWriter, updResp.Result)
	}

	return nil
}
