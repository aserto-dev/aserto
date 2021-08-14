package directory

import (
	"fmt"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/grpcc"
	"github.com/aserto-dev/aserto/pkg/grpcc/authorizer"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	dir "github.com/aserto-dev/proto/aserto/authorizer/directory"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type ListUserAppsCmd struct {
	UserID string `arg:"id" name:"id" required:"" help:"user id or identifier"`
}

func (cmd *ListUserAppsCmd) Run(c *cc.CommonCtx) error {
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

	resp, err := dirClient.ListUserApplications(ctx, &dir.ListUserApplicationsRequest{Id: idResp.Id})
	if err != nil {
		return err
	}
	if resp.Results == nil {
		resp.Results = []string{}
	}

	return jsonx.OutputJSON(c.OutWriter, resp.Results)
}

type SetUserAppCmd struct {
	UserID  string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	AppName string `arg:"name" name:"name" required:"" help:"application name"`
}

func (cmd *SetUserAppCmd) Run(c *cc.CommonCtx) error {
	return status.Error(codes.Unimplemented, "not implemented")
}

type DelUserAppCmd struct {
	UserID  string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	AppName string `arg:"name" name:"name" required:"" help:"application name"`
}

func (cmd *DelUserAppCmd) Run(c *cc.CommonCtx) error {
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

	_, err = dirClient.DeleteUserApplication(ctx, &dir.DeleteUserApplicationRequest{
		Id:   idResp.Id,
		Name: cmd.AppName,
	})
	if err != nil {
		return err
	}

	return nil
}

type GetApplPropsCmd struct {
	UserID  string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	AppName string `arg:"name" name:"name" required:"" help:"application name"`
}

func (cmd *GetApplPropsCmd) Run(c *cc.CommonCtx) error {
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

	resp, err := dirClient.GetApplProperties(ctx, &dir.GetApplPropertiesRequest{Id: idResp.Id, Name: cmd.AppName})
	if err != nil {
		return err
	}
	return jsonx.OutputJSONPB(c.OutWriter, resp.Results)
}

type GetApplRolesCmd struct {
	UserID  string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	AppName string `arg:"name" name:"name" required:"" help:"application name"`
}

func (cmd *GetApplRolesCmd) Run(c *cc.CommonCtx) error {
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

	resp, err := dirClient.GetApplRoles(ctx, &dir.GetApplRolesRequest{Id: idResp.Id, Name: cmd.AppName})
	if err != nil {
		return err
	}
	if resp.Results == nil {
		resp.Results = []string{}
	}
	return jsonx.OutputJSON(c.OutWriter, resp.Results)
}

type GetApplPermsCmd struct {
	UserID  string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	AppName string `arg:"name" name:"name" required:"" help:"application name"`
}

func (cmd *GetApplPermsCmd) Run(c *cc.CommonCtx) error {
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

	resp, err := dirClient.GetApplPermissions(ctx, &dir.GetApplPermissionsRequest{Id: idResp.Id, Name: cmd.AppName})
	if err != nil {
		return err
	}
	if resp.Results == nil {
		resp.Results = []string{}
	}
	return jsonx.OutputJSON(c.OutWriter, resp.Results)
}

type SetApplPropCmd struct {
	UserID  string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	AppName string `arg:"name" name:"name" required:"" help:"application name"`
	Key     string `arg:"key" name:"key" required:"" help:"property key"`
	Value   string `required:"" help:"set property using string value"`
}

func (cmd *SetApplPropCmd) Run(c *cc.CommonCtx) error {
	if _, err := uuid.Parse(cmd.UserID); err != nil {
		return errors.Errorf("argument provided is not a valid user id")
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

	fmt.Fprintf(c.ErrWriter, "set property %s\n", cmd.Key)
	if _, err := dirClient.SetApplProperty(ctx, &dir.SetApplPropertyRequest{
		Id:    idResp.Id,
		Name:  cmd.AppName,
		Key:   cmd.Key,
		Value: structpb.NewStringValue(cmd.Value),
	}); err != nil {
		return err
	}

	return nil
}

type SetApplRoleCmd struct {
	UserID  string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	AppName string `arg:"name" name:"name" required:"" help:"application name"`
	Key     string `arg:"role" name:"role" required:"" help:"role name"`
}

func (cmd *SetApplRoleCmd) Run(c *cc.CommonCtx) error {
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

	if _, err := dirClient.SetApplRole(ctx, &dir.SetApplRoleRequest{
		Id:   idResp.Id,
		Name: cmd.AppName,
		Role: cmd.Key,
	}); err != nil {
		return err
	}
	return nil
}

type SetApplPermCmd struct {
	UserID  string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	AppName string `arg:"name" name:"name" required:"" help:"application name"`
	Key     string `arg:"perm" name:"perm" required:"" help:"permission name"`
}

func (cmd *SetApplPermCmd) Run(c *cc.CommonCtx) error {
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

	if _, err := dirClient.SetApplPermission(ctx, &dir.SetApplPermissionRequest{
		Id:         idResp.Id,
		Name:       cmd.AppName,
		Permission: cmd.Key,
	}); err != nil {
		return err
	}

	return nil
}

type DelApplPropCmd struct {
	UserID  string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	AppName string `arg:"name" name:"name" required:"" help:"application name"`
	Key     string `arg:"prop" name:"prop" required:"" help:"property name"`
}

func (cmd *DelApplPropCmd) Run(c *cc.CommonCtx) error {
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

	fmt.Fprintf(c.ErrWriter, "remove property %s\n", cmd.Key)
	if _, err := dirClient.DeleteApplProperty(ctx, &dir.DeleteApplPropertyRequest{
		Id:   idResp.Id,
		Name: cmd.AppName,
		Key:  cmd.Key,
	}); err != nil {
		return err
	}

	return nil
}

type DelApplRoleCmd struct {
	UserID  string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	AppName string `arg:"name" name:"name" required:"" help:"application name"`
	Key     string `arg:"role" name:"role" required:"" help:"role name"`
}

func (cmd *DelApplRoleCmd) Run(c *cc.CommonCtx) error {
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

	fmt.Fprintf(c.ErrWriter, "remove role %s\n", cmd.Key)
	if _, err := dirClient.DeleteApplRole(ctx, &dir.DeleteApplRoleRequest{
		Id:   idResp.Id,
		Name: cmd.AppName,
		Role: cmd.Key,
	}); err != nil {
		return err
	}

	return nil
}

type DelApplPermCmd struct {
	UserID  string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	AppName string `arg:"name" name:"name" required:"" help:"application name"`
	Key     string `arg:"perm" name:"perm" required:"" help:"permission name"`
}

func (cmd *DelApplPermCmd) Run(c *cc.CommonCtx) error {
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

	fmt.Fprintf(c.ErrWriter, "remove permission %s\n", cmd.Key)
	if _, err := dirClient.DeleteApplPermission(ctx, &dir.DeleteApplPermissionRequest{
		Id:         idResp.Id,
		Name:       cmd.AppName,
		Permission: cmd.Key,
	}); err != nil {
		return err
	}
	return nil
}
