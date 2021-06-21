package dev

import (
	"github.com/aserto-dev/aserto/pkg/cc"

	"github.com/cli/browser"
	"github.com/fatih/color"
)

const (
	webConsoleURL = "https://localhost:8383"
)

type ConsoleCmd struct{}

func (cmd ConsoleCmd) Run(c *cc.CommonCtx) error {
	color.Green(">>> launch onebox console...")

	return browser.OpenURL(webConsoleURL)
}
