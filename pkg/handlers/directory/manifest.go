package directory

import (
	"io"
	"os"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/go-directory/aserto/directory/model/v3"
	"github.com/pkg/errors"
)

type GetManifestMetadataCmd struct{}

func (cmd *GetManifestMetadataCmd) Run(c *cc.CommonCtx) error {
	client, err := c.DirectoryClient()
	if err != nil {
		return errors.Wrap(err, "failed to get directory client")
	}

	respStream, err := client.Model.GetManifest(c.Context, &model.GetManifestRequest{})
	if err != nil {
		return errors.Wrap(err, "get manifest call failed")
	}

	var metadata *model.Metadata
	for {
		manifestResp, err := respStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return errors.Wrap(err, "receive message from get manifest stream failed")
		}

		if manifestResp.GetBody() != nil {
			_ = manifestResp.GetBody()
		} else if manifestResp.GetMetadata() != nil {
			metadata = manifestResp.GetMetadata()
		}
	}

	if err := respStream.CloseSend(); err != nil {
		return errors.Wrap(err, "failed to close stream")
	}

	return jsonx.OutputJSONPB(c.UI.Output(), metadata)
}

type GetManifestCmd struct{}

func (cmd *GetManifestCmd) Run(c *cc.CommonCtx) error {
	client, err := c.DirectoryClient()
	if err != nil {
		return errors.Wrap(err, "failed to get directory client")
	}

	respStream, err := client.Model.GetManifest(c.Context, &model.GetManifestRequest{})
	if err != nil {
		return errors.Wrap(err, "get manifest call failed")
	}

	body := make([]byte, 0)
	for {
		manifestResp, err := respStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return errors.Wrap(err, "receive message from get manifest stream failed")
		}

		if manifestResp.GetBody() != nil {
			body = append(body, manifestResp.GetBody().Data...)
		} else if manifestResp.GetMetadata() != nil {
			_ = manifestResp.GetMetadata()
		}
	}

	if err := respStream.CloseSend(); err != nil {
		return errors.Wrap(err, "failed to close stream")
	}

	c.UI.Normal().Msg(string(body))
	return nil
}

type SetManifestCmd struct {
	ManifestFilePath string `arg:"" name:"manifest_file_path" required:"" help:"absolute path to the manifest file"`
}

func (cmd *SetManifestCmd) Run(c *cc.CommonCtx) error {
	client, err := c.DirectoryClient()
	if err != nil {
		return errors.Wrap(err, "failed to get directory client")
	}

	bytesContent, err := os.ReadFile(cmd.ManifestFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to read manifest file")
	}

	setManifestStream, err := client.Model.SetManifest(c.Context)
	if err != nil {
		return errors.Wrap(err, "set manifest call failed")
	}

	if err := setManifestStream.Send(
		&model.SetManifestRequest{
			Msg: &model.SetManifestRequest_Body{
				Body: &model.Body{Data: bytesContent},
			}}); err != nil {
		return errors.Wrap(err, "failed to send manifest body")
	}

	resp, err := setManifestStream.CloseAndRecv()
	if err != nil {
		return errors.Wrap(err, "failed to close stream")
	}
	return jsonx.OutputJSONPB(c.UI.Output(), resp)
}

type DeleteManifestCmd struct{}

func (cmd *DeleteManifestCmd) Run(c *cc.CommonCtx) error {
	client, err := c.DirectoryClient()
	if err != nil {
		return errors.Wrap(err, "failed to get directory client")
	}

	resp, err := client.Model.DeleteManifest(c.Context, &model.DeleteManifestRequest{})
	if err != nil {
		return errors.Wrap(err, "delete manifest call failed")
	}

	return jsonx.OutputJSONPB(c.UI.Output(), resp)
}
