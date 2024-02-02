package cmd

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/handlers/dev"
	topaz "github.com/aserto-dev/topaz/pkg/cli/cmd"
)

type TopazCmd struct {
	RunCmd            topaz.RunCmd          `cmd:"" name:"run" group:"topaz" help:"run topaz instance in interactive mode"`
	Start             topaz.StartCmd        `cmd:"" group:"topaz" help:"start a topaz instance"`
	Stop              topaz.StopCmd         `cmd:"" group:"topaz" help:"stop a topaz instance"`
	Status            topaz.StatusCmd       `cmd:"" group:"topaz" help:"status of topaz instance"`
	Update            topaz.UpdateCmd       `cmd:"" group:"topaz" help:"download the latest aserto topaz image"`
	Console           topaz.ConsoleCmd      `cmd:"" group:"topaz" help:"launch web console"`
	Configure         topaz.ConfigureCmd    `cmd:"" group:"topaz" help:"configure a policy"`
	List              topaz.ListConfigCmd   `cmd:"" group:"topaz" help:"list topaz configuration files"`
	PolicyFromOpenAPI dev.PolicyFromOpenAPI `cmd:"" group:"topaz" help:"generate an open api policy"`
	Install           topaz.InstallCmd      `cmd:"" group:"topaz" help:"install topaz"`
	Uninstall         topaz.UninstallCmd    `cmd:"" group:"topaz" help:"uninstall topaz, removes all locally installed artifacts"`
}

func (cmd *TopazCmd) Run(c *cc.CommonCtx) error {
	return nil
}
