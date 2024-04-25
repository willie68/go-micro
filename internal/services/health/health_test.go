package health

import (
	"errors"
	"testing"
	"time"

	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
)

var (
	hs *SHealth

	_ Check = &MyCheck{}
)

type MyCheck struct {
	name  string
	fired bool
	times int
	ret   bool
	err   error
}

// Check implements Check.
func (m *MyCheck) Check() (bool, error) {
	m.fired = true
	m.times++
	return m.ret, m.err
}

// CheckName implements Check.
func (m *MyCheck) CheckName() string {
	return m.name
}

func InitHealth(ast *assert.Assertions) {
	if hs == nil {
		cfg := Config{
			Period:     10,
			StartDelay: 1,
		}
		h, err := NewHealthSystem(cfg)
		hs = h
		ast.Nil(err)
		ast.NotNil(hs)
	}
}

func ShutdownHealth() {
	_ = do.Shutdown[*SHealth](nil)
	hs = nil
}

func TestHealthBase(t *testing.T) {
	ast := assert.New(t)

	InitHealth(ast)

	hsdi, err := do.Invoke[*SHealth](nil)
	ast.Nil(err)
	ast.NotNil(hsdi)

	chk := MyCheck{
		name:  "myname",
		fired: false,
		times: 0,
		ret:   true,
		err:   nil,
	}
	hs.Register(&chk)

	time.Sleep(12 * time.Second)

	ast.True(chk.fired)
	ast.Equal(1, chk.times)

	hs.checkHealthCheckTimer()
	ast.True(hs.readyz)
	ast.Equal(0, len(hs.messages))

	ok := hs.Unregister(chk.name)
	ast.True(ok)

	ShutdownHealth()
}

func TestMessage(t *testing.T) {
	ast := assert.New(t)

	n := time.Now()
	InitHealth(ast)
	hs.lastChecked = n
	hs.messages = make([]string, 0)
	msg := hs.Message()
	ast.Equal(n.String(), msg.LastCheck)
	ast.Equal(0, len(msg.Messages))

	ShutdownHealth()
}

func TestHealthUnhealthy(t *testing.T) {
	ast := assert.New(t)

	InitHealth(ast)

	hsdi, err := do.Invoke[*SHealth](nil)
	ast.Nil(err)
	ast.NotNil(hsdi)

	chk := MyCheck{
		name:  "myname",
		fired: false,
		times: 0,
		ret:   false,
		err:   errors.New("error"),
	}
	err = Register(&chk)
	ast.Nil(err)

	time.Sleep(12 * time.Second)

	ast.True(chk.fired)
	ast.Equal(1, chk.times)

	msg := hs.Message()
	ast.NotNil(msg.LastCheck)
	ast.Equal(1, len(msg.Messages))
	ast.Equal("myname: error", msg.Messages[0])

	err = Unregister(chk.name)
	ast.Nil(err)

	ShutdownHealth()
}
