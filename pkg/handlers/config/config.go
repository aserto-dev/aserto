package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aserto-dev/aserto/pkg/cc"
	config "github.com/aserto-dev/aserto/pkg/cc/config"
	aErr "github.com/aserto-dev/aserto/pkg/cc/errors"
	"github.com/aserto-dev/aserto/pkg/handlers/user"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/aserto/pkg/keyring"
	"github.com/aserto-dev/aserto/pkg/x"
	"github.com/aserto-dev/clui"
	"github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/go-grpc/aserto/tenant/account/v1"
	"github.com/samber/lo"
	"gopkg.in/yaml.v2"

	"github.com/pkg/errors"
)

type GetContextsCmd struct{}

func (cmd *GetContextsCmd) Run(c *cc.CommonCtx) error {
	cfgPath, err := config.GetSymlinkConfigPath()
	if err != nil {
		return err
	}

	if !config.FileExists(cfgPath) {
		return aErr.NeedLoginErr
	}

	cfg, err := config.NewConfig(config.Path(cfgPath))
	if err != nil {
		return err
	}

	return jsonx.OutputJSON(c.UI.Output(), cfg.Context.Contexts)
}

type GetActiveContextCmd struct{}

func (cmd *GetActiveContextCmd) Run(c *cc.CommonCtx) error {
	cfgPath, err := config.GetSymlinkConfigPath()
	if err != nil {
		return err
	}

	if !config.FileExists(cfgPath) {
		return aErr.NeedLoginErr
	}

	cfg, err := config.NewConfig(config.Path(cfgPath))
	if err != nil {
		return err
	}

	for _, ctxs := range cfg.Context.Contexts {
		if cfg.Context.ActiveContext == ctxs.Name {
			return jsonx.OutputJSON(c.UI.Output(), ctxs)
		}
	}

	return nil
}

type DeleteContextCmd struct {
	ContextName string `arg:"" name:"context_name" required:"" help:"context name"`
}

