package user

import (
	"context"
	"fmt"
	"time"

	"github.com/aserto-dev/aserto/pkg/auth0/api"
	"github.com/aserto-dev/aserto/pkg/auth0/device"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/clients/tenant"
	"github.com/aserto-dev/aserto/pkg/keyring"
	client "github.com/aserto-dev/go-aserto"

	"github.com/cli/browser"
	"github.com/pkg/errors"
)

type LoginCmd struct {
	Browser   bool   `flag:"browser" negatable:"" default:"true" help:"do not open browser"`
	TenantID  string `flag:"tenant-id" hidden:""`
	TenantKey string `flag:"tenant-key" hidden:""`
}

func (d *LoginCmd) Run(c *cc.CommonCtx) error {
	if d.TenantID != "" && d.TenantKey != "" {
		return d.run(c)
	}

	settings := c.Auth

	flow := device.New(
		device.WithClientID(settings.ClientID),
		device.WithDeviceAuthorizationURL(settings.DeviceAuthorizationURL),
		device.WithTokenURL(settings.TokenURL),
		device.WithAudience(settings.Audience),
		device.WithGrantType(settings.GrantType),
		device.WithScope("openid", "profile", "email"),
	)

	if err := flow.GetDeviceCode(c.Context); err != nil {
		return err
	}

	fmt.Printf("Copy your one-time code: %s\n", flow.GetUserCode())

	if d.Browser {
		fmt.Printf("Press Enter to open browser %s\n", flow.GetVerificationURI())
		fmt.Scanln() //nolint: errcheck
		if err := browser.OpenURL(flow.GetVerificationURI()); err != nil {
			return err
		}
	} else {
		fmt.Printf("Open browser %s\n", flow.GetVerificationURI())
	}

	ctx, cancel := context.WithTimeout(c.Context, flow.ExpiresIn())
	defer cancel()

	for {
		if ok, err := flow.RequestAccessToken(ctx); ok {
			fmt.Fprintln(c.StdOut(), ".")
			break
		} else if err != nil {
			return err
		}

		select {
		case <-time.After(flow.Interval()):
			fmt.Fprint(c.StdOut(), ".")
		case <-ctx.Done():
			return errors.New("canceled")
		}
	}

	if err := UpdateKeyRing(c, flow.AccessToken()); err != nil {
		return err
	}

	c.Con().Info().Msg("Login successful")

	return nil
}

func UpdateKeyRing(c *cc.CommonCtx, token *api.Token) error {
	ctx, cancel := context.WithTimeout(c.Context, time.Second*10)
	defer cancel()

	conn, err := tenant.NewClient(
		ctx,
		client.WithAddr(c.Environment.TenantService.Address),
		client.WithTokenAuth(token.Access),
	)
	if err != nil {
		return err
	}

	if err = getTenantID(ctx, conn, token); err != nil {
		return errors.Wrapf(err, "get tenant id")
	}

	if err = GetConnectionKeys(ctx, conn, token); err != nil {
		return errors.Wrapf(err, "get connection keys")
	}

	kr, err := keyring.NewKeyRing(c.Auth.Issuer)
	if err != nil {
		return err
	}
	if c.TenantID() != "" {
		token.TenantID = c.TenantID()
	}
	if err := kr.SetToken(token); err != nil {
		return err
	}

	c.Con().Info().Msg("Switched to tenant-id %q", c.TenantID())

	return nil
}

func SwitchKeyRing(c *cc.CommonCtx, token *api.Token, tenantID string) error {
	ctx, cancel := context.WithTimeout(c.Context, time.Second*10)
	defer cancel()

	conn, err := tenant.NewClient(
		ctx,
		client.WithAddr(c.Environment.TenantService.Address),
		client.WithTokenAuth(token.Access),
	)
	if err != nil {
		return err
	}

	token.TenantID = tenantID

	if err = GetConnectionKeys(ctx, conn, token); err != nil {
		return errors.Wrapf(err, "get connection keys")
	}

	kr, err := keyring.NewKeyRing(c.Auth.Issuer)
	if err != nil {
		return err
	}

	if err := kr.SetToken(token); err != nil {
		return err
	}

	c.Con().Info().Msg("Switched to tenant-id %q", c.TenantID())

	return nil
}

func (d *LoginCmd) run(c *cc.CommonCtx) error {
	ctx, cancel := context.WithTimeout(c.Context, time.Second*10)
	defer cancel()

	conn, err := tenant.NewClient(
		ctx,
		client.WithAddr(c.Environment.TenantService.Address),
		client.WithAPIKeyAuth(d.TenantKey),
	)
	if err != nil {
		return err
	}

	tok := api.Token{
		Type:      "basic",
		Scope:     "scope",
		Identity:  "identity token",
		Access:    "access token",
		ExpiresIn: 3600,
		ExpiresAt: time.Now().Add(3600 * time.Second),
		TenantID:  d.TenantID,
	}

	if err := GetConnectionKeys(ctx, conn, &tok); err != nil {
		return err
	}

	kr, err := keyring.NewKeyRing(c.Auth.Issuer)
	if err != nil {
		return err
	}

	if err := kr.SetToken(&tok); err != nil {
		return err
	}

	c.Config.TenantID = d.TenantID

	c.Con().Info().Msg("Switched to tenant-id %q", c.TenantID())

	return nil
}
