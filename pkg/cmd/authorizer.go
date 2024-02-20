package cmd

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/cc/config"
	"github.com/aserto-dev/aserto/pkg/filex"
	"github.com/go-http-utils/headers"
	"google.golang.org/grpc/metadata"

	aErr "github.com/aserto-dev/aserto/pkg/cc/errors"
	topazCC "github.com/aserto-dev/topaz/pkg/cli/cc"
	topazClients "github.com/aserto-dev/topaz/pkg/cli/clients"
	topazAuthz "github.com/aserto-dev/topaz/pkg/cli/cmd/authorizer"
)

type AuthorizerCmd struct {
	EvalDecision topazAuthz.EvalDecisionCmd `cmd:"" help:"evaluate policy decision" group:"authorizer"`
	DecisionTree topazAuthz.DecisionTreeCmd `cmd:"" help:"get decision tree" group:"authorizer"`
	ExecQuery    topazAuthz.ExecQueryCmd    `cmd:"" help:"execute query" group:"authorizer"`
	Compile      topazAuthz.CompileCmd      `cmd:"" help:"compile query" group:"authorizer"`
	GetPolicy    topazAuthz.GetPolicyCmd    `cmd:"" help:"get policy" group:"authorizer"`
	ListPolicies topazAuthz.ListPoliciesCmd `cmd:"" help:"list policies" group:"authorizer"`
}

func (cmd *AuthorizerCmd) AfterApply(c *topazCC.CommonCtx) error {
	cfgPath, err := config.GetSymlinkConfigPath()
	if err != nil {
		return err
	}

	if !filex.FileExists(cfgPath) {
		return aErr.NeedLoginErr
	}

	cfg, err := config.NewConfig(config.Path(cfgPath))
	if err != nil {
		return err
	}

	for _, ctxs := range cfg.Context.Contexts {
		if cfg.Context.ActiveContext == ctxs.Name {
			if ctxs.TopazConfigFile != "" {
				err = setServicesConfig(cfg, ctxs.TopazConfigFile)
				if err != nil {
					return err
				}
			}
			// err = os.Setenv(topazClients.EnvTopazAuthorizerSvc, cfg.Services.AuthorizerService.Address)
			// if err != nil {
			// 	return err
			// }
			tenantToken, err := getTenantTokenDetails(ctxs.TenantID, cfg.Auth)
			if err != nil {
				return err
			}

			authorizerConfig := topazClients.AuthorizerConfig{
				Host:     cfg.Services.AuthorizerService.Address,
				APIKey:   "",
				Token:    tenantToken,
				Insecure: cfg.Services.AuthorizerService.Insecure,
				TenantID: ctxs.TenantID,
			}

			c.Context = metadata.AppendToOutgoingContext(c.Context, string(headers.Authorization), "Bearer "+tenantToken)

			cmd.EvalDecision.AuthorizerConfig = authorizerConfig
			cmd.Compile.AuthorizerConfig = authorizerConfig
			cmd.ExecQuery.AuthorizerConfig = authorizerConfig
			cmd.GetPolicy.AuthorizerConfig = authorizerConfig
			cmd.DecisionTree.AuthorizerConfig = authorizerConfig
			cmd.ListPolicies.AuthorizerConfig = authorizerConfig
		}
	}
	return nil
}

func (cmd *AuthorizerCmd) Run(c *cc.CommonCtx) error {
	return nil
}
