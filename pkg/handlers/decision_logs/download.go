package decision_logs //nolint // prefer standardizing name over removing _

import (
	"context"
	"net/http"
	"os"
	"path"
	"strings"
)

func download(ctx context.Context, name, url, downloadDir string) error {
	var err error

	if downloadDir == "" {
		downloadDir, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	httpClient := http.Client{}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return err
	}

	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	tmpFile, err := os.CreateTemp("", "aserto-")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.ReadFrom(httpResp.Body)
	if err != nil {
		return err
	}

	lastSep := strings.LastIndex(name, "/")
	if lastSep != -1 {
		subDir := path.Join(downloadDir, name[:lastSep])
		err = os.MkdirAll(subDir, 0700)
		if err != nil {
			return err
		}
	}

	finalPath := path.Join(downloadDir, name)
	err = os.Rename(tmpFile.Name(), finalPath)

	return err
}
