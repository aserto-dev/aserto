package directory

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/aserto/pkg/pb"
	dir "github.com/aserto-dev/go-grpc/aserto/authorizer/directory/v1"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/structpb"
)

type GetResCmd struct {
	Key string `arg:"key" name:"key" required:"" help:"resource key"`
}

func (cmd *GetResCmd) Run(c *cc.CommonCtx) error {
	client, err := c.AuthorizerClient()
	if err != nil {
		return err
	}

	resp, err := client.Directory.GetResource(c.Context, &dir.GetResourceRequest{
		Key: cmd.Key,
	})
	if err != nil {
		return err
	}

	return jsonx.OutputJSON(c.OutWriter, resp.Value)
}

type SetResCmd struct {
	Key   string `arg:"key" name:"key" required:"" help:"resource key"`
	Value string `optional:"" help:"set resource using string value"`
	Stdin bool   `optional:"" name:"stdin" help:"set resource using from --stdin"`
	File  string `optional:"" type:"existingfile" help:"set resource using file content"`
}

func (cmd *SetResCmd) Run(c *cc.CommonCtx) error {
	var (
		value structpb.Struct
		buf   io.Reader
		err   error
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
		fmt.Fprintf(c.ErrWriter, "reading value flag\n")
		buf = strings.NewReader(cmd.Value)
	default:
		return errors.Errorf("no input option specified [--stdin | --file <filepath> | --value <string>]")
	}

	if buf == nil {
		value = structpb.Struct{}
	} else if err = pb.BufToProto(buf, &value); err != nil { //nolint:gocritic
		return err
	}

	client, err := c.AuthorizerClient()
	if err != nil {
		return err
	}

	resp, err := client.Directory.SetResource(c.Context, &dir.SetResourceRequest{
		Key:   cmd.Key,
		Value: &value,
	})
	if err != nil {
		return err
	}

	return jsonx.OutputJSON(c.OutWriter, resp)
}

type DelResCmd struct {
	Key string `arg:"key" name:"key" required:"" help:"resource key"`
}

func (cmd *DelResCmd) Run(c *cc.CommonCtx) error {
	client, err := c.AuthorizerClient()
	if err != nil {
		return err
	}

	_, err = client.Directory.DeleteResource(c.Context, &dir.DeleteResourceRequest{
		Key: cmd.Key,
	})

	if err != nil {
		return err
	}

	return nil
}

type ListResCmd struct{}

func (cmd *ListResCmd) Run(c *cc.CommonCtx) error {
	client, err := c.AuthorizerClient()
	if err != nil {
		return err
	}

	resp, err := client.Directory.ListResources(c.Context, &dir.ListResourcesRequest{})
	if err != nil {
		return err
	}

	return jsonx.OutputJSON(c.OutWriter, resp)
}
