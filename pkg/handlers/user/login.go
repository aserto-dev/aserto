package user

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/aserto-dev/aserto/pkg/auth0/device"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/client/tenant"
	"github.com/aserto-dev/aserto/pkg/keyring"
	aserto "github.com/aserto-dev/go-aserto/client"

	"github.com/cli/browser"
	"github.com/pkg/errors"
)

type LoginCmd struct {
	Browser bool `flag:"browser" negatable:"" default:"true" help:"do not open browser"`
}

func (d *LoginCmd) Run(c *cc.CommonCtx) error {
	settings := c.Auth

	flow := device.New(
		device.WithClientID(settings.ClientID),
		device.WithDeviceAuthorizationURL(settings.DeviceAuthorizationURL),
		device.WithTokenURL(settings.TokenURL),
		device.WithAudience(settings.Audience),
		device.WithGrantType(settings.GrantType),
		device.WithScope("openid", "profile", "email"),
	)

	// handle Ctrl+C
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := flow.GetDeviceCode(ctx); err != nil {
		return err
	}

	fmt.Printf("Copy your one-time code: %s\n", flow.GetUserCode())

	if d.Browser {
		fmt.Printf("Press Enter to open browser %s\n", flow.GetVerificationURI())
		fmt.Scanln()
		if err := browser.OpenURL(flow.GetVerificationURI()); err != nil {
			return err
		}
	} else {
		fmt.Printf("Open browser %s\n", flow.GetVerificationURI())
	}

	ctx, cancel := context.WithTimeout(ctx, flow.ExpiresIn())
	defer cancel()

	for {
		if ok, err := flow.RequestAccessToken(ctx); ok {
			fmt.Fprint(c.UI.Output(), ".\n")
			break
		} else if err != nil {
			return err
		}

		select {
		case <-time.After(flow.Interval()):
			fmt.Fprint(c.UI.Output(), ".")
		case <-ctx.Done():
			return errors.New("canceled")
		}
	}

	token := flow.AccessToken()

	conn, err := tenant.New(
		c.Context,
		aserto.WithAddr(c.Environment.TenantService.Address),
		aserto.WithTokenAuth(token.Access),
	)
	if err != nil {
		return err
	}

	tenantID, err := getTenantID(ctx, c, conn, token.Subject)
	if err != nil {
		return errors.Wrapf(err, "get tenant ID ")
	}

	token.DefaultTenantID = tenantID
	kr, err := keyring.NewKeyRing(c.Auth.Issuer)
	if err != nil {
		return err
	}
	if err := kr.SetToken(token); err != nil {
		return err
	}

	tenantKr, err := keyring.NewTenantKeyRing(tenantID + "-" + token.Subject)
	if err != nil {
		return err
	}
	tenantToken := &keyring.TenantToken{TenantID: tenantID}

	if err = GetConnectionKeys(c.Context, conn, tenantToken); err != nil {
		return errors.Wrapf(err, "get connection keys")
	}

	if err := tenantKr.SetToken(tenantToken); err != nil {
		return err
	}

	fmt.Fprint(c.UI.Output(), "Login successful\n")

	return nil
}
