package token

import (
	"crypto/rsa"
	"encoding/json"
	tokenException "gitlab.com/iotTracker/brain/security/token/exception"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"gopkg.in/square/go-jose.v2"
)

type JWTValidator struct {
	rsaPublicKey *rsa.PublicKey
}

func NewJWTValidator(rsaPublicKey *rsa.PublicKey) JWTValidator {
	return JWTValidator{rsaPublicKey: rsaPublicKey}
}

func (jwtv *JWTValidator) ValidateJWT(jwt string) (wrappedClaims.WrappedClaims, error) {
	// Parse the jwt. Successful parse means the content of authorisation header was jwt
	jwtObject, err := jose.ParseSigned(jwt)
	if err != nil {
		return wrappedClaims.WrappedClaims{}, tokenException.InvalidJWT{Reasons: []string{err.Error()}}
	}

	// Verify jwt signature and retrieve json marshalled claims
	// Failure indicates jwt was damaged or tampered with
	jsonClaims, err := jwtObject.Verify(jwtv.rsaPublicKey)
	if err != nil {
		return wrappedClaims.WrappedClaims{}, tokenException.JWTVerification{Reasons: []string{err.Error()}}
	}

	// Unmarshal json claims
	wrapped := wrappedClaims.WrappedClaims{}
	err = json.Unmarshal(jsonClaims, &wrapped)
	if err != nil {
		// This is an unknown flop, by now things shouldn't flop
		return wrappedClaims.WrappedClaims{}, tokenException.JWTUnmarshalling{Reasons: []string{err.Error()}}
	}

	// Unwrap the claims and return the result
	return wrapped, nil
}
