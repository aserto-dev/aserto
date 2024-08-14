package controlplane

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/go-grpc/aserto/tenant/connection/v1"
	"github.com/aserto-dev/topaz/pkg/cli/jsonx"

	"github.com/pkg/errors"
)

type ClientCertCmd struct {
	ID        string `arg:"" help:"edge authorizer connection ID"`
	Directory string `flag:"" short:"d" help:"directory to save certificate file" default:"${cwd}"`
	Raw       bool   `flag:"" help:"raw message output"`
}

func (cmd ClientCertCmd) Run(c *cc.CommonCtx) error {
	if cmd.ID == "" {
		return errors.New("connection ID argument not provided")
	}

	cli, err := c.TenantClient(c.Context)
	if err != nil {
		return err
	}

	resp, err := cli.Connections.GetConnection(c.Context, &connection.GetConnectionRequest{
		Id: cmd.ID,
	})
	if err != nil {
		return err
	}

	conn := resp.Result
	if conn == nil {
		return errors.New("invalid empty connection")
	}

	if conn.Kind != api.ProviderKind_PROVIDER_KIND_EDGE_AUTHORIZER {
		return errors.New("not an edge authorizer connection")
	}

	var cfg Config
	if buf, err := conn.Config.MarshalJSON(); err == nil {
		if err := json.Unmarshal(buf, &cfg); err != nil {
			return err
		}
	} else {
		return err
	}

	if len(cfg.APICerts) == 0 {
		return errors.New("invalid connection configuration")
	}

	cert := cfg.APICerts[len(cfg.APICerts)-1]

	if cmd.Raw {
		return jsonx.OutputJSON(c.StdOut(), cert)
	}

	c.Con().Info().Msg("ID : %s", cert.ID)
	c.Con().Info().Msg("CN : %s", cert.CN)
	c.Con().Info().Msg("Exp: %s", cert.Expiration)

	if err := cmd.writeFile(c, sidecarKey, cert.Key); err != nil {
		return err
	}

	if err := cmd.writeFile(c, sidecarCert, cert.Cert); err != nil {
		return err
	}

	return nil
}

const (
	sidecarKey  string = "sidecar.key"
	sidecarCert string = "sidecar.crt"
)

type Config struct {
	APICerts []APICert `json:"api_cert"`
}

type APICert struct {
	ID         string    `json:"id"`
	CN         string    `json:"common_name"`
	Key        string    `json:"private_key"`
	Cert       string    `json:"certificate"`
	Expiration time.Time `json:"expiration"`
}

func (cmd ClientCertCmd) writeFile(c *cc.CommonCtx, name, value string) error {
	if fi, err := os.Stat(cmd.Directory); err != nil || !fi.IsDir() {
		if err != nil {
			return err
		}
		if !fi.IsDir() {
			return errors.Errorf("--directory argument %q is not a directory", cmd.Directory)
		}
	}

	fn := filepath.Join(cmd.Directory, name)

	c.Con().Msg(fn)

	w, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = fmt.Fprintln(w, value)

	return err
}
