package serror

import "errors"

// Error definitions
var (
	ErrNotImplementedYet = errors.New("not implemented yet")
	ErrUnknowError       = errors.New("unknown error")
	ErrAlreadyExists     = errors.New("object already exists")
	ErrLoginFailed       = errors.New("login failed")
	ErrNotExists         = errors.New("object not exists")
	ErrTokenExpired      = errors.New("token expired")
	ErrTokenNotValid     = errors.New("token not valid")
	ErrMissingID         = errors.New("missing id")
)
