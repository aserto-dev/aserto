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

	resp, err := dirClient.ListUserApplications(ctx, &dir.ListUserApplicationsRequest{Id: cmd.UserID})
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

	_, err = dirClient.DeleteUserApplication(ctx, &dir.DeleteUserApplicationRequest{
		Id:   cmd.UserID,
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

	resp, err := dirClient.GetApplProperties(ctx, &dir.GetApplPropertiesRequest{Id: cmd.UserID, Name: cmd.AppName})
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

	resp, err := dirClient.GetApplRoles(ctx, &dir.GetApplRolesRequest{Id: cmd.UserID, Name: cmd.AppName})
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

	resp, err := dirClient.GetApplPermissions(ctx, &dir.GetApplPermissionsRequest{Id: cmd.UserID, Name: cmd.AppName})
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
	Value   string `optional:"" help:"set property using string value"`
	Stdin   bool   `optional:"" name:"stdin" help:"set property using from --stdin"`
	File    string `optional:"" type:"existingfile" help:"set property using file content"`
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
	if _, err := dirClient.SetApplProperty(ctx, &dir.SetApplPropertyRequest{
		Id:    cmd.UserID,
		Name:  cmd.AppName,
		Key:   cmd.Key,
		Value: value,
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

	if _, err := dirClient.SetApplRole(ctx, &dir.SetApplRoleRequest{
		Id:   cmd.UserID,
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

	if _, err := dirClient.SetApplPermission(ctx, &dir.SetApplPermissionRequest{
		Id:         cmd.UserID,
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

	fmt.Fprintf(c.ErrWriter, "remove property %s\n", cmd.Key)
	if _, err := dirClient.DeleteApplProperty(ctx, &dir.DeleteApplPropertyRequest{
		Id:   cmd.UserID,
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

	fmt.Fprintf(c.ErrWriter, "remove role %s\n", cmd.Key)
	if _, err := dirClient.DeleteApplRole(ctx, &dir.DeleteApplRoleRequest{
		Id:   cmd.UserID,
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

	fmt.Fprintf(c.ErrWriter, "remove permission %s\n", cmd.Key)
	if _, err := dirClient.DeleteApplPermission(ctx, &dir.DeleteApplPermissionRequest{
		Id:         cmd.UserID,
		Name:       cmd.AppName,
		Permission: cmd.Key,
	}); err != nil {
		return err
	}
	return nil
}
