package dev

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/clients/tenant"
	decisionlogger "github.com/aserto-dev/aserto/pkg/decision_logger"
	"github.com/aserto-dev/aserto/pkg/x"
	"google.golang.org/protobuf/types/known/structpb"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/go-grpc/aserto/tenant/connection/v1"
	policy "github.com/aserto-dev/go-grpc/aserto/tenant/policy/v1"

	topazConfig "github.com/aserto-dev/topaz/pkg/cc/config"
	topazCC "github.com/aserto-dev/topaz/pkg/cli/cc"
	topazCerts "github.com/aserto-dev/topaz/pkg/cli/cmd/certs"
	topazCommon "github.com/aserto-dev/topaz/pkg/cli/cmd/common"
	topazConfigure "github.com/aserto-dev/topaz/pkg/cli/cmd/configure"
	"github.com/pkg/errors"
)

type ConfigureCmd struct {
	topazConfigure.NewConfigCmd
	EdgeAuthorizer  string `optional:"" help:"id of edge authorizer connection used to register with the Aserto control plane"`
	DecisionLogging bool   `optional:"" help:"enable decision logging"`
}

func (cmd *ConfigureCmd) Validate() error {
	if cmd.DecisionLogging && cmd.EdgeAuthorizer == "" {
		return errors.New("decision logging requires an edge authorizer to be configured")
	}

	return nil
}

func (cmd *ConfigureCmd) Run(c *cc.CommonCtx) error {
	c.Con().Info().Msg(">>> configure policy...")
	c.Con().Msg("tenant id: %s", c.TenantID())

	if cmd.Name == "" && cmd.Resource == "" {
		if cmd.LocalPolicyImage == "" {
			return errors.New("you either need to provide a local policy image or the resource and the policy name for the configuration")
		}
	}
	configFile := cmd.Name.String() + ".yaml"
	if configFile != c.CommonCtx.Config.Active.ConfigFile {
		c.CommonCtx.Config.Active.Config = cmd.Name.String()
		c.CommonCtx.Config.Active.ConfigFile = filepath.Join(topazCC.GetTopazCfgDir(), configFile)
	}

	configGenerator := topazConfig.NewGenerator(cmd.Name.String()).
		WithVersion(topazConfig.ConfigFileVersion).
		WithLocalPolicyImage(cmd.LocalPolicyImage).
		WithPolicyName(cmd.Name.String()).
		WithResource(cmd.Resource).
		WithEdgeDirectory(cmd.EdgeDirectory).
		WithTenantID(c.TenantID())

	_, err := configGenerator.CreateConfigDir()
	if err != nil {
		return err
	}

	if _, err := configGenerator.CreateCertsDir(); err != nil {
		return err
	}
	certGenerator := topazCerts.GenerateCertsCmd{CertsDir: topazCC.GetTopazCertsDir()}
	err = certGenerator.Run(c.CommonCtx)
	if err != nil {
		return err
	}
	if _, err := configGenerator.CreateDataDir(); err != nil {
		return err
	}

	client, err := c.TenantClient(c.Context)
	if err != nil {
		return err
	}
	getDiscovery := true
	policyRef, err := findPolicyRef(c.Context, client, cmd.Name.String())
	if err != nil {
		// policy name not found
		getDiscovery = false
	}
	if getDiscovery {
		discoveryConf, err := getDiscoveryConfig(c.Context, client)
		if err != nil {
			return err
		}
		configGenerator = configGenerator.WithDiscovery(discoveryConf.URL, discoveryConf.APIKey)
	}

	if cmd.EdgeAuthorizer != "" {
		certFile, keyFile, errCerts := getEdgeAuthorizerCerts(c.Context, client, cmd.EdgeAuthorizer, topazCC.GetTopazCertsDir(), policyRef.Name)
		if errCerts != nil {
			return err
		}
		configGenerator = configGenerator.
			WithController(c.Environment.Get(x.ControlPlaneService).Address,
				filepath.Join("${TOPAZ_CERTS_DIR}", certFile),
				filepath.Join("${TOPAZ_CERTS_DIR}", keyFile)).
			WithSelfDecisionLogger(c.Environment.Get(x.EMSService).Address,
				filepath.Join("${TOPAZ_CERTS_DIR}", certFile),
				filepath.Join("${TOPAZ_CERTS_DIR}", keyFile),
				filepath.Join(cmd.Name.String(), decisionlogger.Dir),
			)
	}

	c.Con().Msg("policy name: %s", cmd.Name)

	var w io.Writer

	if cmd.Stdout {
		w = c.StdOut()
	} else {
		if !cmd.Force {
			if _, err := os.Stat(c.CommonCtx.Config.Active.ConfigFile); err == nil {
				c.Con().Warn().Msg("Configuration file %q already exists.", c.CommonCtx.Config.Active.ConfigFile)
				if !topazCommon.PromptYesNo("Do you want to continue?", false) {
					return nil
				}
			}
		}
		w, err = os.Create(c.CommonCtx.Config.Active.ConfigFile)
		if err != nil {
			return err
		}
	}
	if configGenerator.DiscoveryURL != "" {
		return configGenerator.GenerateConfig(w, topazConfig.EdgeTemplate)
	}
	return configGenerator.GenerateConfig(w, topazConfig.Template)
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

func getEdgeAuthorizerCerts(ctx context.Context, client *tenant.Client, connID, configDir, policyName string) (certFile, keyFile string, err error) {
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

	crtName := fmt.Sprintf("%s-client.crt", policyName)
	keyName := fmt.Sprintf("%s-client.key", policyName)

	err = fileFromConfigField(structVal, "certificate", configDir, crtName)
	if err != nil {
		return "", "", err
	}

	err = fileFromConfigField(structVal, "private_key", configDir, keyName)
	if err != nil {
		return "", "", err
	}

	return crtName, keyName, nil
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

	filePath := filepath.Join(configDir, fileName)
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
