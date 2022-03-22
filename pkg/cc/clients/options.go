package clients

import (
	"log"

	aserto "github.com/aserto-dev/aserto-go/client"
	"github.com/aserto-dev/aserto/pkg/cc/errors"
	"github.com/aserto-dev/aserto/pkg/cc/token"
	"github.com/aserto-dev/aserto/pkg/paths"
	"github.com/aserto-dev/aserto/pkg/x"
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

func (b *ServiceOptions) serviceOverrides(svc x.Service) Overrides {
	overrides, ok := b.overrides[svc]
	if !ok {
		return &noOverrides{}
	}

	return overrides
}

type noOverrides struct{}

func (o *noOverrides) Address() string {
	return ""
}

func (o *noOverrides) Key() string {
	return ""
}

func (o *noOverrides) IsAnonymous() bool {
	return false
}

func (o *noOverrides) IsInsecure() bool {
	return false
}

type optionsBuilder struct {
	Overrides

	service     x.Service
	defaultAddr string
	tenantID    string
	token       *token.CachedToken
}

func (c *optionsBuilder) ConnectionOptions() ([]aserto.ConnectionOption, error) {
	authOption, err := c.authOption()
	if err != nil {
		return nil, err
	}

	tenantOption, err := c.tenantOption()
	if err != nil {
		return nil, err
	}

	caCertPathOption := nilOption
	if c.service == x.AuthorizerService && !c.isHosted() {
		// Look for a CA cert
		p, err := paths.New()
		if err == nil {
			caCertPathOption = aserto.WithCACertPath(p.Certs.GRPC.CA)
		} else {
			log.Println("Unable to locate onebox certificates.", err.Error())
		}
	}

	return []aserto.ConnectionOption{
		aserto.WithAddr(c.address()),
		aserto.WithInsecure(c.IsInsecure()),
		authOption,
		tenantOption,
		caCertPathOption,
	}, nil
}

func (c *optionsBuilder) address() string {
	addr := c.Address()
	if addr != "" {
		return addr
	}

	return c.defaultAddr
}

func (c *optionsBuilder) authOption() (aserto.ConnectionOption, error) {
	if c.IsAnonymous() {
		return nilOption, nil
	}

	if c.Key() != "" {
		return aserto.WithAPIKeyAuth(c.Key()), nil
	}

	if err := c.token.Verify(); err != nil {
		return nil, err
	}

	tkn, err := c.token.Get()
	if err != nil {
		return nil, err
	}
	return aserto.WithTokenAuth(tkn.Access), nil
}

func (c *optionsBuilder) tenantOption() (aserto.ConnectionOption, error) {
	if c.tenantID != "" {
		return aserto.WithTenantID(c.tenantID), nil
	}

	if !c.isHosted() {
		return nilOption, nil
	}

	return nil, errors.NeedTenantIDErr
}

func (c *optionsBuilder) isHosted() bool {
	return c.address() == c.defaultAddr
}

func nilOption(*aserto.ConnectionOptions) error {
	return nil
}
