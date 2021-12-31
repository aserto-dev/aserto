package dev

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"text/template"

	"github.com/aserto-dev/aserto-go/client/grpc/tenant"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/filex"
	"google.golang.org/protobuf/types/known/structpb"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/go-grpc/aserto/tenant/connection/v1"
	policy "github.com/aserto-dev/go-grpc/aserto/tenant/policy/v1"

	"github.com/pkg/errors"
)

type ConfigureCmd struct {
	Name   string `arg:"" required:"" help:"policy name"`
	Stdout bool   `short:"p" help:"generated configuration is printed to stdout but not saved"`
}

func (cmd ConfigureCmd) Run(c *cc.CommonCtx) error {
	fmt.Fprintf(c.ErrWriter, ">>> configure policy...\n")
	fmt.Fprintf(c.ErrWriter, "tenant id: %s\n", c.TenantID())

	client, err := tenant.New(c.Context, c.TenantSvcConnectionOptions()...)
	if err != nil {
		return err
	}

	policyRef, err := findPolicyRef(c.Context, client, cmd.Name)
	if err != nil {
		return err
	}

	discoveryConf, err := getDiscoveryConfig(c.Context, client)
	if err != nil {
		return err
	}

	params := templateParams{
		TenantID:     c.TenantID(),
		PolicyName:   policyRef.Name,
		PolicyID:     policyRef.Id,
		DiscoveryURL: discoveryConf.URL,
		TenantKey:    discoveryConf.APIKey,
	}

	if params.TenantKey == "" {
		return errors.Errorf("missing $ASERTO_TENANT_KEY env var")
	}

	fmt.Fprintf(c.ErrWriter, "policy id: %s\n", params.PolicyID)

	var w io.Writer

	if cmd.Stdout {
		w = c.OutWriter
	} else {
		w, err = ConfigFileWriter(params.PolicyName)
		if err != nil {
			return err
		}
	}

	return WriteConfig(w, configTemplate, &params)
}

func findPolicyRef(ctx context.Context, client *tenant.Client, policyName string) (*api.PolicyRef, error) {
	policyRefResp, err := client.Policy.ListPolicyRefs(ctx, &policy.ListPolicyRefsRequest{})
	if err != nil {
		return nil, err
	}

	for _, v := range policyRefResp.Results {
		if v.Name == policyName {
			return v, nil
		}
	}
	return nil, errors.Errorf("policy not found [%s]", policyName)
}

type discoveryConfig struct {
	URL    string
	APIKey string
}

func newDiscoveryConfig(config *structpb.Struct) (*discoveryConfig, error) {
	urlField, ok := config.Fields["url"]
	if !ok {
		return nil, errors.New("missing field: url")
	}

	apiKeyField, ok := config.Fields["api_key"]
	if !ok {
		return nil, errors.New("missing field: api_key")
	}

	return &discoveryConfig{URL: urlField.GetStringValue(), APIKey: apiKeyField.GetStringValue()}, nil
}

func getDiscoveryConfig(ctx context.Context, client *tenant.Client) (*discoveryConfig, error) {
	resp, err := client.Connections.ListConnections(
		ctx,
		&connection.ListConnectionsRequest{Kind: api.ProviderKind_PROVIDER_KIND_DISCOVERY},
	)
	if err != nil {
		return nil, err
	}

	if len(resp.Results) == 0 {
		return nil, errors.New("no discovery connections available for tenant. please contact support@aserto.com")
	}

	for _, conn := range resp.Results {
		conResp, err := client.Connections.GetConnection(ctx, &connection.GetConnectionRequest{Id: conn.Id})
		if err == nil {
			conf, err := newDiscoveryConfig(conResp.Result.Config)
			if err == nil {
				return conf, nil
			}
		}
	}

	return nil, errors.Errorf("cannot find discovery configuration")
}

func ConfigFileWriter(policyName string) (io.Writer, error) {
	configDir, err := CreateConfigDir()
	if err != nil {
		return nil, err
	}

	return os.Create(path.Join(configDir, policyName+".yaml"))
}

func CreateConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := path.Join(home, "/.config/aserto/aserto-one/cfg")
	if filex.DirExists(configDir) {
		return configDir, nil
	}
	return configDir, os.MkdirAll(configDir, 0700)
}

func WriteConfig(w io.Writer, templ string, params *templateParams) error {
	t, err := template.New("config").Parse(templ)
	if err != nil {
		return err
	}

	err = t.Execute(w, params)
	if err != nil {
		return err
	}

	return nil
}
