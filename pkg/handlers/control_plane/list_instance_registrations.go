package controlplane

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/go-grpc/aserto/management/v2"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ListInstanceRegistrationsCmd struct{}

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

	err = jsonx.OutputJSONPBArray(c.TopazContext.StdOut(), instsOut)
	if err != nil {
		return err
	}

	return nil
}
