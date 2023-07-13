package directory

import (
	"io"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

func OutputJSONResults(results []string, writer io.Writer) error {
	if results == nil {
		results = []string{}
	}

	return jsonx.OutputJSON(writer, results)
}

type ListUserAppsCmd struct {
	UserID string `arg:"id" name:"id" required:"" help:"user id or identifier"`
}

func (cmd *ListUserAppsCmd) Run(c *cc.CommonCtx) error {
	return errors.Errorf("NOT IMPLEMENTED")
	// client, identity, err := NewClientWithIdentity(c, cmd.UserID)
	// if err != nil {
	// 	return err
	// }

	// resp, err := client.Directory.ListUserApplications(
	// 	c.Context,
	// 	&dir.ListUserApplicationsRequest{Id: identity.Id},
	// )
	// if err != nil {
	// 	return err
	// }
	// return OutputJSONResults(resp.Results, c.UI.Output())
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
	return errors.Errorf("NOT IMPLEMENTED")

	// client, identity, err := NewClientWithIdentity(c, cmd.UserID)
	// if err != nil {
	// 	return err
	// }

	// _, err = client.Directory.DeleteUserApplication(
	// 	c.Context,
	// 	&dir.DeleteUserApplicationRequest{
	// 		Id:   identity.Id,
	// 		Name: cmd.AppName,
	// 	},
	// )
	// if err != nil {
	// 	return err
	// }

	// return nil
}

type GetApplPropsCmd struct {
	UserID  string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	AppName string `arg:"name" name:"name" required:"" help:"application name"`
}

func (cmd *GetApplPropsCmd) Run(c *cc.CommonCtx) error {
	return errors.Errorf("NOT IMPLEMENTED")

	// client, identity, err := NewClientWithIdentity(c, cmd.UserID)
	// if err != nil {
	// 	return err
	// }

	// resp, err := client.Directory.GetApplProperties(
	// 	c.Context,
	// 	&dir.GetApplPropertiesRequest{
	// 		Id:   identity.Id,
	// 		Name: cmd.AppName,
	// 	},
	// )
	// if err != nil {
	// 	return err
	// }
	// return jsonx.OutputJSONPB(c.UI.Output(), resp.Results)
}

type GetApplRolesCmd struct {
	UserID  string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	AppName string `arg:"name" name:"name" required:"" help:"application name"`
}

func (cmd *GetApplRolesCmd) Run(c *cc.CommonCtx) error {
	return errors.Errorf("NOT IMPLEMENTED")

	// client, identity, err := NewClientWithIdentity(c, cmd.UserID)
	// if err != nil {
	// 	return err
	// }

	// resp, err := client.Directory.GetApplRoles(
	// 	c.Context,
	// 	&dir.GetApplRolesRequest{
	// 		Id:   identity.Id,
	// 		Name: cmd.AppName,
	// 	},
	// )
	// if err != nil {
	// 	return err
	// }
	// return OutputJSONResults(resp.Results, c.UI.Output())
}

type GetApplPermsCmd struct {
	UserID  string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	AppName string `arg:"name" name:"name" required:"" help:"application name"`
}

func (cmd *GetApplPermsCmd) Run(c *cc.CommonCtx) error {
	return errors.Errorf("NOT IMPLEMENTED")

	// client, identity, err := NewClientWithIdentity(c, cmd.UserID)
	// if err != nil {
	// 	return err
	// }

	// resp, err := client.Directory.GetApplPermissions(
	// 	c.Context,
	// 	&dir.GetApplPermissionsRequest{
	// 		Id:   identity.Id,
	// 		Name: cmd.AppName,
	// 	},
	// )
	// if err != nil {
	// 	return err
	// }
	// return OutputJSONResults(resp.Results, c.UI.Output())
}

type SetApplPropCmd struct {
	UserID  string         `arg:"id" name:"id" required:"" help:"user id or identifier"`
	AppName string         `arg:"name" name:"name" required:"" help:"application name"`
	Key     string         `arg:"key" name:"key" required:"" help:"property key"`
	Value   structpb.Value `xor:"group" required:"" name:"value" help:"set property value using json data from argument"`
	Stdin   bool           `xor:"group" required:"" name:"stdin" help:"set property value using json data from --stdin"`
	File    string         `xor:"group" required:"" name:"file" type:"existingfile" help:"set property value using json data from input file"`
}

func (cmd *SetApplPropCmd) Run(c *cc.CommonCtx) error {
	return errors.Errorf("NOT IMPLEMENTED")

	// var (
	// 	value *structpb.Value
	// 	buf   io.Reader
	// 	err   error
	// )

	// client, identity, err := NewClientWithIdentity(c, cmd.UserID)
	// if err != nil {
	// 	return err
	// }

	// switch {
	// case cmd.Stdin:
	// 	fmt.Fprintf(c.UI.Err(), "reading stdin\n")
	// 	buf = os.Stdin

	// 	value, err = pb.BufToValue(buf)
	// 	if err != nil {
	// 		return errors.Wrapf(err, "unmarshal stdin")
	// 	}

	// case cmd.File != "":
	// 	fmt.Fprintf(c.UI.Err(), "reading file [%s]\n", cmd.File)
	// 	buf, err = os.Open(cmd.File)
	// 	if err != nil {
	// 		return errors.Wrapf(err, "opening file [%s]", cmd.File)
	// 	}
	// 	value, err = pb.BufToValue(buf)
	// 	if err != nil {
	// 		return errors.Wrapf(err, "unmarshal file [%s]", cmd.File)
	// 	}

	// default:
	// 	value = &cmd.Value
	// }

	// fmt.Fprintf(c.UI.Err(), "set property [%s]=[%s]\n", cmd.Key, value.String())
	// if _, err := client.Directory.SetApplProperty(
	// 	c.Context,
	// 	&dir.SetApplPropertyRequest{
	// 		Id:    identity.Id,
	// 		Name:  cmd.AppName,
	// 		Key:   cmd.Key,
	// 		Value: value,
	// 	},
	// ); err != nil {
	// 	return err
	// }

	// return nil
}

type SetApplRoleCmd struct {
	UserID  string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	AppName string `arg:"name" name:"name" required:"" help:"application name"`
	Key     string `arg:"role" name:"role" required:"" help:"role name"`
}

func (cmd *SetApplRoleCmd) Run(c *cc.CommonCtx) error {
	return errors.Errorf("NOT IMPLEMENTED")

	// client, identity, err := NewClientWithIdentity(c, cmd.UserID)
	// if err != nil {
	// 	return err
	// }

	// if _, err := client.Directory.SetApplRole(
	// 	c.Context,
	// 	&dir.SetApplRoleRequest{
	// 		Id:   identity.Id,
	// 		Name: cmd.AppName,
	// 		Role: cmd.Key,
	// 	},
	// ); err != nil {
	// 	return err
	// }
	// return nil
}

type SetApplPermCmd struct {
	UserID  string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	AppName string `arg:"name" name:"name" required:"" help:"application name"`
	Key     string `arg:"perm" name:"perm" required:"" help:"permission name"`
}

func (cmd *SetApplPermCmd) Run(c *cc.CommonCtx) error {
	return errors.Errorf("NOT IMPLEMENTED")

	// client, err := c.AuthorizerClient()
	// if err != nil {
	// 	return err
	// }

	// idResp, err := client.Directory.GetIdentity(c.Context, &dir.GetIdentityRequest{
	// 	Identity: cmd.UserID,
	// })
	// if err != nil {
	// 	return errors.Wrapf(err, "resolve identity")
	// }

	// if _, err := client.Directory.SetApplPermission(c.Context, &dir.SetApplPermissionRequest{
	// 	Id:         idResp.Id,
	// 	Name:       cmd.AppName,
	// 	Permission: cmd.Key,
	// }); err != nil {
	// 	return err
	// }

	// return nil
}

type DelApplPropCmd struct {
	UserID  string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	AppName string `arg:"name" name:"name" required:"" help:"application name"`
	Key     string `arg:"prop" name:"prop" required:"" help:"property name"`
}

func (cmd *DelApplPropCmd) Run(c *cc.CommonCtx) error {
	return errors.Errorf("NOT IMPLEMENTED")

	// client, identity, err := NewClientWithIdentity(c, cmd.UserID)
	// if err != nil {
	// 	return err
	// }

	// fmt.Fprintf(c.UI.Err(), "remove property %s\n", cmd.Key)
	// if _, err := client.Directory.DeleteApplProperty(
	// 	c.Context,
	// 	&dir.DeleteApplPropertyRequest{
	// 		Id:   identity.Id,
	// 		Name: cmd.AppName,
	// 		Key:  cmd.Key,
	// 	},
	// ); err != nil {
	// 	return err
	// }

	// return nil
}

type DelApplRoleCmd struct {
	UserID  string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	AppName string `arg:"name" name:"name" required:"" help:"application name"`
	Key     string `arg:"role" name:"role" required:"" help:"role name"`
}

func (cmd *DelApplRoleCmd) Run(c *cc.CommonCtx) error {
	return errors.Errorf("NOT IMPLEMENTED")

	// client, identity, err := NewClientWithIdentity(c, cmd.UserID)
	// if err != nil {
	// 	return err
	// }

	// fmt.Fprintf(c.UI.Err(), "remove role %s\n", cmd.Key)
	// if _, err := client.Directory.DeleteApplRole(
	// 	c.Context,
	// 	&dir.DeleteApplRoleRequest{
	// 		Id:   identity.Id,
	// 		Name: cmd.AppName,
	// 		Role: cmd.Key,
	// 	},
	// ); err != nil {
	// 	return err
	// }

	// return nil
}

type DelApplPermCmd struct {
	UserID  string `arg:"id" name:"id" required:"" help:"user id or identifier"`
	AppName string `arg:"name" name:"name" required:"" help:"application name"`
	Key     string `arg:"perm" name:"perm" required:"" help:"permission name"`
}

func (cmd *DelApplPermCmd) Run(c *cc.CommonCtx) error {
	return errors.Errorf("NOT IMPLEMENTED")

	// client, identity, err := NewClientWithIdentity(c, cmd.UserID)
	// if err != nil {
	// 	return err
	// }

	// fmt.Fprintf(c.UI.Err(), "remove permission %s\n", cmd.Key)
	// if _, err := client.Directory.DeleteApplPermission(
	// 	c.Context,
	// 	&dir.DeleteApplPermissionRequest{
	// 		Id:         identity.Id,
	// 		Name:       cmd.AppName,
	// 		Permission: cmd.Key,
	// 	},
	// ); err != nil {
	// 	return err
	// }
	// return nil
}
