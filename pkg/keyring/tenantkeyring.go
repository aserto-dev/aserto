package keyring

import (
	"encoding/json"
	"log"
	"os/user"

	"github.com/aserto-dev/aserto/pkg/auth0/api"
	"github.com/pkg/errors"
	"github.com/zalando/go-keyring"
)

type TenantKeyRing struct {
	tenantID string
	user     string
}

func NewTenantKeyRing(tenantID string) (*TenantKeyRing, error) {
	u, err := user.Current()
	if err != nil {
		return nil, errors.Wrapf(err, "get username")
	}
	return &TenantKeyRing{
		user:     u.Name,
		tenantID: tenantID,
	}, nil
}

func (kr *TenantKeyRing) GetToken() (*api.TenantToken, error) {
	tokenStr, err := keyring.Get(kr.user, kr.tenantID)
	if err != nil {
		return nil, errors.Wrapf(err, "get token")
	}

	var token api.TenantToken
	if err := json.Unmarshal([]byte(tokenStr), &token); err != nil {
		return nil, errors.Wrapf(err, "unmarshal token")
	}
	return &token, nil
}

func (kr *TenantKeyRing) SetToken(tok *api.TenantToken) error {
	tokenBytes, err := json.Marshal(tok)
	if err != nil {
		return errors.Wrapf(err, "marshal token")
	}

	_ = keyring.Delete(kr.user, kr.tenantID)

	tokenStr := string(tokenBytes)
	if err := keyring.Set(kr.user, kr.tenantID, tokenStr); err != nil {
		return errors.Wrapf(err, "set token")
	}

	return nil
}

func (kr *TenantKeyRing) DelToken() error {
	if err := keyring.Delete(kr.user, kr.tenantID); err != nil {
		log.Printf("keyring delete %v", err)
	}
	return nil
}
