package directory

import (
	"fmt"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/grpcc"
	"github.com/aserto-dev/aserto/pkg/grpcc/authorizer"
	"github.com/aserto-dev/aserto/pkg/jsonx"
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
	Value  string `required:"" help:"set property using string value"`
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

	fmt.Fprintf(c.ErrWriter, "set property %s=%s\n", cmd.Key, cmd.Value)
	if _, err := dirClient.SetUserProperty(ctx, &dir.SetUserPropertyRequest{
		Id:    idResp.Id,
		Key:   cmd.Key,
		Value: structpb.NewStringValue(cmd.Value),
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
