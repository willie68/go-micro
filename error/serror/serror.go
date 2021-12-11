package serror

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Serr error model
type Serr struct {
	Code   int    `json:"code"`
	Key    string `json:"key"`
	Srv    string `json:"service,omitempty"`
	Msg    string `json:"message,omitempty"`
	Origin string `json:"origin,omitempty"`
}

// Service the service name
var Service string
var wrapper = make([]func(err error) *Serr, 0)

// Error returns the error
func (e *Serr) Error() string {
	if e.Key == "" {
		e.Key = "unexpected-error"
	}
	byt, err := json.Marshal(e)
	if err != nil {
		return e.str()
	}
	return string(byt)
}

// New creates an error
func New(code int, args ...string) *Serr {
	return build(&Serr{
		Code: code,
	}, nil, args...)
}

// Wrapper are used to wrap specific errors into service errors
// please do not call Wrap or Wrapc in the implemented function
func Wrapper(fn func(err error) *Serr) {
	wrapper = append(wrapper, fn)
}

// Wrap wraps an error
func Wrap(err error, args ...string) *Serr {
	return Wrapc(err, http.StatusInternalServerError, args...)
}

// Wrapc wraps an error (with code)
func Wrapc(err error, code int, args ...string) *Serr {
	if err == nil {
		return nil
	}
	if err, ok := err.(*Serr); ok {
		return err
	}
	for _, w := range wrapper {
		e := w(err)
		if e != nil {
			return build(e, err, args...)
		}
	}
	return build(&Serr{
		Code: code,
	}, err, args...)
}

// Unauthorized unauthorized error
func Unauthorized(err error, args ...string) *Serr {
	return build(&Serr{
		Key:  "unauthorized",
		Code: http.StatusUnauthorized,
	}, err, args...)
}

// Forbidden forbidden error
func Forbidden(err error, args ...string) *Serr {
	return build(&Serr{
		Key:  "forbidden",
		Code: http.StatusForbidden,
	}, err, args...)
}

// BadRequest bad request error
func BadRequest(err error, args ...string) *Serr {
	return build(&Serr{
		Key:  "badrequest",
		Code: http.StatusBadRequest,
	}, err, args...)
}

// NotFound not found error
func NotFound(typ string, id string, err ...error) *Serr {
	var first error
	if len(err) > 0 {
		first = err[0]
	}
	return build(&Serr{
		Msg:  fmt.Sprintf("could not find %s %s", typ, id),
		Key:  fmt.Sprintf("%s-not-found", typ),
		Code: http.StatusNotFound,
	}, first)
}

// Checks if the error is an service error of the given code
func Is(err error, code int) bool {
	if e, ok := err.(*Serr); ok {
		return e.Code == code
	}
	return false
}

func (e *Serr) str() string {
	s := make([]string, 0)
	if e.Msg != "" {
		s = append(s, e.Msg)
	}
	s = append(s, fmt.Sprintf(", code: %d", e.Code))
	s = append(s, fmt.Sprintf(", key: %s", e.Key))
	if e.Srv != "" {
		s = append(s, fmt.Sprintf(", service: %s", e.Srv))
	}
	if e.Origin != "" {
		s = append(s, fmt.Sprintf(", origin: %s", e.Origin))
	}
	return strings.Join(s, "")
}

func build(e *Serr, err error, args ...string) *Serr {
	if err != nil {
		e.Origin = err.Error()
	}
	if len(args) > 0 {
		e.Key = args[0]
	}
	if len(args) > 1 {
		e.Msg = args[1]
	}
	if Service != "" {
		e.Srv = Service
	}
	return e
}
