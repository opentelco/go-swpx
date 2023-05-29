package activities

import "git.liero.se/opentelco/go-swpx/core"

type Activities struct {
	c *core.Core
}

func New(c *core.Core) *Activities {
	return &Activities{c: c}
}

func (a *Activities) Test() error {
	return nil
}
