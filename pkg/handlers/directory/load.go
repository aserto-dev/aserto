package directory

import "github.com/aserto-dev/aserto/pkg/cc"

type LoadCmd struct{}

func (cmd *LoadCmd) Run(c cc.CommonCtx) error {
	return nil
}
