package user

import (
	"os"
	"path/filepath"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/cc/config"
	"github.com/aserto-dev/aserto/pkg/keyring"
	"github.com/aserto-dev/aserto/pkg/x"
	"github.com/pkg/errors"
)

type LogoutCmd struct {
}

func (cmd *LogoutCmd) Run(c *cc.CommonCtx) error {
	kr, err := keyring.NewKeyRing(c.Auth.Issuer)
	if err != nil {
		return errors.Wrapf(err, "instantiate keyring")
	}

	err = kr.DelToken()
	if err != nil {
		return errors.Wrapf(err, "delete token")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return errors.Wrap(err, "failed to determine user home directory")
	}

	filePath := filepath.Join(home, ".config", x.AppName, config.ConfigPath)

	return os.Remove(filePath)
}
