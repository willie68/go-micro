package producer

import "fmt"

type Producer struct {
	name string
}

func New(name string) *Producer {
	return &Producer{
		name: name,
	}
}

func (p *Producer) Doit(in string) string {
	return fmt.Sprintf("Prod: %s_%s", p.name, in)
}
