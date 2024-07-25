package controlplane

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/go-grpc/aserto/management/v2"
	"github.com/aserto-dev/topaz/pkg/cli/jsonx"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type ListInstanceRegistrationsCmd struct{}

func (cmd ListInstanceRegistrationsCmd) Run(c *cc.CommonCtx) error {
	cli, err := c.ControlPlaneClient(c.Context)
	if err != nil {
		return err
	}

	resp, err := cli.ListInstanceRegistrations(c.Context, &management.ListInstanceRegistrationsRequest{})
	if err != nil {
		return err
	}

	var results []protoreflect.ProtoMessage
	for _, inst := range resp.Result {
		results = append(results, inst)
	}

	err = jsonx.OutputJSONPBArray(c.StdOut(), results)
	if err != nil {
		return err
	}

	return nil
}
