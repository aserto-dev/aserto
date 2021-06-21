// +build darwin,cgo linux windows
// +build amd64

package keyring

import (
	"encoding/json"

	kr99 "github.com/99designs/keyring"
	"github.com/aserto-dev/aserto/pkg/auth0/api"
	x "github.com/aserto-dev/aserto/pkg/x"
	"github.com/pkg/errors"
)

type KeyRing struct {
	Store kr99.Keyring
}

func NewKeyRing() (*KeyRing, error) {
	kr, err := kr99.Open(kr99.Config{
		AllowedBackends: []kr99.BackendType{
			kr99.KeychainBackend,      // MacOS keychain
			kr99.SecretServiceBackend, // Linux libsecret via dbus
			kr99.WinCredBackend,       // Windows credential manager
		},
		ServiceName:                    "aserto",
		KeychainName:                   "login",
		KeychainTrustApplication:       true,
		KeychainSynchronizable:         false,
		KeychainAccessibleWhenUnlocked: true,
		LibSecretCollectionName:        "aserto",
		WinCredPrefix:                  "aserto-",
	})

	if err != nil {
		return &KeyRing{}, err
	}

	return &KeyRing{kr}, nil
}

func key(env string) string {
	if env == x.EnvProduction {
		return x.Aserto
	}
	return x.Aserto + "-" + env
}

func (kr *KeyRing) GetToken(env string) (*api.Token, error) {
	item, err := kr.Store.Get(key(env))
	if err != nil {
		return nil, errors.Wrapf(err, "get token")
	}

	var token api.Token
	if err := json.Unmarshal(item.Data, &token); err != nil {
		return nil, errors.Wrapf(err, "unmarshal token")
	}
	return &token, nil
}

func (kr *KeyRing) SetToken(env string, tok *api.Token) error {
	tokBytes, err := json.Marshal(tok)
	if err != nil {
		return errors.Wrapf(err, "marshal token")
	}

	if err := kr.Store.Set(kr99.Item{
		Key:         key(env),
		Data:        tokBytes,
		Label:       key(env),
		Description: "application credentials",
	}); err != nil {
		return errors.Wrapf(err, "set token")
	}
	return nil
}

func (kr *KeyRing) DelToken(env string) error {
	_ = kr.Store.Remove(key(env))
	return nil
}
