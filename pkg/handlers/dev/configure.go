package dev

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"text/template"

	"github.com/aserto-dev/aserto-go/client/tenant"
	"github.com/aserto-dev/aserto/pkg/cc"
	decisionlogger "github.com/aserto-dev/aserto/pkg/decision_logger"
	"github.com/aserto-dev/aserto/pkg/filex"
	"github.com/aserto-dev/aserto/pkg/x"
	"google.golang.org/protobuf/types/known/structpb"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/go-grpc/aserto/tenant/connection/v1"
	policy "github.com/aserto-dev/go-grpc/aserto/tenant/policy/v1"

	"github.com/pkg/errors"
)

type ConfigureCmd struct {
	Name            string `arg:"" required:"" help:"policy name"`
	Stdout          bool   `short:"p" help:"generated configuration is printed to stdout but not saved"`
	EdgeAuthorizer  string `optional:"" help:"id of edge authorizer connection used to register with the Aserto control plane"`
	DecisionLogging bool   `optional:"" help:"enable decision logging"`
}

func (cmd ConfigureCmd) Validate() error {
	if cmd.DecisionLogging && cmd.EdgeAuthorizer == "" {
		return errors.New("decision logging requires an edge authorizer to be configured")
	}

	return nil
}

func (cmd ConfigureCmd) Run(c *cc.CommonCtx) error {
	fmt.Fprintf(c.UI.Err(), ">>> configure policy...\n")
	fmt.Fprintf(c.UI.Err(), "tenant id: %s\n", c.TenantID())

	configDir, err := CreateConfigDir()
	if err != nil {
		return err
	}

	client, err := c.TenantClient()
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
		TenantID:        c.TenantID(),
		PolicyName:      policyRef.Name,
		PolicyID:        policyRef.Id,
		DiscoveryURL:    discoveryConf.URL,
		TenantKey:       discoveryConf.APIKey,
		DecisionLogging: cmd.DecisionLogging,
	}

	if cmd.EdgeAuthorizer != "" {
		certFile, keyFile, errCerts := getEdgeAuthorizerCerts(c.Context, client, cmd.EdgeAuthorizer, configDir)
		if errCerts != nil {
			return err
		}

		params.ControlPlane.Enabled = true
		params.ControlPlane.Address = c.Environment.Get(x.ControlPlaneService).Address
		params.ControlPlane.ClientCertPath = path.Join("/app/cfg", certFile)
		params.ControlPlane.ClientKeyPath = path.Join("/app/cfg", keyFile)

		//params.DecisionLogger.EMSAddress = c.DecisionLogger.EMSAddress
		params.DecisionLogger.EMSAddress = c.Environment.Get(x.EMSService).Address
		params.DecisionLogger.StorePath = decisionlogger.ContainerPath
		params.DecisionLogger.ClientCertPath = path.Join("/app/cfg", certFile)
		params.DecisionLogger.ClientKeyPath = path.Join("/app/cfg", keyFile)
	}

	if params.TenantKey == "" {
		return errors.Errorf("missing $ASERTO_TENANT_KEY env var")
	}

	fmt.Fprintf(c.UI.Err(), "policy id: %s\n", params.PolicyID)
	fmt.Fprintf(c.UI.Err(), "policy name: %s\n", params.PolicyName)

	var w io.Writer

	if cmd.Stdout {
		w = c.UI.Output()
	} else {
		w, err = os.Create(path.Join(configDir, params.PolicyName+".yaml"))
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

func getEdgeAuthorizerCerts(ctx context.Context, client *tenant.Client, connID, configDir string) (certFile, keyFile string, err error) {
	resp, err := client.Connections.GetConnection(ctx, &connection.GetConnectionRequest{
		Id: connID,
	})
	if err != nil {
		return "", "", err
	}

	conn := resp.Result
	if conn == nil {
		return "", "", errors.New("invalid empty connection")
	}

	if conn.Kind != api.ProviderKind_PROVIDER_KIND_EDGE_AUTHORIZER {
		return "", "", errors.New("not an edge authorizer connection")
	}

	certs := conn.Config.Fields["api_cert"].GetListValue().GetValues()
	if len(certs) == 0 {
		return "", "", errors.New("invalid configuration: api_cert")
	}

	structVal := certs[len(certs)-1].GetStructValue()
	if structVal == nil {
		return "", "", errors.New("invalid configuration: api_cert")
	}

	err = fileFromConfigField(structVal, "certificate", configDir, "client.crt")
	if err != nil {
		return "", "", err
	}

	err = fileFromConfigField(structVal, "private_key", configDir, "client.key")
	if err != nil {
		return "", "", err
	}

	return "client.crt", "client.key", nil
}

func fileFromConfigField(structVal *structpb.Struct, field, configDir, fileName string) error {
	val, ok := structVal.Fields[field]
	if !ok {
		return errors.Errorf("missing field: %s", field)
	}

	strVal := val.GetStringValue()
	if strVal == "" {
		return errors.Errorf("empty field: %s", field)
	}

	filePath := path.Join(configDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(strVal)
	if err != nil {
		return err
	}

	return nil
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
