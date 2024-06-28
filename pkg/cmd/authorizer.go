package cmd

import (
	"errors"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/go-http-utils/headers"
	"github.com/samber/lo"
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

	isTopazConfig := !cc.IsAsertoAccount(cfg.ConfigName)

	if isTopazConfig {
		err = setServicesConfig(cfg, c.Config.Active.ConfigFile)
		if err != nil {
			return err
		}
	}

	token, err := getTenantToken((cfg.Auth))
	if err != nil && !errors.Is(err, errs.NeedLoginErr) {
		return err
	}

	authorizerConfig := topazClients.AuthorizerConfig{
		Host:     cfg.Services.AuthorizerService.Address,
		APIKey:   cfg.Services.AuthorizerService.APIKey,
		Token:    lo.Ternary(isTopazConfig, "", token.Access), // only send access token to hosted services.
		Insecure: cfg.Services.AuthorizerService.Insecure,
		TenantID: lo.Ternary(isTopazConfig, "", cfg.TenantID),
	}

	if !isTopazConfig {
		// only send access token to hosted services.
		c.Context = metadata.AppendToOutgoingContext(c.Context, string(headers.Authorization), "Bearer "+token.Access)
	}

	cmd.AuthorizerCmd.CheckDecision.AuthorizerConfig = authorizerConfig
	cmd.AuthorizerCmd.ExecQuery.AuthorizerConfig = authorizerConfig
	cmd.AuthorizerCmd.DecisionTree.AuthorizerConfig = authorizerConfig
	cmd.AuthorizerCmd.Get.Policy.AuthorizerConfig = authorizerConfig
	cmd.AuthorizerCmd.List.Policies.AuthorizerConfig = authorizerConfig
	cmd.AuthorizerCmd.Test.Exec.AuthorizerConfig = authorizerConfig

	return nil
}

func (cmd *AuthorizerCmd) Run(c *cc.CommonCtx) error {
	return nil
}
