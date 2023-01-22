// Package auth all things related to authentication and authorization
package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

// used context key for parameter given in a std context
var (
	TokenCtxKey = &contextKey{"Token"}
	ErrorCtxKey = &contextKey{"Error"}
)

// defining some common errors
var (
	ErrUnauthorized = errors.New("token is unauthorized")
	ErrExpired      = errors.New("token is expired")
	ErrNBFInvalid   = errors.New("token nbf validation failed")
	ErrIATInvalid   = errors.New("token iat validation failed")
	ErrNoTokenFound = errors.New("no token found")
	ErrAlgoInvalid  = errors.New("algorithm mismatch")
)

// FromContext extract the JWT and a flatten claim structure from a context
func FromContext(ctx context.Context) (*JWT, map[string]any, error) {
	token, ok := ctx.Value(TokenCtxKey).(*JWT)

	var err error
	var claims map[string]any

	if !ok || (token == nil) {
		claims = map[string]any{}
		return token, claims, errors.New("token not present")
	}

	claims = token.Payload
	err, _ = ctx.Value(ErrorCtxKey).(error)
	return token, claims, err
}

// Authenticator returns a handler for authentication
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := FromContext(r.Context())

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if token == nil || !token.IsValid {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

// Verifier returns a handler for verification
func Verifier(ja *JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return Verify(ja, TokenFromHeader, TokenFromCookie)(next)
	}
}

// Verify checking a request
func Verify(ja *JWTAuth, findTokenFns ...func(r *http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token, err := VerifyRequest(ja, r, findTokenFns...)
			ctx = NewContext(ctx, token, err)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}

// VerifyRequest this request -> authorization
func VerifyRequest(ja *JWTAuth, r *http.Request, findTokenFns ...func(r *http.Request) string) (*JWT, error) {
	var tokenString string

	// Extract token string from the request by calling token find functions in
	// the order they where provided. Further extraction stops if a function
	// returns a non-empty string.
	for _, fn := range findTokenFns {
		tokenString = fn(r)
		if tokenString != "" {
			break
		}
	}
	if tokenString == "" {
		return nil, ErrNoTokenFound
	}

	return VerifyToken(ja, tokenString)
}

// VerifyToken verify the bearer token
func VerifyToken(ja *JWTAuth, tokenString string) (*JWT, error) {
	// Decode & verify the token
	token, err := DecodeJWT(tokenString)
	if err != nil {
		return &token, err
	}

	if err := token.Validate(ja.Config); err != nil {
		return &token, err
	}

	// Valid!
	return &token, nil
}

// NewContext creating a new context for the http functions
func NewContext(ctx context.Context, t *JWT, err error) context.Context {
	ctx = context.WithValue(ctx, TokenCtxKey, t)
	ctx = context.WithValue(ctx, ErrorCtxKey, err)
	return ctx
}

// TokenFromHeader tries to retrieve the token string from the
// "Authorization" request header: "Authorization: BEARER T".
func TokenFromHeader(r *http.Request) string {
	// Get token from authorization header.
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

// TokenFromQuery tries to retrieve the token string from the "jwt" URI
// query parameter.
//
// To use it, build our own middleware handler, such as:
//
//	func Verifier(ja *JWTAuth) func(http.Handler) http.Handler {
//		return func(next http.Handler) http.Handler {
//			return Verify(ja, TokenFromQuery, TokenFromHeader, TokenFromCookie)(next)
//		}
//	}
func TokenFromQuery(r *http.Request) string {
	// Get token from query param named "jwt".
	return r.URL.Query().Get("jwt")
}

// TokenFromCookie tries to retrieve the token string from a cookie named
// "jwt".
func TokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return ""
	}
	return cookie.Value
}

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an any without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "jwtauth context value " + k.name
}
