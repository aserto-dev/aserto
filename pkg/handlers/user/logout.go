package user

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/keyring"
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

	return nil
}
