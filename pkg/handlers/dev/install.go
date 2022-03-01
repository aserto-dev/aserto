package dev

import (
	"context"
	"fmt"
	"os"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/dockerx"
	"github.com/aserto-dev/aserto/pkg/filex"
	"github.com/aserto-dev/aserto/pkg/handlers/dev/certs"
	"github.com/aserto-dev/aserto/pkg/orasx"
	localpaths "github.com/aserto-dev/aserto/pkg/paths"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"oras.land/oras-go/pkg/content"
	"oras.land/oras-go/pkg/oras"
)

type InstallCmd struct {
	TrustCert bool `optional:"" default:"false" help:"add onebox certificate to the system's trusted CAs"`

	ImageName    string `optional:""  default:"authorizer-onebox" help:"image name"`
	ImageVersion string `optional:""  default:"latest" help:"image version" `
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

	paths, err := localpaths.Create()
	if err != nil {
		return errors.Wrap(err, "failed to create configuration directory")
	}

	cfgLocal := paths.LocalConfig()
	if !filex.FileExists(cfgLocal) {
		fmt.Fprintf(c.UI.Output(), "creating %s\n", cfgLocal)
		if err := createLocalConfig(cfgLocal); err != nil {
			return errors.Wrap(err, "create local configuration")
		}
	}

	edsFile := paths.LocalEDS()
	if !filex.FileExists(edsFile) {
		fmt.Fprintf(c.UI.Output(), "creating %s\n", edsFile)
		if err := createDefaultEds(edsFile); err != nil {
			return errors.Wrap(err, "create default eds")
		}
	}

	// Create onebox certs if none exist.
	if err := certs.GenerateCerts(c.UI.Output(), c.UI.Err(), paths.Certs.GRPC, paths.Certs.Gateway); err != nil {
		return errors.Wrap(err, "failed to create dev certificates")
	}

	if cmd.TrustCert {
		fmt.Fprintln(c.UI.Output(), "Adding developer certificate to system store. You may need to provide credentials.")
		if err := certs.AddTrustedCert(paths.Certs.Gateway.CA); err != nil {
			return errors.Wrap(err, "add gateway-ca cert to trusted certs")
		}
	}

	return dockerx.DockerWith(map[string]string{
		"IMAGE_NAME":    cmd.ImageName,
		"IMAGE_VERSION": cmd.ImageVersion,
	},
		"pull", "ghcr.io/aserto-dev/$IMAGE_NAME:$IMAGE_VERSION",
	)
}

func createLocalConfig(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return errors.Wrapf(err, "creating %s", path)
	}

	return WriteConfig(f, configTemplateLocal, &templateParams{TenantID: localTenantID})
}

func createDefaultEds(edsFile string) error {
	ctx := context.Background()

	resolver := orasx.NewResolver("", "", false, false, []string{}...)

	fileStore := content.NewFileStore(edsFile)
	defer fileStore.Close()

	allowedMediaTypes := []string{content.DefaultBlobMediaType, content.DefaultBlobDirMediaType}

	pullOpts := []oras.PullOpt{
		oras.WithAllowedMediaTypes(allowedMediaTypes),
		oras.WithPullStatusTrack(os.Stdout),
	}

	ref := "ghcr.io/aserto-demo/assets/eds:v9"

	_, _, err := oras.Pull(ctx, resolver, ref, fileStore, pullOpts...)
	if err != nil {
		return errors.Wrap(err, "pull assets/eds:v9")
	}

	return nil
}
