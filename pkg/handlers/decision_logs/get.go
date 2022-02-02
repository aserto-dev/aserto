package decision_logs //nolint // prefer standardizing name over removing _

import (
	"context"
	"io/fs"

	"github.com/aserto-dev/aserto/pkg/cc"
	dl "github.com/aserto-dev/go-grpc/aserto/decision_logs/v1"
	"google.golang.org/protobuf/proto"
)

type GetCmd struct {
	Name     string   `arg:"" optional:"" help:"download decision logs"`
	Path     string   `optional:"" help:"download path"`
	Info     bool     `xor:"group2" optional:"" help:"get information about the logs, don't download"`
	Policies []string `optional:"" sep:"," help:"ID of policies to get logs for (all if not specified)"`
}

func (cmd GetCmd) Run(c *cc.CommonCtx, apiKey APIKey) error {
	impl := getImpl{
		c:         c,
		id:        cmd.Name,
		info:      cmd.Info,
		localPath: cmd.Path,
		apiKey:    apiKey,
		getter:    &cmd,
		dirPaths:  cmd.Policies,
	}

	return impl.run()
}

func (cmd *GetCmd) list(ctx context.Context, cli dl.DecisionLogsClient, paths []string) ([]proto.Message, error) {
	return listDecisionLogs(ctx, cli, paths)
}

func (cmd *GetCmd) get(ctx context.Context, cli dl.DecisionLogsClient, id string) (proto.Message, error) {
	resp, err := cli.GetDecisionLog(ctx, &dl.GetDecisionLogRequest{
		Name: id,
	})

	if err != nil {
		return nil, err
	}

	return resp.Log, nil
}

func (cmd *GetCmd) idFromListItem(item proto.Message) string {
	listItem := item.(*dl.DecisionLogItem)
	return listItem.Name
}

func (cmd *GetCmd) urlFromItem(item proto.Message) string {
	log := item.(*dl.DecisionLog)
	return log.Url
}

func (cmd *GetCmd) shouldFetch(finfo fs.FileInfo, item proto.Message) bool {
	return false
}
