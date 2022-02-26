package token

import (
	"log"

	"github.com/aserto-dev/aserto/pkg/auth0/api"
	"github.com/aserto-dev/aserto/pkg/cc/errors"
	"github.com/aserto-dev/aserto/pkg/keyring"
)

type CachedToken struct {
	token *api.Token
}

func New(token *api.Token) CachedToken {
	return CachedToken{token: token}
}

type CacheKey string

func Load(key CacheKey) CachedToken {
	return CachedToken{token: loadToken(key)}
}

func (t CachedToken) Get() *api.Token {
	return t.token
}

func (t CachedToken) Verify() error {
	if t.token == nil || t.token.Access == "" {
		return errors.NeedLoginErr
	}

	if t.token.IsExpired() {
		return errors.TokenExpiredErr
	}

	return nil
}

func (t CachedToken) TenantID() string {
	if t.token != nil && !t.token.IsExpired() {
		return t.token.TenantID
	}

	return ""
}

func loadToken(key CacheKey) *api.Token {
	kr, err := keyring.NewKeyRing(string(key))
	if err != nil {
		log.Printf("token: instantiating keyring, %s", err.Error())
		return nil
	}

	token, err := kr.GetToken()
	if err != nil {
		return nil
	}

	return token
}
