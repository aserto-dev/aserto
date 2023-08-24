package token

import (
	"log"
	"sync"

	"github.com/aserto-dev/aserto/pkg/auth0/api"
	errs "github.com/aserto-dev/aserto/pkg/cc/errors"
	"github.com/aserto-dev/aserto/pkg/keyring"
)

type CachedToken struct {
	token    *api.Token
	tenantID string

	verifyOnce sync.Once
	errVerify  error
}

func New(token *api.Token) *CachedToken {
	tenantID := token.DefaultTenantID

	return &CachedToken{token: token, tenantID: tenantID}
}

type CacheKey string

func Load(key CacheKey) *CachedToken {
	token := loadToken(key)
	tenantID := token.DefaultTenantID

	return &CachedToken{token: token, tenantID: tenantID}
}

func (t *CachedToken) Get() (*api.Token, error) {
	t.verifyOnce.Do(func() {
		t.errVerify = t.Verify()
	})

	return t.token, t.errVerify
}

func (t *CachedToken) Verify() error {
	if t.token == nil || t.token.Access == "" {
		return errs.NeedLoginErr
	}

	if t.token.IsExpired() {
		return errs.TokenExpiredErr
	}

	return nil
}

func (t *CachedToken) TenantID() string {
	return t.tenantID
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
