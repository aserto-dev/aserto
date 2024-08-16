package clients

import (
	"strings"

	"github.com/aserto-dev/aserto/pkg/cc/config"
	"github.com/aserto-dev/aserto/pkg/cc/errors"
	"github.com/aserto-dev/aserto/pkg/cc/token"
	"github.com/aserto-dev/aserto/pkg/x"
	client "github.com/aserto-dev/go-aserto"
)

// Overrides are options that modify the default behavior of connections to aserto services.
type Overrides interface {
	// Address overrides the default address used to connect to a service.
	Address() string

	// Key provides an API key to be used instead of the default access token.
	Key() string

	// IsAnonymous means no credentials are sent to the service.
	IsAnonymous() bool

	// IsInsecure indicates that no TLS verification is performed.
	IsInsecure() bool
}

type ServiceOptions struct {
	needsToken bool
	overrides  map[x.Service]Overrides
}

func NewServiceOptions() *ServiceOptions {
	return &ServiceOptions{overrides: map[x.Service]Overrides{}}
}

func (b *ServiceOptions) Override(svc x.Service, overrides Overrides) {
	b.overrides[svc] = overrides
}

func (b *ServiceOptions) RequireToken() {
	b.needsToken = true
}

func (b *ServiceOptions) ConfigOverrider(cfg *config.Config) {
	for svc, overrides := range b.overrides {
		options := cfg.Services.Get(svc)
		options.Address = overrides.Address()
		options.APIKey = overrides.Key()
		options.Anonymous = overrides.IsAnonymous()
		options.Insecure = overrides.IsInsecure()
	}
}

type optionsBuilder struct {
	service     x.Service
	options     *x.ServiceOptions
	defaultAddr string
	tenantID    string
	token       *token.CachedToken
}

func (c *optionsBuilder) ConnectionOptions() ([]client.ConnectionOption, error) {
	authOption, err := c.authOption()
	if err != nil {
		return nil, err
	}

	tenantOption, err := c.tenantOption()
	if err != nil {
		return nil, err
	}

	caCertPathOption := nilOption

	return []client.ConnectionOption{
		client.WithAddr(c.address()),
		client.WithInsecure(c.options.Insecure),
		authOption,
		tenantOption,
		caCertPathOption,
	}, nil
}

func (c *optionsBuilder) address() string {
	addr := c.options.Address
	if addr != "" {
		return addr
	}

	return c.defaultAddr
}

func (c *optionsBuilder) authOption() (client.ConnectionOption, error) {
	if c.options.Anonymous {
		return nilOption, nil
	}

	if c.options.APIKey != "" {
		return client.WithAPIKeyAuth(c.options.APIKey), nil
	}

	if err := c.token.Verify(); err != nil {
		return nil, err
	}

	tkn, err := c.token.Get()
	if err != nil {
		return nil, err
	}
	return client.WithTokenAuth(tkn.Access), nil
}

func (c *optionsBuilder) tenantOption() (client.ConnectionOption, error) {
	if c.tenantID != "" {
		return client.WithTenantID(c.tenantID), nil
	}

	if c.token.TenantID() != "" {
		return client.WithTenantID(c.token.TenantID()), nil
	}

	if !c.isHosted() {
		return nilOption, nil
	}

	return nil, errors.NeedTenantIDErr
}

func (c *optionsBuilder) isHosted() bool {
	return strings.Contains(c.address(), "aserto.com")
}

func nilOption(*client.ConnectionOptions) error {
	return nil
}
