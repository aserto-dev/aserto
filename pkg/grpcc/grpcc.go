package grpcc

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	clientTimeout = time.Duration(5) * time.Second
)

type Client struct {
	Conn *grpc.ClientConn // TODO: change into grpc.ClientConnInterface type.
}

// TokenAuth bearer token based authentication.
type TokenAuth struct {
	token string
}

func NewTokenAuth(token string) *TokenAuth {
	return &TokenAuth{
		token: token,
	}
}

func (t TokenAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		authorization: bearer + " " + t.token,
	}, nil
}

func (TokenAuth) RequireTransportSecurity() bool {
	return true
}

// APIKeyAuth API key based authentication
type APIKeyAuth struct {
	key string
}

func NewAPIKeyAuth(key string) *APIKeyAuth {
	return &APIKeyAuth{
		key: key,
	}
}

func (k *APIKeyAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		authorization: basic + " " + k.key,
	}, nil
}

func (k *APIKeyAuth) RequireTransportSecurity() bool {
	return true
}

func NewClient(ctx context.Context, dialAddr string, creds credentials.PerRPCCredentials) (*Client, error) {
	var insecure bool
	if strings.Contains(dialAddr, "localhost") {
		insecure = true
	}

	tlsConf, err := tlsConfig(insecure)
	if err != nil {
		return nil, errors.Wrap(err, "failed to setup tls configuration")
	}

	clientCreds := credentials.NewTLS(tlsConf)

	conn, err := grpc.DialContext(
		ctx,
		dialAddr,
		grpc.WithTransportCredentials(clientCreds),
		grpc.WithPerRPCCredentials(creds),
		grpc.WithBlock(),
		grpc.WithTimeout(clientTimeout), //nolint:staticcheck // can't release a with timeout context in this method
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to setup grpc dial context to %s", dialAddr)
	}

	return &Client{
		Conn: conn,
	}, nil
}
