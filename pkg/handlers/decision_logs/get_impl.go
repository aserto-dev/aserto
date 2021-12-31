package decision_logs //nolint // prefer standardizing name over removing _

import (
	"context"
	"io/fs"
	"os"
	"path"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	dl "github.com/aserto-dev/go-grpc/aserto/decision_logs/v1"
	"google.golang.org/protobuf/proto"
)

type getter interface {
	list(context.Context, dl.DecisionLogsClient) ([]proto.Message, error)
	get(context.Context, dl.DecisionLogsClient, string) (proto.Message, error)
	idFromListItem(proto.Message) string
	urlFromItem(proto.Message) string
	shouldFetch(fs.FileInfo, proto.Message) bool
}

type getImpl struct {
	c         *cc.CommonCtx
	id        string
	info      bool
	localPath string
	apiKey    APIKey
	getter    getter
}

func (impl *getImpl) run() error {
	ctx := impl.c.Context
	cli, err := newClient(impl.c, impl.apiKey)
	if err != nil {
		return err
	}

	var ids []string

	if impl.id == "" {
		ids, err = impl.calculateGet(ctx, cli)
		if err != nil {
			return err
		}
	} else {
		ids = []string{impl.id}
	}

	users, err := impl.get(ctx, cli, ids)
	if err != nil {
		return err
	}

	if impl.info {
		return jsonx.OutputJSONPBMap(impl.c.OutWriter, users)
	}

	for id, msg := range users {
		url := impl.getter.urlFromItem(msg)
		err := download(ctx, id, url, impl.localPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (impl *getImpl) calculateGet(ctx context.Context, cli dl.DecisionLogsClient) ([]string, error) {
	localPath := impl.localPath
	if !path.IsAbs(localPath) {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		localPath = path.Join(wd, localPath)
	}

	logFS := os.DirFS(localPath)
	oldItems := map[string]fs.FileInfo{}
	err := fs.WalkDir(logFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() != "." {
				return fs.SkipDir
			}

			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		oldItems[path] = info
		return nil
	})
	if err != nil {
		return nil, err
	}

	items, err := impl.getter.list(ctx, cli)
	if err != nil {
		return nil, err
	}

	newItems := []string{}

	for _, msg := range items {
		id := impl.getter.idFromListItem(msg)
		local, ok := oldItems[id]
		if ok && !impl.getter.shouldFetch(local, msg) {
			continue
		}

		newItems = append(newItems, id)
	}

	return newItems, nil
}

func (impl *getImpl) get(ctx context.Context, cli dl.DecisionLogsClient, ids []string) (map[string]proto.Message, error) {
	items := map[string]proto.Message{}

	for _, id := range ids {
		item, err := impl.getter.get(ctx, cli, id)
		if err != nil {
			return nil, err
		}
		items[id] = item
	}

	return items, nil
}
