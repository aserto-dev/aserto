package grpcc

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func SetTenantContext(ctx context.Context, tenantID string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, asertoTenantID, tenantID)
}

func SetAsertoAPIKey(ctx context.Context, key string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, asertoAPIKey, authzBasicHeader(key))
}

func authzBasicHeader(key string) string {
	return basic + " " + key
}

func SetAccountContext(ctx context.Context, accountID string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, asertoAccountID, accountID)
}
