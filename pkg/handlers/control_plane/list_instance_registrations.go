package controlplane

import (
	"github.com/aserto-dev/go-grpc/management/v2"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ListInstanceRegistrationsCmd struct {
}

func (cmd ListInstanceRegistrationsCmd) Run(c *cc.CommonCtx) error {
	cli, err := c.ControlPlaneClient()
	if err != nil {
		return err
	}

	resp, err := cli.ListInstanceRegistrations(c.Context, &management.ListInstanceRegistrationsRequest{})
	if err != nil {
		return err
	}

	var instsOut []protoreflect.ProtoMessage
	for _, inst := range resp.Result {
		instsOut = append(instsOut, inst)
	}

	err = jsonx.OutputJSONPBArray(c.UI.Output(), instsOut)
	if err != nil {
		return err
	}

	return nil
}
