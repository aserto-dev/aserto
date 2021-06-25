package directory

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/grpcc"
	"github.com/aserto-dev/aserto/pkg/grpcc/authorizer"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/aserto/pkg/pb"
	dir "github.com/aserto-dev/proto/aserto/authorizer/directory"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/structpb"
)

type GetUserPropsCmd struct {
	UserID string `arg:"id" name:"id" required:"" help:"user id or identifier"`
}

func (cmd *GetUserPropsCmd) Run(c *cc.CommonCtx) error {
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

	resp, err := dirClient.GetUserProperties(ctx, &dir.GetUserPropertiesRequest{Id: idResp.Id})
	if err != nil {
		return err
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp.Results)
}

type GetUserRolesCmd struct {
	UserID string `arg:"id" name:"id" required:"" help:"user id or identifier"`
}

func (cmd *GetUserRolesCmd) Run(c *cc.CommonCtx) error {
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

	resp, err := dirClient.GetUserRoles(ctx, &dir.GetUserRolesRequest{Id: idResp.Id})
	if err != nil {
		return err
	}

	if resp.Results == nil {
		resp.Results = []string{}
	}

	return jsonx.OutputJSON(c.OutWriter, resp.Results)
}

type GetUserPermsCmd struct {
	UserID string `arg:"id" name:"id" required:"" help:"user id or identifier"`
}

func (cmd *GetUserPermsCmd) Run(c *cc.CommonCtx) error {
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

	resp, err := dirClient.GetUserPermissions(ctx, &dir.GetUserPermissionsRequest{Id: idResp.Id})
	if err != nil {
		return err
	}

	if resp.Results == nil {
		resp.Results = []string{}
	}

	return jsonx.OutputJSON(c.OutWriter, resp.Results)
}

type SetUserPropCmd struct {
	UserID string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	Key    string `arg:"key" name:"key" required:"" help:"property key"`
	Value  string `optional:"" help:"set property using string value"`
	Stdin  bool   `optional:"" name:"stdin" help:"set property using from --stdin"`
	File   string `optional:"" type:"existingfile" help:"set property using file content"`
}

func (cmd *SetUserPropCmd) Run(c *cc.CommonCtx) error {
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

	var (
		value *structpb.Value
		buf   io.Reader
	)

	switch {
	case cmd.Stdin:
		fmt.Fprintf(c.ErrWriter, "reading stdin\n")
		buf = os.Stdin

	case cmd.File != "":
		fmt.Fprintf(c.ErrWriter, "reading file [%s]\n", cmd.File)
		buf, err = os.Open(cmd.File)
		if err != nil {
			return errors.Wrapf(err, "opening file [%s]", cmd.File)
		}
	case cmd.Value != "":
		fmt.Fprintf(c.ErrWriter, "reading string value\n")
		buf = strings.NewReader(cmd.Value)

	default:
		return errors.Errorf("no input option specified, [--stdin | --file=filepath | --value=string]")
	}

	if buf == nil {
		value = &structpb.Value{}
	} else if err := pb.BufToProto(buf, value); err != nil {
		return err
	}

	fmt.Fprintf(c.ErrWriter, "set property %s\n", cmd.Key)
	if _, err := dirClient.SetUserProperty(ctx, &dir.SetUserPropertyRequest{
		Id:    idResp.Id,
		Key:   cmd.Key,
		Value: value,
	}); err != nil {
		return err
	}

	return nil
}

type SetUserRoleCmd struct {
	UserID string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	Key    string `arg:"role" name:"role" required:"" help:"role name"`
}

func (cmd *SetUserRoleCmd) Run(c *cc.CommonCtx) error {
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

	fmt.Fprintf(c.ErrWriter, "set role %s\n", cmd.Key)
	if _, err := dirClient.SetUserRole(ctx, &dir.SetUserRoleRequest{
		Id:   idResp.Id,
		Role: cmd.Key,
	}); err != nil {
		return err
	}
	return nil
}

type SetUserPermCmd struct {
	UserID string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	Key    string `arg:"perm" name:"perm" required:"" help:"permission name"`
}

func (cmd *SetUserPermCmd) Run(c *cc.CommonCtx) error {
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

	fmt.Fprintf(c.ErrWriter, "set permission %s\n", cmd.Key)
	if _, err := dirClient.SetUserPermission(ctx, &dir.SetUserPermissionRequest{
		Id:         idResp.Id,
		Permission: cmd.Key,
	}); err != nil {
		return err
	}
	return nil
}

type DelUserPropCmd struct {
	UserID string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	Key    string `arg:"prop" name:"prop" required:"" help:"property name"`
}

func (cmd *DelUserPropCmd) Run(c *cc.CommonCtx) error {
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

	fmt.Fprintf(c.ErrWriter, "removing property [%s]\n", cmd.Key)
	if _, err := dirClient.DeleteUserProperty(ctx, &dir.DeleteUserPropertyRequest{
		Id:  idResp.Id,
		Key: cmd.Key,
	}); err != nil {
		return err
	}

	return nil
}

type DelUserRoleCmd struct {
	UserID string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	Key    string `arg:"role" name:"role" required:"" help:"role name"`
}

func (cmd *DelUserRoleCmd) Run(c *cc.CommonCtx) error {
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

	fmt.Fprintf(c.ErrWriter, "removing role [%s]\n", cmd.Key)
	if _, err := dirClient.DeleteUserRole(ctx, &dir.DeleteUserRoleRequest{
		Id:   idResp.Id,
		Role: cmd.Key,
	}); err != nil {
		return err
	}

	return nil
}

type DelUserPermCmd struct {
	UserID string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	Key    string `arg:"perm" name:"perm" required:"" help:"permission name"`
}

func (cmd *DelUserPermCmd) Run(c *cc.CommonCtx) error {
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

	fmt.Fprintf(c.ErrWriter, "removing permission [%s]\n", cmd.Key)
	if _, err := dirClient.DeleteUserPermission(ctx, &dir.DeleteUserPermissionRequest{
		Id:         idResp.Id,
		Permission: cmd.Key,
	}); err != nil {
		return err
	}

	return nil
}
