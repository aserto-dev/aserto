package keyring

import (
	"encoding/json"
	"log"
	"os/user"

	"github.com/aserto-dev/aserto/pkg/auth0/api"
	"github.com/pkg/errors"
	"github.com/zalando/go-keyring"
)

type KeyRing struct {
	service string
	user    string
}

func NewKeyRing(key string) (*KeyRing, error) {
	u, err := user.Current()
	if err != nil {
		return nil, errors.Wrapf(err, "get username")
	}
	return &KeyRing{
		service: key,
		user:    u.Username,
	}, nil
}

func (kr *KeyRing) GetToken() (*api.Token, error) {
	tokenStr, err := keyring.Get(kr.service, kr.user)
	if err != nil {
		return nil, errors.Wrapf(err, "get token")
	}

	var token api.Token
	if err := json.Unmarshal([]byte(tokenStr), &token); err != nil {
		return nil, errors.Wrapf(err, "unmarshal token")
	}
	return &token, nil
}

func (kr *KeyRing) SetToken(tok *api.Token) error {
	tokenBytes, err := json.Marshal(tok)
	if err != nil {
		return errors.Wrapf(err, "marshal token")
	}

	tokenStr := string(tokenBytes)
	if err := keyring.Set(kr.service, kr.user, tokenStr); err != nil {
		return errors.Wrapf(err, "set token")
	}

	return nil
}

func (kr *KeyRing) DelToken() error {
	if err := keyring.Delete(kr.service, kr.user); err != nil {
		log.Printf("keyring delete %v", err)
	}
	return nil
}
