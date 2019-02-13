package token

import (
	"crypto/rsa"
	"encoding/json"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/security/claims"
	"gopkg.in/square/go-jose.v2"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
)

type JWTValidator struct {
	rsaPublicKey *rsa.PublicKey
}

func NewJWTValidator(rsaPublicKey *rsa.PublicKey) JWTValidator {
	return JWTValidator{rsaPublicKey: rsaPublicKey}
}

func (jwtv *JWTValidator) ValidateJWT(jwt string) (claims.Claims, error) {
	// Parse the jwt. Successful parse means the content of authorisation header was jwt
	jwtObject, err := jose.ParseSigned(jwt)
	if err != nil {
		log.Warn("Invalid JWT Submitted!")
		return nil, err
	}

	// Verify jwt signature and retrieve json marshalled claims
	// Failure indicates jwt was damaged or tampered with
	jsonClaims, err := jwtObject.Verify(jwtv.rsaPublicKey)
	if err != nil {
		log.Warn("JWT verification failure!")
		return nil, err
	}

	// Unmarshal json claims
	wrapped := wrappedClaims.WrappedClaims{}
	err = json.Unmarshal(jsonClaims, &wrapped)
	if err != nil {
		// This is an unknown flop, by now things shouldn't flop
		log.Warn("Unable to Unmarshal login claims!")
		return nil, err
	}

	// Unwrap the claims and return the result
	return wrapped.Unwrap()
}
