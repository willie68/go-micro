package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/willie68/go-micro/internal/config"
)

// JWTAuthConfig authentication/Authorisation configuration for JWT authentification
type JWTAuthConfig struct {
	Active      bool
	Validate    bool
	TenantClaim string
	Strict      bool
}

// JWT struct for the decoded jwt token
type JWT struct {
	Token     string
	Header    map[string]any
	Payload   map[string]any
	Signature string
	IsValid   bool
}

// JWTAuth the jwt authentication struct
type JWTAuth struct {
	Config JWTAuthConfig
}

// JWTConfig for the service
var JWTConfig = JWTAuthConfig{
	Active: false,
}

// InitJWT initialise the JWT for this service
func InitJWT(cnfg JWTAuthConfig) JWTAuth {
	JWTConfig = cnfg
	return JWTAuth{
		Config: cnfg,
	}
}

// ParseJWTConfig building up the dynamical configuration for this
func ParseJWTConfig(cfg config.Authentication) (JWTAuthConfig, error) {
	jwtcfg := JWTAuthConfig{
		Active: true,
	}
	var err error
	jwtcfg.Validate, err = config.GetConfigValueAsBool(cfg.Properties, "validate")
	if err != nil {
		return jwtcfg, err
	}
	return jwtcfg, nil
}

// DecodeJWT simple decode the jwt token string
func DecodeJWT(token string) (JWT, error) {
	jwt := JWT{
		Token:   token,
		IsValid: false,
	}

	if token == "" {
		return JWT{}, errors.New("missing token string")
	}

	if len(token) > 7 && strings.ToUpper(token[0:6]) == "BEARER" {
		token = token[7:]
	}

	// decode JWT token without verifying the signature
	jwtParts := strings.Split(token, ".")
	if len(jwtParts) < 2 {
		err := errors.New("token missing payload part")
		return jwt, err
	}
	var err error

	jwt.Header, err = jwtDecodePart(jwtParts[0])
	if err != nil {
		err = fmt.Errorf("token header parse error, %v", err)
		return jwt, err
	}

	jwt.Payload, err = jwtDecodePart(jwtParts[1])
	if err != nil {
		err = fmt.Errorf("token payload parse error, %v", err)
		return jwt, err
	}
	if len(jwtParts) > 2 {
		jwt.Signature = jwtParts[2]
	}
	jwt.IsValid = true
	return jwt, nil
}

func jwtDecodePart(payload string) (map[string]any, error) {
	var result map[string]any
	payloadData, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(payload)
	if err != nil {
		err = fmt.Errorf("token payload can't be decoded: %v", err)
		return nil, err
	}
	err = json.Unmarshal(payloadData, &result)
	if err != nil {
		err = fmt.Errorf("token payload parse error, %v", err)
		return nil, err
	}
	return result, nil
}

// Validate validation of the token is not implemented
func (j *JWT) Validate(_ JWTAuthConfig) error {
	//TODO here should be the implementation of the validation of the token
	return nil
}
