package dev

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/pkg/errors"
)

type GenerateFromOpenAPI struct {
	URL  string `arg:"" required:"" help:"URL of the openapi.yaml"`
	Name string `arg:"" optional:"" help:"The name for the policy"`
}

const packageTemplate = `package %s

default allowed = false`

func parseURI(uri string) []string {
	result := []string{}
	parts := strings.Split(uri, "/")
	for _, part := range parts[1:] {
		if strings.Contains(part, "{") {
			clean := strings.Replace(strings.Replace(part, "{", "", -1), "}", "", -1)
			result = append(result, "__"+clean)
		} else {
			result = append(result, part)
		}
	}
	return result
}

func generatePackageName(root, verb, uri string) string {
	parts := []string{root, verb}
	parts = append(parts, parseURI(uri)...)
	return strings.Join(parts, ".")
}

func (cmd *GenerateFromOpenAPI) Run(c *cc.CommonCtx) error {

	specURL, err := url.Parse(cmd.URL)
	if err != nil {
		return errors.Wrapf(err, "Failed to parse spec URL [%s]", cmd.URL)
	}

	doc, err := openapi3.NewLoader().LoadFromURI(specURL)

	if err != nil {
		return errors.Wrapf(err, "Failed to load spec from URL [%s]", cmd.URL)
	}

	root := cmd.Name
	if cmd.Name == "" {
		root = strings.Replace(strings.ToLower(doc.Info.Title), " ", "_", -1)
	}

	packages := []string{}

	for uri, path := range doc.Paths {
		if path.Get != nil {
			packages = append(packages, generatePackageName(root, "GET", uri))
		}
		if path.Post != nil {
			packages = append(packages, generatePackageName(root, "POST", uri))
		}
		if path.Put != nil {
			packages = append(packages, generatePackageName(root, "PUT", uri))
		}
		if path.Delete != nil {
			packages = append(packages, generatePackageName(root, "DELETE", uri))
		}
		if path.Options != nil {
			packages = append(packages, generatePackageName(root, "OPTIONS", uri))
		}
	}

	policiesDirectoryName := root + "/src/policies"
	if _, err := os.Stat(root); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(policiesDirectoryName, os.ModePerm)

		if err != nil {
			return errors.Wrapf(err, "Failed to create the directory [%s]", policiesDirectoryName)
		}
	}

	for _, pkg := range packages {
		policy := fmt.Sprintf(packageTemplate, pkg)
		filename := pkg + ".rego"
		path := policiesDirectoryName + "/" + filename
		destination, err := os.Create(path)
		if err != nil {
			return errors.Wrapf(err, "Error creating the policy module file [%s]", path)
		}

		_, writeErr := fmt.Fprint(destination, policy)
		if writeErr != nil {
			return errors.Wrapf(writeErr, "Error writing to the policy module file [%s]", path)
		}

		defer destination.Close()
	}

	return nil
}
