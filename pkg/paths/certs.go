package paths

import (
	"path"

	"github.com/aserto-dev/aserto/pkg/filex"
)

type CertPaths struct {
	Name string
	Cert string
	CA   string
	Key  string
	Dir  string
}

func NewCertPaths(dir, prefix string) *CertPaths {
	return &CertPaths{
		Name: "authorizer-" + prefix,
		Cert: path.Join(dir, prefix+".crt"),
		CA:   path.Join(dir, prefix+"-ca.crt"),
		Key:  path.Join(dir, prefix+".key"),
		Dir:  dir,
	}
}

func (c *CertPaths) FindExisting() []string {
	existing := []string{}
	for _, cert := range []string{c.Cert, c.CA, c.Key} {
		if filex.FileExists(cert) {
			existing = append(existing, cert)
		}
	}

	return existing
}
