package cc

import (
	"log"

	"github.com/aserto-dev/aserto/pkg/auth0/api"
	"github.com/aserto-dev/aserto/pkg/keyring"
)

type CachedToken struct {
	token *api.Token
}

func NewCachedToken(env string) CachedToken {
	token := loadToken(env)

	return CachedToken{token: token}
}

func (t CachedToken) Get() *api.Token {
	return t.token
}

func (t CachedToken) Verify() error {
	if t.token == nil || t.token.Access == "" {
		return NeedLoginErr
	}

	if t.token.IsExpired() {
		return TokenExpiredErr
	}

	return nil
}

func (t CachedToken) TenantID() string {
	if t.token != nil && !t.token.IsExpired() {
		return t.token.TenantID
	}

	return ""
}

func loadToken(env string) *api.Token {
	kr, err := keyring.NewKeyRing(env)
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
