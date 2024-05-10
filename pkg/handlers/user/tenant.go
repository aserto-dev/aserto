package user

import (
	"context"

	auth0 "github.com/aserto-dev/aserto/pkg/auth0/api"
	"github.com/aserto-dev/aserto/pkg/client/tenant"
	aserto "github.com/aserto-dev/go-aserto/client"
	"github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/go-grpc/aserto/tenant/account/v1"
	"github.com/aserto-dev/go-grpc/aserto/tenant/connection/v1"

	"github.com/pkg/errors"
)

func getTenantID(ctx context.Context, client *tenant.Client, token *auth0.Token) error {
	resp, err := client.Account.GetAccount(ctx, &account.GetAccountRequest{})
	if err != nil {
		return errors.Wrapf(err, "get account")
	}

	token.TenantID = resp.Result.DefaultTenant

	return err
}

func GetConnectionKeys(ctx context.Context, client *tenant.Client, token *auth0.Token) error {
	resp, err := client.Connections.ListConnections(
		aserto.SetTenantContext(ctx, token.TenantID),
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
				return errors.Wrapf(err, "get authorizer connection [%s]", cn.Id)
			}
		case cn.Kind == api.ProviderKind_PROVIDER_KIND_DECISION_LOGS:
			if respX, err := GetConnection(ctx, client, cn.Id); err == nil {
				token.DecisionLogsKey = respX.Result.Config.Fields["api_key"].GetStringValue()
			} else {
				return errors.Wrapf(err, "get decision-logs connection [%s]", cn.Id)
			}
		case cn.Kind == api.ProviderKind_PROVIDER_KIND_DISCOVERY:
			if respX, err := GetConnection(ctx, client, cn.Id); err == nil {
				token.DiscoveryKey = respX.Result.Config.Fields["api_key"].GetStringValue()
			} else {
				return errors.Wrapf(err, "get discovery connection [%s]", cn.Id)
			}
		case cn.Kind == api.ProviderKind_PROVIDER_KIND_DIRECTORY && cn.Name == "aserto-directory":
			if respX, err := GetConnection(ctx, client, cn.Id); err == nil {
				token.DirectoryReadKey = respX.Result.Config.Fields["api_key_read"].GetStringValue()
				token.DirectoryWriteKey = respX.Result.Config.Fields["api_key_write"].GetStringValue()
			} else {
				return errors.Wrapf(err, "get directory connection [%s]", cn.Id)
			}
		case cn.Kind == api.ProviderKind_PROVIDER_KIND_POLICY_REGISTRY && cn.Name == "aserto-policy-registry":
			if respX, err := GetConnection(ctx, client, cn.Id); err == nil {
				token.RegistryDownloadKey = respX.Result.Config.Fields["api_key_read"].GetStringValue()
				token.RegistryUploadKey = respX.Result.Config.Fields["api_key_write"].GetStringValue()
			} else {
				return errors.Wrapf(err, "get registry connection [%s]", cn.Id)
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
