package cmd

import (
	"errors"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/go-http-utils/headers"
	"google.golang.org/grpc/metadata"

	topazCC "github.com/aserto-dev/topaz/pkg/cli/cc"
	topazClients "github.com/aserto-dev/topaz/pkg/cli/clients"
	"github.com/aserto-dev/topaz/pkg/cli/cmd/authorizer"

	errs "github.com/aserto-dev/aserto/pkg/cc/errors"
)

const (
	BearerToken = `Bearer `
	ConfigFlag  = `config`
)

type AuthorizerCmd struct {
	authorizer.AuthorizerCmd
}

func (cmd *AuthorizerCmd) AfterApply(context *kong.Context, c *topazCC.CommonCtx) error {
	cfg, err := getConfig(context)
	if err != nil {
		return err
	}

	if !cc.IsAsertoAccount(cfg.ConfigName) {
		err = setServicesConfig(cfg, c.Config.Active.ConfigFile)
		if err != nil {
			return err
		}
	}

	tenantToken, err := getTenantTokenDetails(cfg.Auth)
	if err != nil && !errors.Is(err, errs.NeedLoginErr) {
		return err
	}
	useTenantID := ""
	if tenantToken == "" {
		useTenantID = cfg.TenantID
	}

	authorizerConfig := topazClients.AuthorizerConfig{
		Host:     cfg.Services.AuthorizerService.Address,
		APIKey:   cfg.Services.AuthorizerService.APIKey,
		Token:    tenantToken,
		Insecure: cfg.Services.AuthorizerService.Insecure,
		TenantID: useTenantID,
	}

	c.Context = metadata.AppendToOutgoingContext(c.Context, string(headers.Authorization), BearerToken+tenantToken)

	cmd.AuthorizerCmd.CheckDecision.AuthorizerConfig = authorizerConfig
	cmd.AuthorizerCmd.ExecQuery.AuthorizerConfig = authorizerConfig
	cmd.AuthorizerCmd.GetPolicy.AuthorizerConfig = authorizerConfig
	cmd.AuthorizerCmd.DecisionTree.AuthorizerConfig = authorizerConfig
	cmd.AuthorizerCmd.ListPolicies.AuthorizerConfig = authorizerConfig

	return nil
}

func (cmd *AuthorizerCmd) Run(c *cc.CommonCtx) error {
	return nil
}
