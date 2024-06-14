package cmd

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/handlers/dev"
	topazConfig "github.com/aserto-dev/topaz/pkg/cli/cmd/configure"
	topaz "github.com/aserto-dev/topaz/pkg/cli/cmd/topaz"
)

type TopazCmd struct {
	RunCmd    topaz.RunCmd       `cmd:"" name:"run" group:"topaz" help:"run topaz instance in interactive mode"`
	Start     topaz.StartCmd     `cmd:"" group:"topaz" help:"start a topaz instance"`
	Stop      topaz.StopCmd      `cmd:"" group:"topaz" help:"stop a topaz instance"`
	Status    topaz.StatusCmd    `cmd:"" group:"topaz" help:"status of topaz instance"`
	Update    topaz.UpdateCmd    `cmd:"" group:"topaz" help:"download the latest aserto topaz image"`
	Console   topaz.ConsoleCmd   `cmd:"" group:"topaz" help:"launch web console"`
	Config    AsertoConfigCmd    `cmd:"" group:"topaz" help:"configure a policy"`
	Install   topaz.InstallCmd   `cmd:"" group:"topaz" help:"install topaz"`
	Uninstall topaz.UninstallCmd `cmd:"" group:"topaz" help:"uninstall topaz, removes all locally installed artifacts"`
}

type AsertoConfigCmd struct {
	New    dev.ConfigureCmd            `cmd:"" help:"create new configuration"`
	Rename topazConfig.RenameConfigCmd `cmd:"" help:"rename configuration"`
	Delete topazConfig.DeleteConfigCmd `cmd:"" help:"delete configuration"`
}

func (cmd *TopazCmd) Run(c *cc.CommonCtx) error {
	return nil
}
