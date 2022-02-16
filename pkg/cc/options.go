package cc

import (
	"log"

	aserto "github.com/aserto-dev/aserto-go/client"
	"github.com/aserto-dev/aserto/pkg/client"
	"github.com/aserto-dev/aserto/pkg/paths"
	"github.com/aserto-dev/aserto/pkg/x"
)

type optionsBuilder struct {
	client.Overrides

	service     x.Service
	defaultAddr string
	tenantID    string
	token       CachedToken
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

	return aserto.WithTokenAuth(c.token.Get().Access), nil
}

func (c *optionsBuilder) tenantOption() (aserto.ConnectionOption, error) {
	if c.tenantID != "" {
		return aserto.WithTenantID(c.tenantID), nil
	}

	if !c.isHosted() {
		return nilOption, nil
	}

	return nil, NeedTenantIDErr
}

func (c *optionsBuilder) isHosted() bool {
	return c.address() == c.defaultAddr
}

func nilOption(*aserto.ConnectionOptions) error {
	return nil
}
