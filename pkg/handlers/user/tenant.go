package user

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/cc/config"
	"github.com/aserto-dev/aserto/pkg/client/tenant"
	"github.com/aserto-dev/aserto/pkg/keyring"
	"github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/go-grpc/aserto/tenant/account/v1"
	"github.com/aserto-dev/go-grpc/aserto/tenant/connection/v1"
	"gopkg.in/yaml.v2"

	"github.com/pkg/errors"
)

func getTenantID(ctx context.Context, c *cc.CommonCtx, client *tenant.Client, userIdentity string) (string, error) {
	resp, err := client.Account.GetAccount(ctx, &account.GetAccountRequest{})
	if err != nil {
		return "", errors.Wrapf(err, "get account")
	}

	if err := writeContexts(c, resp.Result.Tenants, resp.Result.DefaultTenant, userIdentity); err != nil {
		return "", err
	}

	return resp.Result.DefaultTenant, nil
}

func writeContexts(c *cc.CommonCtx, tenants []*api.Tenant, defaultTenant, userIdentity string) error {
	cfgSymlinkPath, err := config.GetSymlinkConfigPath()
	if err != nil {
		return err
	}
	cfgSymWOExt := strings.TrimSuffix(cfgSymlinkPath, ".yaml")
	configPath := fmt.Sprintf("%s-%s.yaml", cfgSymWOExt, userIdentity)

	cfg := &config.Config{}
	if config.FileExists(cfgSymlinkPath) {
		// user already logged in
		if err == nil {
			cfg, err = config.NewConfig(config.Path(cfgSymlinkPath))
			if err != nil {
				return err
			}
			if len(cfg.Context.Contexts) != 0 {
				return nil
			}
		}
	}

	if config.FileExists(configPath) {
		cfg, err = config.NewConfig(config.Path(configPath))
		if err != nil {
			return err
		}
		// if the user was logged in before, just create the symlink
		if !config.FileExists(cfgSymlinkPath) {
			err := os.Symlink(configPath, cfgSymlinkPath)
			if err != nil {
				return err
			}
		}
		if len(cfg.Context.Contexts) != 0 {
			return nil
		}
	}

	var activeTenant string
	cfg.Context.Contexts = make([]config.Ctx, 0)
	for _, tnt := range tenants {
		cfg.Context.Contexts = append(cfg.Context.Contexts, config.Ctx{Name: tnt.Name, TenantID: tnt.Id})
		if defaultTenant == tnt.Id {
			activeTenant = tnt.Name
		}
	}

	cfg.Context.ActiveContext = activeTenant
	cfg.Services = c.Environment
	cfg.Auth = &config.Auth{Issuer: c.Auth.Issuer, ClientID: c.Auth.ClientID, Audience: c.Auth.Audience}

	fileContent, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, fileContent, 0600)
	if err != nil {
		return err
	}

	return os.Symlink(configPath, cfgSymlinkPath)
}

func GetConnectionKeys(ctx context.Context, client *tenant.Client, token *keyring.TenantToken) error {
	client.SetTenantID(token.TenantID)

	resp, err := client.Connections.ListConnections(
		ctx,
		&connection.ListConnectionsRequest{
			Kind: api.ProviderKind_PROVIDER_KIND_UNKNOWN,
		})
	if err != nil {
		return errors.Wrapf(err, "list connections account")
	}

	//nolint:exhaustive // we only care about these two provider kinds.
	for _, cn := range resp.Results {
		switch {
		case cn.Kind == api.ProviderKind_PROVIDER_KIND_AUTHORIZER:
			if respX, err := GetConnection(ctx, client, cn.Id); err == nil {
				token.AuthorizerAPIKey = respX.Result.Config.Fields["api_key"].GetStringValue()
			} else {
				return errors.Wrapf(err, "get connection [%s]", cn.Id)
			}
		case cn.Kind == api.ProviderKind_PROVIDER_KIND_DECISION_LOGS:
			if respX, err := GetConnection(ctx, client, cn.Id); err == nil {
				token.DecisionLogsKey = respX.Result.Config.Fields["api_key"].GetStringValue()
			} else {
				return errors.Wrapf(err, "get connection [%s]", cn.Id)
			}
		case cn.Kind == api.ProviderKind_PROVIDER_KIND_DISCOVERY:
			if respX, err := GetConnection(ctx, client, cn.Id); err == nil {
				token.DiscoveryKey = respX.Result.Config.Fields["api_key"].GetStringValue()
			} else {
				return errors.Wrapf(err, "get connection [%s]", cn.Id)
			}
		case cn.Kind == api.ProviderKind_PROVIDER_KIND_DIRECTORY && cn.Name == "aserto-directory":
			if respX, err := GetConnection(ctx, client, cn.Id); err == nil {
				token.DirectoryReadKey = respX.Result.Config.Fields["api_key_read"].GetStringValue()
				token.DirectoryWriteKey = respX.Result.Config.Fields["api_key_write"].GetStringValue()
			} else {
				return errors.Wrapf(err, "get connection [%s]", cn.Id)
			}
		case cn.Kind == api.ProviderKind_PROVIDER_KIND_POLICY_REGISTRY && cn.Name == "aserto-policy-registry":
			if respX, err := GetConnection(ctx, client, cn.Id); err == nil {
				token.RegistryDownloadKey = respX.Result.Config.Fields["api_key_read"].GetStringValue()
				token.RegistryUploadKey = respX.Result.Config.Fields["api_key_write"].GetStringValue()
			} else {
				return errors.Wrapf(err, "get connection [%s]", cn.Id)
			}
		}
	}

	return nil
}

func GetConnection(
	ctx context.Context,
	client *tenant.Client,
	connectionID string,
) (*connection.GetConnectionResponse, error) {
	return client.Connections.GetConnection(
		ctx,
		&connection.GetConnectionRequest{Id: connectionID},
	)
}
