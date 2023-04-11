package cmd

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/handlers/dev"
)

type DeveloperCmd struct {
	Start             dev.StartCmd          `cmd:"" group:"developer" help:"start sidecar instance"`
	Stop              dev.StopCmd           `cmd:"" group:"developer" help:"stop sidecar instance"`
	Status            dev.StatusCmd         `cmd:"" group:"developer" help:"status of sidecar instance"`
	Update            dev.UpdateCmd         `cmd:"" group:"developer" help:"download the latest aserto onebox image"`
	Console           dev.ConsoleCmd        `cmd:"" group:"developer" help:"launch web console"`
	Configure         dev.ConfigureCmd      `cmd:"" group:"developer" help:"configure a policy"`
	PolicyFromOpenAPI dev.PolicyFromOpenAPI `cmd:"" group:"developer" help:"generate an open api policy"`
	Install           dev.InstallCmd        `cmd:"" group:"developer" help:"install aserto onebox"`
	Uninstall         dev.UninstallCmd      `cmd:"" group:"developer" help:"uninstall aserto onebox, removes all locally installed artifacts"`
}

func (cmd *DeveloperCmd) Run(c *cc.CommonCtx) error {
	return nil
}
