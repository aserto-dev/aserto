package decision_logs //nolint // prefer standardizing name over removing _

import (
	"context"
	"io/fs"

	"github.com/aserto-dev/aserto/pkg/cc"
	dl "github.com/aserto-dev/go-grpc/aserto/decision_logs/v1"
	"google.golang.org/protobuf/proto"
)

type GetUserCmd struct {
	ID   string `xor:"group" optional:"" help:"download decision logs user information"`
	Path string `arg:"" optional:"" help:"download path"`
	Info bool   `xor:"group2" optional:"" help:"get information about the logs, don't download"`
}

func (cmd GetUserCmd) Run(c *cc.CommonCtx, apiKey APIKey) error {
	impl := getImpl{
		c:         c,
		id:        cmd.ID,
		info:      cmd.Info,
		localPath: cmd.Path,
		apiKey:    apiKey,
		getter:    &cmd,
	}

	return impl.run()
}

func (cmd *GetUserCmd) list(ctx context.Context, cli dl.DecisionLogsClient, paths []string) ([]proto.Message, error) {
	return listUsers(ctx, cli)
}

func (cmd *GetUserCmd) get(ctx context.Context, cli dl.DecisionLogsClient, id string) (proto.Message, error) {
	resp, err := cli.GetUser(ctx, &dl.GetUserRequest{
		Id: id,
	})

	if err != nil {
		return nil, err
	}

	return resp.User, nil
}

func (cmd *GetUserCmd) idFromListItem(item proto.Message) string {
	listItem := item.(*dl.UserItem)
	return listItem.Id
}

func (cmd *GetUserCmd) urlFromItem(item proto.Message) string {
	log := item.(*dl.User)
	return log.Url
}

func (cmd *GetUserCmd) shouldFetch(finfo fs.FileInfo, item proto.Message) bool {
	userItem := item.(*dl.UserItem)
	return userItem.UpdatedAt.AsTime().After(finfo.ModTime())
}
