package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto/pkg/filex"
	"github.com/pkg/errors"
)

var ErrConfigNotFound = errors.New("cannot find configuration file")

// ConfigFileMapper is a kong.Mapper that resolves config files.
//
// When applied to a CLI flag, it attempts to find a configuration file that best matches the specified name using
// the following rules:
//  1. If the value is a full or relative path to an existing file, that file is chosen.
//  2. If the value is a file name (without a path separator) with an extension (e.g. "config.yaml") and a file with that
//     name exists in the config directory, that file is chosen.
//  3. If the value is a string without an dot (e.g. "eng") and a file with that name (i.e. "eng.*") exists in the
//     config directory, that file is chosen. If multiple files match, the first one is chosen and a warning is printed
//     to stderr.
type ConfigFileMapper string

func (m ConfigFileMapper) Decode(ctx *kong.DecodeContext, target reflect.Value) error {
	if target.Kind() != reflect.String {
		return errors.Errorf(`"conf" type must be applied to a string not %s`, target.Type())
	}

	var path string
	if err := ctx.Scan.PopValueInto("file", &path); err != nil {
		return err
	}

	if path != "-" {
		if resolved, err := m.find(path); err != nil {
			return err
		} else {
			path = resolved
		}
	}

	target.SetString(path)
	return nil
}

func (m ConfigFileMapper) find(path string) (string, error) {
	expanded := kong.ExpandPath(path)
	if filex.FileExists(expanded) {
		return expanded, nil
	}

	if strings.ContainsRune(path, filepath.Separator) {
		return "", errors.Wrap(ErrConfigNotFound, path)
	}

	// It's a filename with no path. Look in config directory.
	expanded = filepath.Join(string(m), path)
	if filepath.Ext(path) != "" {
		// It's a filename with an extension. Try to find the file.
		if filex.FileExists(expanded) {
			return expanded, nil
		}
	} else if matches, err := filepath.Glob(expanded + ".*"); err == nil && len(matches) > 0 {
		// Its just a name without an extension. Look for a match.
		if len(matches) > 1 {
			fmt.Fprintf(
				os.Stderr,
				"WARNING: The specified configuration ('%s') matches multiple configuration files: %s. Using '%s'",
				path,
				matches,
				matches[0],
			)
		}
		return matches[0], nil
	}

	return "", errors.Wrap(ErrConfigNotFound, path)
}
