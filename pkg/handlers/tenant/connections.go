package tenant

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/aserto-dev/aserto-go/client/grpc/tenant"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	connection "github.com/aserto-dev/go-grpc/aserto/tenant/connection/v1"

	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type ListConnectionsCmd struct {
	Kind string `help:"provider kind"`
}

func (cmd ListConnectionsCmd) Run(c *cc.CommonCtx) error {
	client, err := tenant.New(c.Context, c.TenantSvcConnectionOptions()...)
	if err != nil {
		return err
	}

	resp, err := client.Connections.ListConnections(
		c.Context,
		&connection.ListConnectionsRequest{
			Kind: ProviderKind(cmd.Kind),
		})
	if err != nil {
		return errors.Wrapf(err, "list connections")
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}

type GetConnectionCmd struct {
	ID string `arg:"" required:"" help:"connection id"`
}

func (cmd GetConnectionCmd) Run(c *cc.CommonCtx) error {
	client, err := tenant.New(c.Context, c.TenantSvcConnectionOptions()...)
	if err != nil {
		return err
	}

	req := &connection.GetConnectionRequest{
		Id: cmd.ID,
	}

	resp, err := client.Connections.GetConnection(c.Context, req)
	if err != nil {
		return errors.Wrapf(err, "get connection [%s]", cmd.ID)
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}

type VerifyConnectionCmd struct {
	ID string `arg:"" required:"" help:"connection id"`
}

func (cmd VerifyConnectionCmd) Run(c *cc.CommonCtx) error {
	client, err := tenant.New(c.Context, c.TenantSvcConnectionOptions()...)
	if err != nil {
		return err
	}

	req := &connection.VerifyConnectionRequest{
		Id: cmd.ID,
	}

	if _, err = client.Connections.VerifyConnection(c.Context, req); err != nil {
		st := status.Convert(err)
		re := regexp.MustCompile(`\r?\n`)

		fmt.Fprintf(c.ErrWriter, "verification    : failed\n")
		fmt.Fprintf(c.ErrWriter, "code            : %d\n", st.Code())
		fmt.Fprintf(c.ErrWriter, "message         : %s\n",
			re.ReplaceAllString(st.Message(), " | "))
		fmt.Fprintf(c.ErrWriter, "error           : %s\n",
			re.ReplaceAllString(st.Err().Error(), " | "))

		for _, detail := range st.Details() {
			if t, ok := detail.(*errdetails.ErrorInfo); ok {
				fmt.Fprintf(c.ErrWriter, "domain          : %s\n", t.Domain)
				fmt.Fprintf(c.ErrWriter, "reason          : %s\n", t.Reason)

				for k, v := range t.Metadata {
					fmt.Fprintf(c.ErrWriter, "detail          : %s (%s)\n", v, k)
				}
			}
		}
	} else {
		fmt.Fprintf(c.ErrWriter, "verification: succeeded\n")
	}

	return nil
}

type UpdateConnectionCmd struct {
	ID string `arg:"" required:"" help:"connection id"`

	Name        string            `optional:"" help:"connection name"`
	Description string            `optional:"" help:"connection description"`
	Kind        string            `optional:"" help:"connection kind: use 'tenant list-provider-kinds' for list of allowed values"`
	ProviderID  string            `optional:"" help:"id of the provider used by the connection"`
	Config      map[string]string `optional:"" help:"connection config values (--config key1=val1 --config key2=val2 ...)"`
}

var (
	ErrUnsupportedConfigType = errors.New("found unsupported type")
	ErrUnknownConfigOption   = errors.New("unknown config option")
	ErrInvalidProviderKind   = errors.New("invalid provider kind")
)

func (cmd *UpdateConnectionCmd) Run(c *cc.CommonCtx) error {
	client, err := tenant.New(c.Context, c.TenantSvcConnectionOptions()...)
	if err != nil {
		return err
	}

	getResponse, err := client.Connections.GetConnection(c.Context, &connection.GetConnectionRequest{Id: cmd.ID})
	if err != nil {
		return errors.Wrapf(err, "get connection [%s]", cmd.ID)
	}

	conn := getResponse.Result

	if cmd.Name != "" {
		conn.Name = cmd.Name
	}
	if cmd.Description != "" {
		conn.Description = cmd.Description
	}
	if cmd.Kind != "" {
		value, ok := api.ProviderKind_value[cmd.Kind]
		if !ok {
			return errors.Wrap(ErrInvalidProviderKind, cmd.Kind)
		}

		conn.Kind = api.ProviderKind(value)
	}
	if cmd.ProviderID != "" {
		conn.ProviderId = cmd.ProviderID
	}

	if err := applyConfigOverrides(conn.Config.Fields, cmd.Config); err != nil {
		return err
	}

	if _, err := client.Connections.UpdateConnection(
		c.Context,
		&connection.UpdateConnectionRequest{Connection: conn},
	); err != nil {
		return errors.Wrap(err, "update connection")
	}

	return jsonx.OutputJSONPB(c.OutWriter, conn)
}

func applyConfigOverrides(config map[string]*structpb.Value, overrides map[string]string) error {
	for key, override := range overrides {

		value, ok := config[key]
		if !ok {
			return errors.Wrap(ErrUnknownConfigOption, key)
		}

		switch value.GetKind().(type) {
		case *structpb.Value_StringValue:
			config[key] = structpb.NewStringValue(override)

		case *structpb.Value_NumberValue:
			asFloat, err := strconv.ParseFloat(override, 64)
			if err != nil {
				return typeMismatch(err, key, "number", override)
			}

			config[key] = structpb.NewNumberValue(asFloat)

		case *structpb.Value_BoolValue:
			asBool, err := strconv.ParseBool(override)
			if err != nil {
				return typeMismatch(err, key, "bool", override)
			}

			config[key] = structpb.NewBoolValue(asBool)

		default:
			return errors.Wrapf(ErrUnsupportedConfigType, "key '%s': type '%T', value: '%v'", key, override, override)
		}
	}

	return nil
}

func typeMismatch(err error, key, expectedType, actualValue string) error {
	return errors.Wrapf(err,
		"type mismatch: config value '%s': expected type '%s': received: '%s'",
		key,
		expectedType,
		actualValue,
	)
}

type SyncConnectionCmd struct {
	ID string `arg:"" required:"" help:"connection id"`
}

func (cmd SyncConnectionCmd) Run(c *cc.CommonCtx) error {
	client, err := tenant.New(c.Context, c.TenantSvcConnectionOptions()...)
	if err != nil {
		return err
	}

	getReq := &connection.GetConnectionRequest{
		Id: cmd.ID,
	}

	curConn, err := client.Connections.GetConnection(c.Context, getReq)
	if err != nil {
		return errors.Wrapf(err, "get connection [%s]", cmd.ID)
	}

	if curConn.Result.Kind != api.ProviderKind_PROVIDER_KIND_IDP {
		return errors.Errorf("connection must be of kind IDP (provided %s)", curConn.Result.Kind.Enum().String())
	}

	updReq := &connection.UpdateConnectionRequest{
		Connection: curConn.Result,
		Force:      false,
	}

	if _, err = client.Connections.UpdateConnection(c.Context, updReq); err != nil {
		return errors.Wrapf(err, "update connection [%s]", cmd.ID)
	}

	return nil
}
