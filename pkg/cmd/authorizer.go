package cmd

import (
	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/cc/config"
	"github.com/go-http-utils/headers"
	"google.golang.org/grpc/metadata"

	topazCC "github.com/aserto-dev/topaz/pkg/cli/cc"
	topazClients "github.com/aserto-dev/topaz/pkg/cli/clients"
	topazAuthz "github.com/aserto-dev/topaz/pkg/cli/cmd/authorizer"
)

const (
	BearerToken = `Bearer `
)

type AuthorizerCmd struct {
	EvalDecision topazAuthz.EvalCmd         `cmd:"" help:"evaluate policy decision" group:"authorizer"`
	DecisionTree topazAuthz.DecisionTreeCmd `cmd:"" help:"get decision tree" group:"authorizer"`
	ExecQuery    topazAuthz.QueryCmd        `cmd:"" help:"execute query" group:"authorizer"`
	GetPolicy    topazAuthz.GetPolicyCmd    `cmd:"" help:"get policy" group:"authorizer"`
	ListPolicies topazAuthz.ListPoliciesCmd `cmd:"" help:"list policies" group:"authorizer"`
}

func (cmd *AuthorizerCmd) AfterApply(context *kong.Context, c *topazCC.CommonCtx) error {
	var cfg *config.Config
	var err error

	allFlags := context.Flags()
	for _, f := range allFlags {
		if f.Name == "config" {
			configPath := context.FlagValue(f).(string)
			cfg, err = config.NewConfig(config.Path(configPath))
			if err != nil {
				return err
			}
		}
	}

	if cfg.TopazConfigFile != "" {
		err = setServicesConfig(cfg, cfg.TopazConfigFile)
		if err != nil {
			return err
		}
	}

	tenantToken, err := getTenantTokenDetails(cfg.Auth)
	if err != nil {
		return err
	}

	authorizerConfig := topazClients.AuthorizerConfig{
		Host:     cfg.Services.AuthorizerService.Address,
		APIKey:   cfg.Services.AuthorizerService.APIKey,
		Token:    tenantToken,
		Insecure: cfg.Services.AuthorizerService.Insecure,
		TenantID: cfg.TenantID,
	}

	c.Context = metadata.AppendToOutgoingContext(c.Context, string(headers.Authorization), BearerToken+tenantToken)

	cmd.EvalDecision.AuthorizerConfig = authorizerConfig
	cmd.ExecQuery.AuthorizerConfig = authorizerConfig
	cmd.GetPolicy.AuthorizerConfig = authorizerConfig
	cmd.DecisionTree.AuthorizerConfig = authorizerConfig
	cmd.ListPolicies.AuthorizerConfig = authorizerConfig

	return nil
}

func (cmd *AuthorizerCmd) Run(c *cc.CommonCtx) error {
	return nil
}
