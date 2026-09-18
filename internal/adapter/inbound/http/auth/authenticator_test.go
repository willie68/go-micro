package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromContext(t *testing.T) {
	ast := assert.New(t)

	j, err := DecodeJWT(testToken)

	ctx := context.TODO()
	ctx = context.WithValue(ctx, TokenCtxKey, &j)

	jt, mp, err := FromContext(ctx)
	ast.Nil(err)
	ast.NotNil(jt)

	header := jt.Header
	ast.NotNil(header)
	ast.Equal("RS256", header["alg"])
	ast.Equal("JWT", header["typ"])

	payload := jt.Payload
	ast.NotNil(payload)
	ast.Equal("83e94672-94f8-4760-a63f-ce0f069a1351", payload["sub"])
	ast.Equal("Wilfriedd Klaas", payload["name"])
	ast.Equal(float64(1619853483), payload["iat"])

	sig := jt.Signature
	ast.NotNil(sig)
	ast.Equal(testTokenSignature, sig)

	ast.Equal("Wilfriedd Klaas", mp["name"])
}
