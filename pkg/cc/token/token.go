package token

import (
	"log"
	"sync"

	"github.com/aserto-dev/aserto/pkg/auth0/api"
	"github.com/aserto-dev/aserto/pkg/cc/errors"
	"github.com/aserto-dev/aserto/pkg/keyring"
)

type CachedToken struct {
	token *api.Token

	verifyOnce sync.Once
	errVerify  error
}

func New(token *api.Token) *CachedToken {
	return &CachedToken{token: token}
}

type CacheKey string

func Load(key CacheKey) *CachedToken {
	return &CachedToken{token: loadToken(key)}
}

func (t *CachedToken) Get() (*api.Token, error) {
	t.verifyOnce.Do(func() {
		t.errVerify = t.Verify()
	})

	return t.token, t.errVerify
}

func (t *CachedToken) Verify() error {
	if t.token == nil || t.token.Access == "" {
		return errors.NeedLoginErr
	}

	if t.token.IsExpired() {
		return errors.TokenExpiredErr
	}

	return nil
}

func (t *CachedToken) TenantID() string {
	token, err := t.Get()
	if err != nil {
		return ""
	}

	return token.TenantID
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
