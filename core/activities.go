package core

type Activities struct {
	c *Core
}

func NewActivities(c *Core) *Activities {
	return &Activities{c: c}
}

func (a *Activities) Test() error {
	return nil
}
