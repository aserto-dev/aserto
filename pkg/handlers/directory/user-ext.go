package directory

import (
	"fmt"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	dir "github.com/aserto-dev/go-grpc/aserto/authorizer/directory/v1"

	"google.golang.org/protobuf/types/known/structpb"
)

type GetUserPropsCmd struct {
	UserID string `arg:"id" name:"id" required:"" help:"user id or identifier"`
}

func (cmd *GetUserPropsCmd) Run(c *cc.CommonCtx) error {
	client, identity, err := NewClientWithIdentity(c, cmd.UserID)
	if err != nil {
		return err
	}

	resp, err := client.Directory.GetUserProperties(
		c.Context,
		&dir.GetUserPropertiesRequest{Id: identity.Id},
	)
	if err != nil {
		return err
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp.Results)
}

type GetUserRolesCmd struct {
	UserID string `arg:"id" name:"id" required:"" help:"user id or identifier"`
}

func (cmd *GetUserRolesCmd) Run(c *cc.CommonCtx) error {
	client, identity, err := NewClientWithIdentity(c, cmd.UserID)
	if err != nil {
		return err
	}

	resp, err := client.Directory.GetUserRoles(
		c.Context,
		&dir.GetUserRolesRequest{Id: identity.Id},
	)
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
	client, identity, err := NewClientWithIdentity(c, cmd.UserID)
	if err != nil {
		return err
	}

	resp, err := client.Directory.GetUserPermissions(
		c.Context,
		&dir.GetUserPermissionsRequest{Id: identity.Id},
	)
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
	client, identity, err := NewClientWithIdentity(c, cmd.UserID)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.ErrWriter, "set property %s=%s\n", cmd.Key, cmd.Value)
	if _, err := client.Directory.SetUserProperty(
		c.Context,
		&dir.SetUserPropertyRequest{
			Id:    identity.Id,
			Key:   cmd.Key,
			Value: structpb.NewStringValue(cmd.Value),
		},
	); err != nil {
		return err
	}

	return nil
}

type SetUserRoleCmd struct {
	UserID string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	Key    string `arg:"role" name:"role" required:"" help:"role name"`
}

func (cmd *SetUserRoleCmd) Run(c *cc.CommonCtx) error {
	client, identity, err := NewClientWithIdentity(c, cmd.UserID)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.ErrWriter, "set role %s\n", cmd.Key)
	if _, err := client.Directory.SetUserRole(
		c.Context,
		&dir.SetUserRoleRequest{
			Id:   identity.Id,
			Role: cmd.Key,
		},
	); err != nil {
		return err
	}
	return nil
}

type SetUserPermCmd struct {
	UserID string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	Key    string `arg:"perm" name:"perm" required:"" help:"permission name"`
}

func (cmd *SetUserPermCmd) Run(c *cc.CommonCtx) error {
	client, identity, err := NewClientWithIdentity(c, cmd.UserID)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.ErrWriter, "set permission %s\n", cmd.Key)
	if _, err := client.Directory.SetUserPermission(
		c.Context,
		&dir.SetUserPermissionRequest{
			Id:         identity.Id,
			Permission: cmd.Key,
		},
	); err != nil {
		return err
	}
	return nil
}

type DelUserPropCmd struct {
	UserID string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	Key    string `arg:"prop" name:"prop" required:"" help:"property name"`
}

func (cmd *DelUserPropCmd) Run(c *cc.CommonCtx) error {
	client, identity, err := NewClientWithIdentity(c, cmd.UserID)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.ErrWriter, "removing property [%s]\n", cmd.Key)
	if _, err := client.Directory.DeleteUserProperty(
		c.Context,
		&dir.DeleteUserPropertyRequest{
			Id:  identity.Id,
			Key: cmd.Key,
		},
	); err != nil {
		return err
	}

	return nil
}

type DelUserRoleCmd struct {
	UserID string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	Key    string `arg:"role" name:"role" required:"" help:"role name"`
}

func (cmd *DelUserRoleCmd) Run(c *cc.CommonCtx) error {
	client, identity, err := NewClientWithIdentity(c, cmd.UserID)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.ErrWriter, "removing role [%s]\n", cmd.Key)
	if _, err := client.Directory.DeleteUserRole(
		c.Context,
		&dir.DeleteUserRoleRequest{
			Id:   identity.Id,
			Role: cmd.Key,
		},
	); err != nil {
		return err
	}

	return nil
}

type DelUserPermCmd struct {
	UserID string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	Key    string `arg:"perm" name:"perm" required:"" help:"permission name"`
}

func (cmd *DelUserPermCmd) Run(c *cc.CommonCtx) error {
	client, identity, err := NewClientWithIdentity(c, cmd.UserID)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.ErrWriter, "removing permission [%s]\n", cmd.Key)
	if _, err := client.Directory.DeleteUserPermission(
		c.Context,
		&dir.DeleteUserPermissionRequest{
			Id:         identity.Id,
			Permission: cmd.Key,
		},
	); err != nil {
		return err
	}

	return nil
}
