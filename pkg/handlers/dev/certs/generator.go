package certs

import (
	"fmt"
	"io"

	"github.com/aserto-dev/aserto/pkg/paths"
	"github.com/aserto-dev/go-utils/certs"
	"github.com/aserto-dev/go-utils/logger"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func GenerateCerts(logOut, errOut io.Writer, certPaths ...*paths.CertPaths) error {
	existingFiles := []string{}

	for _, cert := range certPaths {
		existingFiles = append(existingFiles, cert.FindExisting()...)
	}

	if len(existingFiles) != 0 {
		fmt.Fprintln(logOut, "Some cert files already exist. Skipping generation.", existingFiles)
		return nil
	}

	return generate(logOut, errOut, certPaths...)
}

func generate(logOut, errOut io.Writer, certPaths ...*paths.CertPaths) error {
	zerologLogger, err := logger.NewLogger(
		logOut, errOut,
		&logger.Config{Prod: false, LogLevel: "warn", LogLevelParsed: zerolog.WarnLevel},
	)
	if err != nil {
		return errors.Wrap(err, "failed to create logger")
	}

	generator := certs.NewGenerator(zerologLogger)

	for _, certPaths := range certPaths {
		if err := generator.MakeDevCert(&certs.CertGenConfig{
			CommonName:       certPaths.Name,
			CertKeyPath:      certPaths.Key,
			CertPath:         certPaths.Cert,
			CACertPath:       certPaths.CA,
			DefaultTLSGenDir: certPaths.Dir,
		}); err != nil {
			return errors.Wrap(err, "failed to create dev certs")
		}
	}

	return nil
}
