package paths

import (
	"os"
	"path"

	"github.com/aserto-dev/aserto/pkg/filex"
	"github.com/pkg/errors"
)

const (
	DefaultConfigRoot = ".config/aserto/sidecar"
	DefaultCacheRoot  = ".cache/aserto/sidecar"

	cfgSubdir  = "cfg"
	certSubdir = "certs"
	edsSubdir  = "eds"

	localConfigFile = "local.yaml"
)

type Certs struct {
	Root    string
	GRPC    *CertPaths
	Gateway *CertPaths
}

type Paths struct {
	Config string
	EDS    string
	Data   string

	Certs Certs
}

func (p *Paths) LocalConfig() string {
	return path.Join(p.Config, localConfigFile)
}

func New() (*Paths, error) {
	return newPaths("")
}

func NewWithDataRoot(dataRoot string) (*Paths, error) {
	return newPaths(dataRoot)
}

func newPaths(dataRoot string) (*Paths, error) {
	confRoot, cacheRoot, err := DefaultRoots()
	if err != nil {
		return nil, err
	}

	if dataRoot != "" && !filex.DirExists(dataRoot) {
		return nil, errors.Wrapf(os.ErrNotExist, "directory '%s' does not exist", dataRoot)
	}

	return NewIn(confRoot, cacheRoot, dataRoot), nil
}

func DefaultRoots() (confRoot, cacheRoot string, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}

	return path.Join(home, DefaultConfigRoot), path.Join(home, DefaultCacheRoot), nil
}

func NewIn(confRoot, cacheRoot, dataRoot string) *Paths {
	certDir := CertDir(confRoot)

	return &Paths{
		Config: ConfigDir(confRoot),
		EDS:    EdsDir(cacheRoot),
		Data:   dataRoot,

		Certs: Certs{
			Root:    certDir,
			GRPC:    GRPCCerts(certDir),
			Gateway: GatewayCerts(certDir),
		},
	}
}

func Create() (*Paths, error) {
	confRoot, cacheRoot, err := DefaultRoots()
	if err != nil {
		return nil, err
	}

	return CreateIn(confRoot, cacheRoot)
}

func CreateIn(confRoot, cacheRoot string) (*Paths, error) {
	paths := NewIn(confRoot, cacheRoot, "")
	for _, confDir := range []string{
		paths.Config,
		paths.Certs.Root,
		paths.EDS,
	} {
		if err := createDir(confDir); err != nil {
			return nil, errors.Wrap(err, confDir)
		}
	}

	return paths, nil
}

func ConfigDir(confRoot string) string {
	return path.Join(confRoot, cfgSubdir)
}

func CertDir(confRoot string) string {
	return path.Join(confRoot, certSubdir)
}

func EdsDir(cacheRoot string) string {
	return path.Join(cacheRoot, edsSubdir)
}

func GRPCCerts(certDir string) *CertPaths {
	return NewCertPaths(certDir, "grpc")
}

func GatewayCerts(certsDir string) *CertPaths {
	return NewCertPaths(certsDir, "gateway")
}

func createDir(dir string) error {
	if !filex.DirExists(dir) {
		return os.MkdirAll(dir, 0700)
	}

	return nil
}
