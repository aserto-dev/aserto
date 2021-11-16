package directory

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/dirx"
)

const (
	providerJSON  string = "json"
	providerAuth0 string = "auth0"
)

type LoadUserExtCmd struct {
	Provider string `required:"" help:"load users provider (json | auth0)" enum:"json,auth0"`
	Profile  string `optional:"" type:"existingfile" help:"provider profile file (.env)"`
	File     string `optional:"" type:"existingfile" help:"input file (.json)"`
}

func (cmd *LoadUserExtCmd) Run(c *cc.CommonCtx) error {
	loader := UserLoader{Provider: cmd.Provider, Profile: cmd.Profile, File: cmd.File}
	return loader.Load(c, dirx.UserExtensionsRequestFactory)
}
