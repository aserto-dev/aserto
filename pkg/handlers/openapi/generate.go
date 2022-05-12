package openapi

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/getkin/kin-openapi/openapi3"
)

// Test with: https://github.com/splunk/splunk-cloud-sdk-go/tree/master/services/action

type GenerateOpenAPI struct {
	Path string `arg:"" required:"" help:"path to openapi.yaml"`
	Name string `arg:"" required:"" help:"path to openapi.yaml"`
}

func parseURI(uri string) []string {
	result := []string{}
	parts := strings.Split(uri, "/")
	for _, part := range parts[1:] {
		if strings.Contains(part, "{") {
			var clean = strings.Replace(strings.Replace(part, "{", "", -1), "}", "", -1)
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

func (cmd *GenerateOpenAPI) Run(c *cc.CommonCtx) error {

	specURL, err := url.Parse(cmd.Path)
	if err != nil {
		log.Fatal("Failed to parse spec URL", err)
	}
	root := cmd.Name

	doc, err := openapi3.NewLoader().LoadFromURI(specURL)

	if err != nil {
		log.Fatal("Failed to load spec from URL", err)
	}

	paths := doc.Paths

	var packages []string

	for uri, path := range paths {
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

	var policiesDirectoryName = root + "/src/policies"
	if _, err := os.Stat(root); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(policiesDirectoryName, os.ModePerm)

		if err != nil {
			log.Fatal(err)
		}
	}

	for _, pkg := range packages {
		var packageHeader = "package " + pkg
		var policy = packageHeader + "\n\n" + "default allowed = false"
		var filename = pkg + ".rego"
		destination, err := os.Create(policiesDirectoryName + "/" + filename)
		if err != nil {
			fmt.Println("os.Create:", err)
			return nil
		}
		destination.Close()

		fmt.Fprint(destination, policy)
	}

	return nil
}
