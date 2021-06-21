package cmd

type Globals struct {
	Verbose            bool   `help:"verbose output"`
	AuthorizerOverride string `name:"authorizer" env:"ASERTO_AUTHORIZER" help:"authorizer override"`
	TenantOverride     string `name:"tenant-id" env:"ASERTO_TENANT_ID" help:"tenant id override"`
	Environment        string `name:"env" default:"prod" env:"ASERTO_ENV" hidden:""`
}
