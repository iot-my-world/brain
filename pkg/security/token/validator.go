package token

import (
	"crypto/rsa"
	"encoding/json"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	"github.com/iot-my-world/brain/pkg/security/token/exception"
	"gopkg.in/square/go-jose.v2"
)

type JWTValidator struct {
	rsaPublicKey *rsa.PublicKey
}

func NewJWTValidator(rsaPublicKey *rsa.PublicKey) JWTValidator {
	return JWTValidator{rsaPublicKey: rsaPublicKey}
}

func (jwtv *JWTValidator) ValidateJWT(jwt string) (wrappedClaims.Wrapped, error) {
	// Parse the jwt. Successful parse means the content of authorisation header was jwt
	jwtObject, err := jose.ParseSigned(jwt)
	if err != nil {
		return wrappedClaims.Wrapped{}, exception.InvalidJWT{Reasons: []string{err.Error()}}
	}

	// Verify jwt signature and retrieve json marshalled claims
	// Failure indicates jwt was damaged or tampered with
	jsonClaims, err := jwtObject.Verify(jwtv.rsaPublicKey)
	if err != nil {
		return wrappedClaims.Wrapped{}, exception.JWTVerification{Reasons: []string{err.Error()}}
	}

	// Unmarshal json claims
	wrapped := wrappedClaims.Wrapped{}
	err = json.Unmarshal(jsonClaims, &wrapped)
	if err != nil {
		// This is an unknown flop, by now things shouldn't flop
		return wrappedClaims.Wrapped{}, exception.JWTUnmarshalling{Reasons: []string{err.Error()}}
	}

	// Unwrap the claims and return the result
	return wrapped, nil
}
