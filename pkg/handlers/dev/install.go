package dev

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/dockerx"
	"github.com/aserto-dev/aserto/pkg/filex"
	"github.com/aserto-dev/aserto/pkg/orasx"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"oras.land/oras-go/pkg/content"
	"oras.land/oras-go/pkg/oras"
)

type InstallCmd struct {
	ContainerName    string `optional:""  default:"aserto-one" help:"container name"`
	ContainerVersion string `optional:""  default:"latest" help:"container version" `
}

func (cmd InstallCmd) Run(c *cc.CommonCtx) error {
	if running, err := dockerx.IsRunning(dockerx.AsertoOne); running || err != nil {
		if err != nil {
			return err
		}
		color.Yellow("!!! onebox is already running")
		return nil
	}

	color.Green(">>> installing onebox...")

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := path.Join(home, "/.config/aserto/aserto-one/cfg")
	if !filex.DirExists(configDir) {
		if err := os.MkdirAll(configDir, 0700); err != nil {
			return err
		}
	}

	cfgLocal := path.Join(configDir, "/local.yaml")
	if !filex.FileExists(cfgLocal) {
		fmt.Fprintf(c.OutWriter, "creating %s\n", cfgLocal)
		if err := ioutil.WriteFile(cfgLocal, []byte(configTemplateLocal), 0600); err != nil {
			return errors.Wrapf(err, "writing %s", cfgLocal)
		}
	}

	cacheDir := path.Join(home, "/.cache/aserto/aserto-one/eds")
	if !filex.DirExists(configDir) {
		if err := os.MkdirAll(configDir, 0700); err != nil {
			return err
		}
	}

	edsFile := path.Join(cacheDir, "/eds-acmecorp-v4.db")
	if !filex.FileExists(edsFile) {
		fmt.Fprintf(c.OutWriter, "creating %s\n", edsFile)
		ctx := context.Background()

		resolver := orasx.NewResolver("", "", false, false, []string{}...)

		fileStore := content.NewFileStore(cacheDir)
		defer fileStore.Close()

		allowedMediaTypes := []string{content.DefaultBlobMediaType, content.DefaultBlobDirMediaType}

		pullOpts := []oras.PullOpt{
			oras.WithAllowedMediaTypes(allowedMediaTypes),
			oras.WithPullStatusTrack(os.Stdout),
		}

		ref := "ghcr.io/aserto-demo/assets/eds:v4"

		_, _, err := oras.Pull(ctx, resolver, ref, fileStore, pullOpts...)
		if err != nil {
			return errors.Wrap(err, "pull assets/eds:v4")
		}
	}

	return dockerx.DockerWith(map[string]string{
		"CONTAINER_NAME":    cmd.ContainerName,
		"CONTAINER_VERSION": cmd.ContainerVersion,
	},
		"pull", "ghcr.io/aserto-dev/$CONTAINER_NAME:$CONTAINER_VERSION",
	)
}

const configTemplateLocal = `
---
logging:
  prod: false
  log_level: debug

directory_service:
  path: "/app/eds/eds-acmecorp-v4.db"

api:
  grpc:
    connection_timeout_seconds: 2
    certs:
      tls_key_path: "/root/.config/aserto/aserto-one/certs/grpc.key"
      tls_cert_path: "/root/.config/aserto/aserto-one/certs/grpc.crt"
      tls_ca_cert_path: "/root/.config/aserto/aserto-one/certs/grpc-ca.crt"
  gateway:
    certs:
      tls_key_path: "/root/.config/aserto/aserto-one/certs/gateway.key"
      tls_cert_path: "/root/.config/aserto/aserto-one/certs/gateway.crt"
      tls_ca_cert_path: "/root/.config/aserto/aserto-one/certs/gateway-ca.crt"

opa:
  instance_id: "0fb5d7eb-8190-4f9d-ac7f-db0ba8374cb7"
  store: aserto
  graceful_shutdown_period_seconds: 2
  local_bundles:
    paths: []
    skip_verification: true
`
