package dev

import (
	"fmt"
	"io"
	"os"
	"path"
	"text/template"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/filex"
	"github.com/aserto-dev/aserto/pkg/grpcc"
	"github.com/aserto-dev/aserto/pkg/grpcc/tenant"

	"github.com/aserto-dev/proto/aserto/api"
	"github.com/aserto-dev/proto/aserto/tenant/connection"
	"github.com/aserto-dev/proto/aserto/tenant/policy"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

type ConfigureCmd struct {
	Name string `arg:"" required:"" help:"policy name"`
}

// nolint:funlen // tbd
func (cmd ConfigureCmd) Run(c *cc.CommonCtx) error {
	color.Green(">>> configure policy...")
	params := templateParams{}

	conn, err := tenant.Connection(
		c.Context,
		c.TenantService(),
		grpcc.NewTokenAuth(c.AccessToken()),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.OutWriter, "tenant id: %s\n", c.TenantID())

	ctx := grpcc.SetTenantContext(c.Context, c.TenantID())

	policyClient := conn.PolicyClient()
	policyRefResp, err := policyClient.ListPolicyRefs(
		ctx,
		&policy.ListPolicyRefsRequest{},
	)
	if err != nil {
		return err
	}

	var (
		pack  *api.PolicyRef
		found bool
	)

	for _, v := range policyRefResp.Results {
		if v.Name == cmd.Name {
			pack = v
			found = true
			break
		}
	}
	if !found {
		return errors.Errorf("policy not found [%s]", cmd.Name)
	}

	params.TenantID = c.TenantID()
	params.PolicyName = pack.Name
	params.PolicyID = pack.Id
	params.RegistrySvc = c.RegistrySvc()

	fmt.Fprintf(c.OutWriter, "policy id: %s\n", params.PolicyID)

	connClient := conn.ConnectionManagerClient()
	listResp, err := connClient.ListConnections(
		ctx,
		&connection.ListConnectionsRequest{
			Kind: api.ProviderKind_POLICY_REGISTRY,
		},
	)
	if err != nil {
		return err
	}
	if len(listResp.Results) != 1 {
		return errors.Errorf("policy registry connection not found")
	}

	connResp, err := connClient.GetConnection(
		ctx,
		&connection.GetConnectionRequest{
			Id: listResp.Results[0].Id,
		},
	)
	if err != nil {
		return err
	}

	connConfigMap := connResp.Result.Config.AsMap()
	downloadKey, ok := connConfigMap["download_key"].(string)
	if !ok {
		return errors.Errorf("download key not found")
	}

	params.DownloadAPIKey = downloadKey

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := path.Join(home, "/.config/aserto/aserto-one/cfg")
	if !filex.DirExists(configDir) {
		if err := os.MkdirAll(configDir, 0700); err != nil {
			return err
		}
	}

	w, err := os.Create(path.Join(configDir, params.PolicyName+".yaml"))
	if err != nil {
		return err
	}

	err = CreateConfig(w, &params)

	return err
}

func CreateConfig(w io.Writer, params *templateParams) error {
	t, err := template.New("config").Parse(configTemplate)
	if err != nil {
		return err
	}

	err = t.Execute(w, params)
	if err != nil {
		return err
	}

	return nil
}

type templateParams struct {
	TenantID       string
	DownloadAPIKey string
	PolicyName     string
	PolicyID       string
	RegistrySvc    string
}

const configTemplate = `
---
logging:
  prod: false
  log_level: debug

directory_service:
  path: "/app/eds/eds-{{ .PolicyName }}.db"

api:
  grpc:
    connection_timeout_seconds: 2
    certs:
      tls_key_path: "/root/.config/aserto/aserto-one/certs/grpc.key"
      tls_cert_path: "/root/.config/aserto/aserto-one/certs/grpc.crt"
      tls_ca_cert_path: "/root/.config/aserto/aserto-one/certs/grpc-ca.crt"
  gateway:
    certs:
      tls_key_path: "/root/.config/aserto/aserto-one/certs/gateway.key"
      tls_cert_path: "/root/.config/aserto/aserto-one/certs/gateway.crt"
      tls_ca_cert_path: "/root/.config/aserto/aserto-one/certs/gateway-ca.crt"

opa:
  instance_id: "{{ .TenantID }}"
  store: aserto
  graceful_shutdown_period_seconds: 2
  config:
    services:
      acmecorp:
        url: {{ .RegistrySvc }}/{{ .TenantID }}
        response_header_timeout_seconds: 5
        credentials:
          bearer:
            token: "{{ .DownloadAPIKey }}"
    bundles:
      {{ .PolicyID }}:
        service: acmecorp
        resource: "/{{ .PolicyID }}/bundle.tar.gz"
        persist: true
        polling:
          min_delay_seconds: 10
          max_delay_seconds: 30
`
