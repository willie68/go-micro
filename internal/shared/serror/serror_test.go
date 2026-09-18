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

func TestUnauthorizedForbiddenConflict(t *testing.T) {
	Service = service
	u := Unauthorized(nil, "no-auth", "please login")
	if u.Code != 401 || u.Key != "no-auth" {
		t.Fatalf("unauthorized: %+v", u)
	}
	f := Forbidden(nil, "no-access", "denied")
	if f.Code != 403 {
		t.Fatalf("forbidden: %+v", f)
	}
	b := BadRequest(nil, "bad", "invalid")
	if b.Code != 400 {
		t.Fatalf("bad request: %+v", b)
	}
	n := NotFound("address", "1")
	if n.Code != 404 {
		t.Fatalf("not found: %+v", n)
	}
	i := InternalServerError(errors.New("boom"))
	if i.Code != 500 {
		t.Fatalf("internal: %+v", i)
	}
	c := Conflict(errors.New("dup"))
	if c.Code != 409 {
		t.Fatalf("conflict: %+v", c)
	}
}

func TestIsAndStr(t *testing.T) {
	e := New(404, "missing", "gone")
	if !Is(e, 404) {
		t.Fatal("expected matching code")
	}
	if Is(e, 500) || Is(errors.New("x"), 404) {
		t.Fatal("unexpected Is result")
	}
	e.Key = ""
	_ = e.Error()
	e.Key = "k"
	e.Msg = "m"
	e.Srv = "s"
	e.Origin = "o"
	s := e.str()
	if s == "" {
		t.Fatal("expected str output")
	}
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

func TestWrapcNilAndSerr(t *testing.T) {
	if Wrapc(nil, 500) != nil {
		t.Fatal("nil wrap should stay nil")
	}
	orig := New(400, "k", "m")
	if Wrapc(orig, 500) != orig {
		t.Fatal("existing Serr should be returned as-is")
	}
}
