package cmd

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/handlers/config"
	"github.com/aserto-dev/aserto/pkg/handlers/dev"
	topazConfig "github.com/aserto-dev/topaz/pkg/cli/cmd/configure"
)

type ConfigCmd struct {
	Use    config.UseConfigCmd         `cmd:"" help:"use a topaz configuration"`
	New    dev.ConfigureCmd            `cmd:"" help:"create new configuration"`
	List   config.ListConfigCmd        `cmd:"" help:"list configurations"`
	Rename topazConfig.RenameConfigCmd `cmd:"" help:"rename configuration"`
	Delete topazConfig.DeleteConfigCmd `cmd:"" help:"delete configuration"`
	Info   topazConfig.InfoConfigCmd   `cmd:"" help:"display configuration information"`
}

func (cmd *ConfigCmd) Run(c *cc.CommonCtx) error {
	return nil
}
