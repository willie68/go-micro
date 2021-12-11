package serror

import (
	"errors"
	"fmt"
	"testing"
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
	Service = "my-service"
	e := New(500, "my-error-key", "this is a message")
	t.Log(e.Error())
}

func TestWrap(t *testing.T) {
	Service = "my-service"
	e := Wrap(errors.New("my-error"), "my-error-key", "this is a message")
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
	Service = "my-service"
	e := Wrap(myerr, "my-error-key", "this is a message")
	t.Log(e.Error())
}
