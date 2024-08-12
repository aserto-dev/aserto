package user

import (
	"context"

	auth0 "github.com/aserto-dev/aserto/pkg/auth0/api"
	"github.com/aserto-dev/aserto/pkg/clients/tenant"
	client "github.com/aserto-dev/go-aserto"
	"github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/go-grpc/aserto/tenant/account/v1"
	"github.com/aserto-dev/go-grpc/aserto/tenant/connection/v1"

	"github.com/pkg/errors"
)

func getTenantID(ctx context.Context, tc *tenant.Client, token *auth0.Token) error {
	resp, err := tc.Account.GetAccount(ctx, &account.GetAccountRequest{})
	if err != nil {
		return errors.Wrapf(err, "get account")
	}

	token.TenantID = resp.Result.DefaultTenant

	return err
}

func GetConnectionKeys(ctx context.Context, tc *tenant.Client, token *auth0.Token) error {
	tenantContext, cancel := context.WithCancel(client.SetTenantContext(ctx, token.TenantID))
	defer cancel()

	resp, err := tc.Connections.ListConnections(
		tenantContext,
		&connection.ListConnectionsRequest{
			Kind: api.ProviderKind_PROVIDER_KIND_UNKNOWN,
		})
	if err != nil {
		return errors.Wrapf(err, "list connections account")
	}

	//nolint:exhaustive // we only care about these provider kinds.
	for _, cn := range resp.Results {
		switch {
		case cn.Kind == api.ProviderKind_PROVIDER_KIND_AUTHORIZER:
			if respX, err := GetConnection(tenantContext, tc, cn.Id); err == nil {
				token.AuthorizerAPIKey = respX.Result.Config.Fields["api_key"].GetStringValue()
			} else {
				return errors.Wrapf(err, "get authorizer connection [%s]", cn.Id)
			}
		case cn.Kind == api.ProviderKind_PROVIDER_KIND_DECISION_LOGS:
			if respX, err := GetConnection(tenantContext, tc, cn.Id); err == nil {
				token.DecisionLogsKey = respX.Result.Config.Fields["api_key"].GetStringValue()
			} else {
				return errors.Wrapf(err, "get decision-logs connection [%s]", cn.Id)
			}
		case cn.Kind == api.ProviderKind_PROVIDER_KIND_DISCOVERY:
			if respX, err := GetConnection(tenantContext, tc, cn.Id); err == nil {
				token.DiscoveryKey = respX.Result.Config.Fields["api_key"].GetStringValue()
			} else {
				return errors.Wrapf(err, "get discovery connection [%s]", cn.Id)
			}
		case cn.Kind == api.ProviderKind_PROVIDER_KIND_DIRECTORY && cn.Name == "aserto-directory":
			if respX, err := GetConnection(tenantContext, tc, cn.Id); err == nil {
				token.DirectoryReadKey = respX.Result.Config.Fields["api_key_read"].GetStringValue()
				token.DirectoryWriteKey = respX.Result.Config.Fields["api_key_write"].GetStringValue()
			} else {
				return errors.Wrapf(err, "get directory connection [%s]", cn.Id)
			}
		case cn.Kind == api.ProviderKind_PROVIDER_KIND_POLICY_REGISTRY && cn.Name == "aserto-policy-registry":
			if respX, err := GetConnection(tenantContext, tc, cn.Id); err == nil {
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
	tc *tenant.Client,
	connectionID string,
) (*connection.GetConnectionResponse, error) {
	return tc.Connections.GetConnection(
		ctx,
		&connection.GetConnectionRequest{Id: connectionID},
	)
}
