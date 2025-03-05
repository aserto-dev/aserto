package conf_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/zenizh/go-capturer"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto/pkg/cmd/conf"
)

type testCLI struct {
	Cfg string `name:"config" short:"c" type:"conf" env:"TEST_ENV"`
}

func newParser(t *testing.T, cli *testCLI, options ...kong.Option) *kong.Kong {
	t.Helper()

	parser, err := kong.New(
		cli,
		options...,
	)
	require.NoError(t, err)

	return parser
}

type configDir struct {
	t     *testing.T
	Dir   string
	Files []string
}

func NewConfigDir(t *testing.T, dir string, files ...string) *configDir {
	t.Helper()

	confDir := &configDir{t, dir, []string{}}
	confDir.AddEmptyFiles(files...)

	return confDir
}

func (d *configDir) File() string {
	d.t.Helper()

	require.Len(d.t, d.Files, 1, "must have exactly one file")
	return d.FirstFile()
}

func (d *configDir) FirstFile() string {
	return d.Files[0]
}

func (d *configDir) AddEmptyFiles(filenames ...string) {
	d.t.Helper()
	for _, name := range filenames {
		path := filepath.Join(d.Dir, name)
		f, err := os.Create(path)
		require.NoError(d.t, err)
		f.Close()

		d.Files = append(d.Files, path)
	}

	sort.Strings(d.Files)
}

type params struct {
	t          *testing.T
	configPath string
	parser     *kong.Kong
	cli        *testCLI
}

func TestConfigFileMapper(t *testing.T) {
	tests := []struct {
		Name string
		Run  func(p *params)
	}{
		{
			"File name with extension",
			func(p *params) {
				dir := NewConfigDir(p.t, p.configPath, "test.yaml")
				_, err := p.parser.Parse([]string{"-c", "test.yaml"})
				require.NoError(p.t, err)
				require.Equal(p.t, dir.File(), p.cli.Cfg)
			},
		},
		{
			"File name without extension",
			func(p *params) {
				dir := NewConfigDir(p.t, p.configPath, "test.yaml")
				_, err := p.parser.Parse([]string{"-c", "test"})
				require.NoError(p.t, err)
				require.Equal(p.t, dir.File(), p.cli.Cfg)
			},
		},
		{
			"Multiple files with same name",
			func(p *params) {
				NewConfigDir(p.t, p.configPath, "test.yaml", "test.json")
				stderr := capturer.CaptureStderr(func() {
					_, err := p.parser.Parse([]string{"-c", "test"})
					require.NoError(p.t, err)
				})
				require.Regexp(p.t, `WARNING: The specified configuration \('test'\) matches multiple configuration files.*`, stderr)
			},
		},
		{
			"File outside config dir",
			func(p *params) {
				dir := NewConfigDir(p.t, p.t.TempDir(), "test.yaml")
				_, err := p.parser.Parse([]string{"-c", dir.File()})
				require.NoError(p.t, err)
				require.Equal(p.t, dir.File(), p.cli.Cfg)
			},
		},
		{
			"File in env var",
			func(p *params) {
				dir := NewConfigDir(p.t, p.configPath, "test.yaml")
				p.t.Setenv("TEST_ENV", "test")
				_, err := p.parser.Parse([]string{})
				require.NoError(p.t, err)
				require.Equal(p.t, dir.File(), p.cli.Cfg)
			},
		},
		{
			"File doesn't exist",
			func(p *params) {
				_, err := p.parser.Parse([]string{"-c", "test"})
				require.ErrorIs(p.t, errors.Cause(err), conf.ErrConfigNotFound)
			},
		},
		{
			"Path doesn't exist",
			func(p *params) {
				_, err := p.parser.Parse([]string{"-c", "path/test"})
				require.ErrorIs(p.t, errors.Cause(err), conf.ErrConfigNotFound)
			},
		},
	}

	for _, test := range tests {
		t.Run(
			test.Name,
			func(tt *testing.T) {
				cli := &testCLI{}
				tempDir := tt.TempDir()
				p := newParser(
					t,
					cli,
					kong.NamedMapper("conf", conf.ConfigFileMapper(tempDir)),
				)
				test.Run(&params{tt, tempDir, p, cli})
			},
		)
	}
}
