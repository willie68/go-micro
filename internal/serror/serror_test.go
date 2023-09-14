package serror

import (
	"errors"
	"fmt"
	"testing"
)

const (
	service = "my-service"
	errKey  = "my-error-key"
	msg     = "this is a message"
)

type mySpecialError struct {
	SomeCode    int
	SomeMessage string
}

// Error returns the error
func (e *mySpecialError) Error() string {
	return fmt.Sprintf("%d %s", e.SomeCode, e.SomeMessage)
}

func TestNew(t *testing.T) {
	Service = service
	e := New(500, errKey, msg)
	t.Log(e.Error())
}

func TestWrap(t *testing.T) {
	Service = service
	e := Wrap(errors.New("my-error"), errKey, msg)
	t.Log(e.Error())
}

func TestWrapper(t *testing.T) {
	Wrapper(func(err error) *Serr {
		if merr, ok := err.(*mySpecialError); ok {
			return New(merr.SomeCode, "my-special-error", merr.SomeMessage)
		}
		return nil
	})
	myerr := &mySpecialError{
		SomeCode:    404,
		SomeMessage: "not found",
	}
	Service = service
	e := Wrap(myerr, errKey)
	t.Log(e.Error())
}
