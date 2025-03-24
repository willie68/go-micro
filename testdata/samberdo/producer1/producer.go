package producer1

import "fmt"

type Producer1 struct {
	name string
}

func New(name string) *Producer1 {
	return &Producer1{
		name: name,
	}
}

func (p *Producer1) Doit(in string) string {
	return fmt.Sprintf("Prod 1: %s_%s", p.name, in)
}
