package tenant

import (
	"context"

	"github.com/aserto-dev/go-aserto/client"
	"google.golang.org/grpc"

	info "github.com/aserto-dev/go-grpc/aserto/common/info/v1"
	account "github.com/aserto-dev/go-grpc/aserto/tenant/account/v1"
	connection "github.com/aserto-dev/go-grpc/aserto/tenant/connection/v1"
	onboarding "github.com/aserto-dev/go-grpc/aserto/tenant/onboarding/v1"
	policy "github.com/aserto-dev/go-grpc/aserto/tenant/policy/v1"
	policy_builder "github.com/aserto-dev/go-grpc/aserto/tenant/policy_builder/v1"
	profile "github.com/aserto-dev/go-grpc/aserto/tenant/profile/v1"
	provider "github.com/aserto-dev/go-grpc/aserto/tenant/provider/v1"
	registry "github.com/aserto-dev/go-grpc/aserto/tenant/registry/v1"
	scc "github.com/aserto-dev/go-grpc/aserto/tenant/scc/v1"
	v2 "github.com/aserto-dev/go-grpc/aserto/tenant/v2"

	"github.com/pkg/errors"
)

// Client provides access to the Aserto control plane services.
type Client struct {
	conn *client.Connection

	// Account provides methods for managing a customer account.
	Account account.AccountClient

	// Connections provides methods to create and manage connections.
	Connections connection.ConnectionClient

	// Onboarding provides methods to create tenants and invite users.
	Onboarding onboarding.OnboardingClient

	// Policy provides methods for creating, listing, and deleting policy references.
	Policy policy.PolicyClient

	// PolicyBuilder provides methods for creating and managing policy builders.
	PolicyBuilder policy_builder.PolicyBuilderClient

	// Profile provides methods for managing user invitations.
	Profile profile.ProfileClient

	// Provider provides methods for viewing the providers available to create connections.
	Provider provider.ProviderClient

	// Registry provides methods for managing registry repositories.
	Registry registry.RegistryClient

	// SCC provides methods for interacting with a tenant's configured source-control repositories.
	SCC scc.SourceCodeCtlClient

	// Info provides methods for retrieving system information and configuration.
	Info info.InfoClient

	// V2Policy provides methods for for creating and managing policies
	V2Policy v2.PolicyClient

	// V2Repository provides methods for for creating and managing policy's repositories
	V2Repository v2.RepositoryClient

	// V2Source provides methods for for creating and managing policy's sources
	V2Source v2.SourceClient

	// V2Instance provides methods for for creating and managing policy's instances
	V2Instance v2.InstanceClient

	// V2Tenant provides methods for for creating and managing tenant instances
	V2Tenant v2.TenantClient
}

// New creates a tenant Client with the specified connection options.
func New(ctx context.Context, opts ...client.ConnectionOption) (*Client, error) {
	conn, err := client.NewConnection(ctx, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "create grpc client failed")
	}

	return &Client{
		conn:          conn,
		Account:       account.NewAccountClient(conn.Conn),
		Connections:   connection.NewConnectionClient(conn.Conn),
		Onboarding:    onboarding.NewOnboardingClient(conn.Conn),
		Policy:        policy.NewPolicyClient(conn.Conn),
		PolicyBuilder: policy_builder.NewPolicyBuilderClient(conn.Conn),
		Profile:       profile.NewProfileClient(conn.Conn),
		Provider:      provider.NewProviderClient(conn.Conn),
		Registry:      registry.NewRegistryClient(conn.Conn),
		SCC:           scc.NewSourceCodeCtlClient(conn.Conn),
		Info:          info.NewInfoClient(conn.Conn),
		V2Policy:      v2.NewPolicyClient(conn.Conn),
		V2Repository:  v2.NewRepositoryClient(conn.Conn),
		V2Source:      v2.NewSourceClient(conn.Conn),
		V2Instance:    v2.NewInstanceClient(conn.Conn),
		V2Tenant:      v2.NewTenantClient(conn.Conn),
	}, err
}

// SetTenantID provides a tenantID to be included in outgoing messages.
func (c *Client) SetTenantID(tenantID string) {
	c.conn.TenantID = tenantID
}

// Connection returns the underlying grpc connection.
func (c *Client) Connection() grpc.ClientConnInterface {
	return c.conn.Conn
}
