package directory

import (
	"encoding/json"
	"os"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/clui"
	"github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
)

type GetGraphCmd struct {
	Request  string `arg:""  type:"existingfile" name:"request" optional:"" help:"file path to get graph request or '-' to read from stdin"`
	Template bool   `name:"template" help:"prints a get graph request template on stdout"`
}

func (cmd *GetGraphCmd) Run(c *cc.CommonCtx) error {
	if cmd.Template {
		return printGetGraphRequest(c.UI)
	}

	client, err := c.DirectoryClient()
	if err != nil {
		return errors.Wrap(err, "failed to get directory client")
	}

	if cmd.Request == "" {
		return errors.New("request argument is required")
	}

	var req reader.GetGraphRequest
	if cmd.Request == "-" {
		decoder := json.NewDecoder(os.Stdin)

		err = decoder.Decode(&req)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal request from stdin")
		}
	} else {
		dat, err := os.ReadFile(cmd.Request)
		if err != nil {
			return errors.Wrapf(err, "opening file [%s]", cmd.Request)
		}

		err = protojson.Unmarshal(dat, &req)
		if err != nil {
			return errors.Wrapf(err, "failed to unmarshal request from file [%s]", cmd.Request)
		}
	}

	resp, err := client.Reader.GetGraph(c.Context, &req)
	if err != nil {
		return errors.Wrap(err, "get graph call failed")
	}

	return jsonx.OutputJSONPB(c.UI.Output(), resp)
}

func printGetGraphRequest(ui *clui.UI) error {
	req := &reader.GetGraphRequest{
		AnchorType:      "",
		AnchorId:        "",
		ObjectType:      "",
		ObjectId:        "",
		Relation:        "",
		SubjectType:     "",
		SubjectId:       "",
		SubjectRelation: "",
	}
	return jsonx.OutputJSONPB(ui.Output(), req)
}
