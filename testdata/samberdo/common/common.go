package common

type Common struct {
	name string
}

func New(name string) *Common {
	return &Common{
		name: name,
	}
}

func (c *Common) Name() string {
	return c.name
}
