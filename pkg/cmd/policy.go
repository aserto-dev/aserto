package cmd

import "github.com/aserto-dev/aserto/pkg/handlers/policy"

type PolicyCmd struct {
	Build     policy.BuildCmd `cmd:"" help:"Build policies." group:"policy"`
	List      policy.CmdCmd   `cmd:"" help:"List policy images." group:"policy"`
	Pull      policy.PullCmd  `cmd:"" help:"Push policies to a registry." group:"policy"`
	Push      policy.PushCmd  `cmd:"" help:"Pull policies from a registry." group:"policy"`
	Login     policy.CmdCmd   `cmd:"" help:"Login to a registry." group:"policy"`
	Logout    policy.CmdCmd   `cmd:"" help:"Logout from a registry." group:"policy"`
	Save      policy.CmdCmd   `cmd:"" help:"Save a policy to a local bundle tarball." group:"policy"`
	Tag       policy.CmdCmd   `cmd:"" help:"Create a new tag for an existing policy." group:"policy"`
	Rm        policy.CmdCmd   `cmd:"" help:"Removes a policy from the local registry." group:"policy"`
	Inspect   policy.CmdCmd   `cmd:"" help:"Displays information about a policy." group:"policy"`
	Repl      policy.CmdCmd   `cmd:"" help:"Sets you up with a shell for running queries using an OPA instance with a policy loaded." group:"policy"`
	Templates policy.CmdCmd   `cmd:"" help:"List and apply templates" group:"policy"`
	Version   policy.CmdCmd   `cmd:"" help:"Prints version information." group:"policy"`
}
