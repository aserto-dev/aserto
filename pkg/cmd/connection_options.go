package cmd

import (
	"github.com/aserto-dev/aserto/pkg/client"
	"github.com/aserto-dev/aserto/pkg/x"
)

type ConnectionOptions struct {
	APIKey   string `xor:"creds" env:"KEY" help:"service api key"`
	NoAuth   bool   `xor:"creds" help:"do not provide any credentials"`
	Insecure bool   `help:"skip TLS verification"`
}

func (so *ConnectionOptions) Address() string {
	return ""
}

func (so *ConnectionOptions) Key() string {
	return so.APIKey
}

func (so *ConnectionOptions) IsAnonymous() bool {
	return so.NoAuth
}

func (so *ConnectionOptions) IsInsecure() bool {
	return so.Insecure
}

type AuthorizerOptions struct {
	AddressOverride string `name:"authorizer" env:"ADDRESS" help:"authorizer override"`

	ConnectionOptions
}

func (ao *AuthorizerOptions) Address() string {
	return ao.AddressOverride
}

type ConnectionOverrides interface {
	Override(svc x.Service, overrides client.Overrides)
	RequireToken()
}
