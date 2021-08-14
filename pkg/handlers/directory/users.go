package directory

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/grpcc"
	"github.com/aserto-dev/aserto/pkg/grpcc/authorizer"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/proto/aserto/api"
	dir "github.com/aserto-dev/proto/aserto/authorizer/directory"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type ListUsersCmd struct {
	Base   bool     `name:"base" optional:"" help:"return base user object (without extensions)"`
	Count  bool     `name:"count" optional:"" help:"only return user count"`
	Fields []string `name:"fields" optional:"" help:"fields mask, like --fields=id,email"`
}

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

	mask, err := fieldmaskpb.New(&api.User{}, cmd.Fields...)
	if err != nil {
		return err
	}

	pageSize := int32(100)
	if cmd.Count {
		pageSize = int32(-2)
	}

	token := ""
	first := true
	count := int32(0)

	dirClient := conn.DirectoryClient()

	opts := jsonx.MaskedMarshalOpts()

	for {
		resp, err := dirClient.ListUsers(ctx, &dir.ListUsersRequest{
			Page: &api.PaginationRequest{
				Size:  pageSize,
				Token: token,
			},
			Fields: &api.Fields{
				Mask: mask,
			},
			Base: cmd.Base,
		})

		if err != nil {
			return errors.Wrapf(err, "list users")
		}

		if cmd.Count {
			return jsonx.OutputJSONPB(c.OutWriter, resp.Page, opts)
		}

		if first {
			_, _ = c.OutWriter.Write([]byte("[\n"))
			first = false
		}

		for _, u := range resp.Results {
			if count > 0 {
				_, _ = c.OutWriter.Write([]byte(",\n"))
			}

			_ = jsonx.EncodeJSONPB(c.OutWriter, u, opts)

			count++
		}

		if resp.Page.NextToken == "" {
			break
		}

		token = resp.Page.NextToken
	}

	_, _ = c.OutWriter.Write([]byte("\n]\n"))

	return nil
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
	ID   string `arg:"id" name:"id" required:"" help:"user id or identity"`
	Base bool   `name:"base" optional:"" help:"return base user object (without extensions)"`
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
		Id:   idResp.Id,
		Base: cmd.Base,
	})
	if err != nil {
		return errors.Wrapf(err, "get user")
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}
