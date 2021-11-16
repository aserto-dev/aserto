package directory

import (
	aserto "github.com/aserto-dev/aserto-go/client"
	"github.com/aserto-dev/aserto-go/client/grpc"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	dir "github.com/aserto-dev/go-grpc/aserto/authorizer/directory/v1"

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

	client, err := grpc.New(
		c.Context,
		aserto.WithAddr(c.AuthorizerService()),
		aserto.WithTokenAuth(c.AccessToken()),
		aserto.WithTenantID(c.TenantID()),
	)
	if err != nil {
		return err
	}

	dirClient := client.Directory
	idResp, err := dirClient.GetIdentity(c.Context, &dir.GetIdentityRequest{
		Identity: cmd.UserID,
	})
	if err != nil {
		return errors.Wrapf(err, "resolve identity")
	}

	resp, err := dirClient.GetUser(c.Context, &dir.GetUserRequest{Id: idResp.Id})
	if err != nil {
		return err
	}

	user := resp.Result
	user.Enabled = &cmd.Enable

	updResp, err := dirClient.UpdateUser(c.Context, &dir.UpdateUserRequest{
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
