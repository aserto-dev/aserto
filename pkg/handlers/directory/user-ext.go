package directory

import (
	"fmt"
	"io"
	"os"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/aserto/pkg/pb"
	dir "github.com/aserto-dev/go-grpc/aserto/authorizer/directory/v1"
	"github.com/pkg/errors"

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

	return jsonx.OutputJSONPB(c.UI.Output(), resp.Results)
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

	return jsonx.OutputJSON(c.UI.Output(), resp.Results)
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

	return jsonx.OutputJSON(c.UI.Output(), resp.Results)
}

type SetUserPropCmd struct {
	UserID string         `arg:"id" name:"id" required:"" help:"user id or identifier"`
	Key    string         `arg:"key" name:"key" required:"" help:"property key"`
	Value  structpb.Value `xor:"group" required:"" name:"value" help:"set property value using json data from argument"`
	Stdin  bool           `xor:"group" required:"" name:"stdin" help:"set property value using json data from --stdin"`
	File   string         `xor:"group" required:"" name:"file" type:"existingfile" help:"set property value using json data from input file"`
}

func (cmd *SetUserPropCmd) Run(c *cc.CommonCtx) error {
	var (
		value *structpb.Value
		buf   io.Reader
		err   error
	)

	client, identity, err := NewClientWithIdentity(c, cmd.UserID)
	if err != nil {
		return err
	}

	switch {
	case cmd.Stdin:
		fmt.Fprintf(c.UI.Err(), "reading stdin\n")
		buf = os.Stdin

		value, err = pb.BufToValue(buf)
		if err != nil {
			return errors.Wrapf(err, "unmarshal stdin")
		}

	case cmd.File != "":
		fmt.Fprintf(c.UI.Err(), "reading file [%s]\n", cmd.File)
		buf, err = os.Open(cmd.File)
		if err != nil {
			return errors.Wrapf(err, "opening file [%s]", cmd.File)
		}
		value, err = pb.BufToValue(buf)
		if err != nil {
			return errors.Wrapf(err, "unmarshal file [%s]", cmd.File)
		}

	default:
		value = &cmd.Value
	}

	fmt.Fprintf(c.UI.Err(), "set property [%s]=[%s]\n", cmd.Key, value.String())
	if _, err := client.Directory.SetUserProperty(
		c.Context,
		&dir.SetUserPropertyRequest{
			Id:    identity.Id,
			Key:   cmd.Key,
			Value: value,
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

	fmt.Fprintf(c.UI.Err(), "set role %s\n", cmd.Key)
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

	fmt.Fprintf(c.UI.Err(), "set permission %s\n", cmd.Key)
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

	fmt.Fprintf(c.UI.Err(), "removing property [%s]\n", cmd.Key)
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

	fmt.Fprintf(c.UI.Err(), "removing role [%s]\n", cmd.Key)
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

	fmt.Fprintf(c.UI.Err(), "removing permission [%s]\n", cmd.Key)
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