func (cmd *DeleteContextCmd) Run(c *cc.CommonCtx) error {
	cfgPath, err := config.GetSymlinkConfigPath()
	if err != nil {
		return err
	}

	if !config.FileExists(cfgPath) {
		return aErr.NeedLoginErr
	}

	currentUserFilePath, err := os.Readlink(cfgPath)
	if err != nil {
		return err
	}

	cfg, err := config.NewConfig(config.Path(currentUserFilePath))
	if err != nil {
		return err
	}

	if cfg.Context.ActiveContext == cmd.ContextName {
		return errors.Errorf("the context in use cannot be deleted")
	}

	for index, ctxs := range cfg.Context.Contexts {
		if ctxs.Name == cmd.ContextName {
			cfg.Context.Contexts = append(cfg.Context.Contexts[:index], cfg.Context.Contexts[index+1:]...)
			break
		}
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(currentUserFilePath, data, 0600)
}

type SetContextCmd struct {
	Context  string `arg:""  type:"existingfile" name:"context" optional:"" help:"file path to context or '-' to read from stdin"`
	Template bool   `name:"template" help:"prints a context template on stdout"`
	Force    bool   `name:"force" help:"if a context wit the same name exists, forces overwrite"`
}

func (cmd *SetContextCmd) Run(c *cc.CommonCtx) error {
	if cmd.Template {
		return printContext(c.UI)
	}

	if cmd.Context == "" {
		return errors.New("context argument is required")
	}

	var req config.Ctx
	if cmd.Context == "-" {
		decoder := json.NewDecoder(os.Stdin)

		err := decoder.Decode(&req)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal request from stdin")
		}
	} else {
		dat, err := os.ReadFile(cmd.Context)
		if err != nil {
			return errors.Wrapf(err, "opening file [%s]", cmd.Context)
		}

		err = json.Unmarshal(dat, &req)
		if err != nil {
			return errors.Wrapf(err, "failed to unmarshal request from file [%s]", cmd.Context)
		}
	}

	cfgPath, err := config.GetSymlinkConfigPath()
	if err != nil {
		return err
	}

	if !config.FileExists(cfgPath) {
		return aErr.NeedLoginErr
	}

	currentUserFilePath, err := os.Readlink(cfgPath)
	if err != nil {
		return err
	}

	cfg, err := config.NewConfig(config.Path(currentUserFilePath))
	if err != nil {
		return err
	}

	_, idx, _ := lo.FindIndexOf(cfg.Context.Contexts, func(c config.Ctx) bool { return c.Name == req.Name })

	if idx > -1 {
		if !cmd.Force {
			c.UI.Exclamation().Msg("A context with this name already exists; please choose another name or use --force flag to overwrite the existing one")
			c.UI.Note().Msg("Aborting...")
			return nil
		}

		cfg.Context.Contexts[idx] = req
	} else {
		cfg.Context.Contexts = append(cfg.Context.Contexts, req)
	}

	if req.TenantID != "" && c.Factory.TenantID() != "" {
		err = changeTokenToTenantID(c, req.TenantID)
		if err != nil {
			return err
		}
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(currentUserFilePath, data, 0600)
}

type UseContextCmd struct {
	ContextName string `arg:"" name:"context_name" required:"" help:"context name"`
}

func (cmd *UseContextCmd) Run(c *cc.CommonCtx) error {
	cfgPath, err := config.GetSymlinkConfigPath()
	if err != nil {
		return err
	}

	if !config.FileExists(cfgPath) {
		return aErr.NeedLoginErr
	}

	currentUserFilePath, err := os.Readlink(cfgPath)
	if err != nil {
		return err
	}

	cfg, err := config.NewConfig(config.Path(currentUserFilePath))
	if err != nil {
		return err
	}

	_, found := lo.Find(cfg.Context.Contexts, func(c config.Ctx) bool { return c.Name == cmd.ContextName })

	if !found {
		return errors.Errorf("the context name provided doesn't exists")
	}

	cfg.Context.ActiveContext = cmd.ContextName

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(currentUserFilePath, data, 0600)
}

func changeTokenToTenantID(c *cc.CommonCtx, tenantID string) error {
	conn, err := c.TenantClient()
	if err != nil {
		return err
	}

	getAccntResp, err := conn.Account.GetAccount(c.Context, &account.GetAccountRequest{})
	if err != nil {
		return errors.Wrapf(err, "get account")
	}

	var tnt *api.Tenant
	for _, t := range getAccntResp.Result.Tenants {
		if t.Id == tenantID {
			tnt = t
			break
		}
	}

	if tnt == nil {
		return errors.Errorf("tenant id does not exist in users tenant collection [%s]", tenantID)
	}

	fmt.Fprintf(c.UI.Err(), "tenant %s - %s\n", tnt.Id, tnt.Name)

	tkn, err := c.CachedToken.Get()
	if err != nil {
		return err
	}

	tenantKr, err := keyring.NewTenantKeyRing(tenantID + "-" + tkn.Subject)
	if err != nil {
		return err
	}

	token, err := tenantKr.GetToken()
	if err == nil && token != nil {
		// token already set
		return nil
	}

	tenantToken := &keyring.TenantToken{TenantID: tenantID}

	if err = user.GetConnectionKeys(c.Context, conn, tenantToken); err != nil {
		return errors.Wrapf(err, "get connection keys")
	}

	if err := tenantKr.SetToken(tenantToken); err != nil {
		return err
	}

	return nil
}

func printContext(ui *clui.UI) error {
	req := config.Ctx{
		Name:     "context_name",
		TenantID: "tenant_id",
		AuthorizerService: x.ServiceOptions{
			Address:  "address:port",
			APIKey:   "key",
			Insecure: true,
		},
	}
	return jsonx.OutputJSON(ui.Output(), req)
}
