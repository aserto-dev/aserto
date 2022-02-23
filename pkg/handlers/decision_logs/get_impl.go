package decision_logs //nolint // prefer standardizing name over removing _

import (
	"context"
	"io/fs"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	dl "github.com/aserto-dev/go-grpc/aserto/decision_logs/v1"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

const concurrency = 10

type getter interface {
	list(context.Context, dl.DecisionLogsClient, []string) ([]proto.Message, error)
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
	getter    getter
	dirPaths  []string
}

type item struct {
	id  string
	msg proto.Message
}

func (impl *getImpl) run() error {
	ctx := impl.c.Context
	cli, err := impl.c.DecisionLogsClient()
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

	items, err := impl.get(ctx, cli, ids)
	if err != nil {
		return err
	}

	if impl.info {
		return jsonx.OutputJSONPBMap(impl.c.UI.Output(), items)
	}

	itemCh := make(chan item)
	errCh := make(chan error, len(items))
	done := sync.WaitGroup{}

	done.Add(concurrency)
	for worker := 0; worker < concurrency; worker++ {
		go func() {
			defer done.Done()

			for itm := range itemCh {
				url := impl.getter.urlFromItem(itm.msg)
				dlerr := download(ctx, itm.id, url, impl.localPath)
				if dlerr != nil {
					errCh <- errors.Wrapf(dlerr, "error downloading '%s'", itm.id)
				}
			}
		}()
	}

	for id, msg := range items {
		itemCh <- item{
			id:  id,
			msg: msg,
		}
	}
	close(itemCh)
	done.Wait()

	for done := false; !done; {
		select {
		case dlerr := <-errCh:
			err = multierror.Append(err, dlerr)
		default:
			done = true
		}
	}
	close(errCh)

	if err != nil {
		return err
	}

	return nil
}

func (impl *getImpl) isInPaths(p string) bool {
	if impl.dirPaths == nil {
		return true
	}
	for _, v := range impl.dirPaths {
		if strings.HasPrefix(v, p) {
			return true
		}
	}

	return false
}

func (impl *getImpl) getOldItems() (map[string]fs.FileInfo, error) {
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

	walkFunc := func(fileDirPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "." {
				return nil
			}

			newPath := path.Join(fileDirPath, d.Name())
			if !impl.isInPaths(newPath) {
				return fs.SkipDir
			}

			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		oldItems[fileDirPath] = info
		return nil
	}

	err := fs.WalkDir(logFS, ".", walkFunc)
	if err != nil {
		return nil, err
	}

	return oldItems, nil
}

func (impl *getImpl) calculateGet(ctx context.Context, cli dl.DecisionLogsClient) ([]string, error) {
	var oldItems map[string]fs.FileInfo
	var err error

	if !impl.info {
		oldItems, err = impl.getOldItems()
		if err != nil {
			return nil, err
		}
	}

	items, err := impl.getter.list(ctx, cli, impl.dirPaths)
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
