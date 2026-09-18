package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/willie68/go-micro/internal/config"
)

func TestTokenFromHeaderQueryCookie(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?jwt=fromquery", nil)
	assert.Empty(t, TokenFromHeader(req))
	req.Header.Set("Authorization", "Bearer abc.def")
	assert.Equal(t, "abc.def", TokenFromHeader(req))
	assert.Equal(t, "fromquery", TokenFromQuery(req))

	req.AddCookie(&http.Cookie{Name: "jwt", Value: "fromcookie"})
	assert.Equal(t, "fromcookie", TokenFromCookie(req))
	assert.Empty(t, TokenFromCookie(httptest.NewRequest(http.MethodGet, "/", nil)))
}

func TestAuthenticator(t *testing.T) {
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	h := Authenticator(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	jt, err := DecodeJWT(testToken)
	assert.NoError(t, err)
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(NewContext(req.Context(), &jt, nil, false))
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(NewContext(req.Context(), &jt, nil, true))
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestVerifierAndValidate(t *testing.T) {
	ja := &JWTAuth{Config: JWTAuthConfig{IgnorePages: []string{"/health"}}}
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	h := Verifier(ja)(next)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	assert.True(t, called)

	jt, err := DecodeJWT(testToken)
	assert.NoError(t, err)
	assert.NoError(t, jt.Validate(ja))

	_, err = VerifyToken(ja, testToken)
	assert.NoError(t, err)
	_, err = VerifyRequest(ja, httptest.NewRequest(http.MethodGet, "/", nil))
	assert.ErrorIs(t, err, ErrNoTokenFound)
}

func TestParseJWTConfigAndInit(t *testing.T) {
	cfg, err := ParseJWTConfig(config.Authentication{
		Type:       "jwt",
		Properties: map[string]any{"validate": true},
	})
	assert.NoError(t, err)
	assert.True(t, cfg.Validate)
	auth := InitJWT(cfg)
	assert.True(t, auth.Config.Validate)
	assert.Equal(t, "jwtauth context value Token", TokenCtxKey.String())
}

func TestDecodeJWTErrors(t *testing.T) {
	_, err := DecodeJWT("")
	assert.Error(t, err)
	_, err = DecodeJWT("onlyonepart")
	assert.Error(t, err)
}
