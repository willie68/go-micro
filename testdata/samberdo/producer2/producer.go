package producer2

import "fmt"

type Producer2 struct {
	name string
}

func New(name string) *Producer2 {
	return &Producer2{
		name: name,
	}
}

func (p *Producer2) Doit(in string) string {
	return fmt.Sprintf("Prod 2: %s_%s", p.name, in)
}
