package directory

import (
	aserto "github.com/aserto-dev/aserto-go/client"
	"github.com/aserto-dev/aserto-go/client/authorizer"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	dir "github.com/aserto-dev/go-grpc/aserto/authorizer/directory/v1"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type ListUsersCmd struct {
	Base   bool     `name:"base" optional:"" help:"return base user object (without extensions)"`
	Count  bool     `name:"count" optional:"" help:"only return user count"`
	Fields []string `name:"fields" optional:"" help:"fields mask, like --fields=id,email"`
}

func (cmd *ListUsersCmd) Run(c *cc.CommonCtx) error {
	client, err := authorizer.New(
		c.Context,
		aserto.WithAddr(c.AuthorizerService()),
		aserto.WithTokenAuth(c.AccessToken()),
		aserto.WithTenantID(c.TenantID()),
	)
	if err != nil {
		return err
	}
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

	opts := jsonx.MaskedMarshalOpts()

	for {
		resp, err := client.Directory.ListUsers(c.Context, &dir.ListUsersRequest{
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
	_, identity, err := NewClientWithIdentity(c, cmd.Identity)
	if err != nil {
		return err
	}

	return jsonx.OutputJSONPB(c.OutWriter, identity)
}

type GetUserCmd struct {
	ID   string `arg:"id" name:"id" required:"" help:"user id or identity"`
	Base bool   `name:"base" optional:"" help:"return base user object (without extensions)"`
}

func (cmd *GetUserCmd) Run(c *cc.CommonCtx) error {
	client, identity, err := NewClientWithIdentity(c, cmd.ID)
	if err != nil {
		return err
	}

	resp, err := client.Directory.GetUser(c.Context, &dir.GetUserRequest{
		Id:   identity.Id,
		Base: cmd.Base,
	})
	if err != nil {
		return errors.Wrapf(err, "get user")
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}
