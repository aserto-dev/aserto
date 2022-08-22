package directory

import (
	"io"
	"os"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/aserto/pkg/pb"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
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

	return jsonx.OutputJSON(c.UI.Output(), resp.Value)
}

type SetResCmd struct {
	Key   string         `arg:"key" name:"key" required:"" help:"resource key"`
	Value structpb.Value `xor:"group" required:"" name:"value" help:"set resource value using json data from argument"`
	Stdin bool           `xor:"group" required:"" name:"stdin" help:"set resource value using json data from --stdin"`
	File  string         `xor:"group" required:"" name:"file" type:"existingfile" help:"set resource value using json data file"`
}

func (cmd *SetResCmd) Run(c *cc.CommonCtx) error {
	var (
		value *structpb.Value
		buf   io.Reader
		err   error
	)

	switch {
	case cmd.Stdin:
		buf = os.Stdin

		value, err = pb.BufToValue(buf)
		if err != nil {
			return errors.Wrapf(err, "unmarshal stdin")
		}

	case cmd.File != "":
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

	structValue := pb.NewStruct()
	structValue.Fields[cmd.Key] = value

	client, err := c.AuthorizerClient()
	if err != nil {
		return err
	}

	_, err = client.Directory.SetResource(c.Context, &dir.SetResourceRequest{
		Key:   cmd.Key,
		Value: structValue,
	})
	if err != nil {
		return err
	}

	return nil
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

type ListResCmd struct {
	Count bool `name:"count" optional:"" help:"only return resource count"`
}

func (cmd *ListResCmd) Run(c *cc.CommonCtx) error {
	client, err := c.AuthorizerClient()
	if err != nil {
		return err
	}

	token := ""
	pageSize := int32(100)
	var apiRes []string

	if cmd.Count {
		pageSize = int32(-2)
	}

	for {
		resp, err := client.Directory.ListResources(c.Context, &dir.ListResourcesRequest{
			Page: &api.PaginationRequest{
				Size:  pageSize,
				Token: token,
			},
		})
		if err != nil {
			return err
		}

		if cmd.Count {
			return jsonx.OutputJSONPB(c.UI.Output(), resp.Page)
		}

		apiRes = append(apiRes, resp.Results...)
		if resp.Page.NextToken == "" {
			break
		}

		token = resp.Page.NextToken
	}

	return jsonx.OutputJSON(c.UI.Output(), apiRes)
}
