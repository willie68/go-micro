package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/willie68/go-micro/internal/config"
)

type JWTAuthConfig struct {
	Validate  bool
	TenantKey string
}

type JWT struct {
	Token     string
	Header    map[string]interface{}
	Payload   map[string]interface{}
	Signature string
	IsValid   bool
}

type JWTAuth struct {
	Config JWTAuthConfig
}

func ParseJWTConfig(cfg config.Authentcation) (JWTAuthConfig, error) {
	jwtcfg := JWTAuthConfig{}
	var err error
	jwtcfg.Validate, err = config.GetConfigValueAsBool(cfg.Properties, "validate")
	if err != nil {
		return jwtcfg, err
	}
	//	jwtcfg.TenantKey, err = config.GetConfigValueAsString(cfg.Properties, "tenantKey")
	//	if err != nil {
	//		return jwtcfg, err
	//	}
	return jwtcfg, nil
}

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

func jwtDecodePart(payload string) (map[string]interface{}, error) {
	var result map[string]interface{}
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

func (j *JWT) Validate(cfg JWTAuthConfig) error {
	//TODO here should be the implementation of the validation of the token
	return nil
}
