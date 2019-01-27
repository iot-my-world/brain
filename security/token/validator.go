package token

import (
	"crypto/rsa"
	"gitlab.com/iotTracker/brain/security/claims"
	"gopkg.in/square/go-jose.v2"
	"gitlab.com/iotTracker/brain/log"
	"encoding/json"
	"time"
	"errors"
)

type JWTValidator struct {
	rsaPublicKey *rsa.PublicKey
}

func NewJWTValidator (rsaPublicKey *rsa.PublicKey) JWTValidator {
	return JWTValidator{rsaPublicKey: rsaPublicKey}
}


func (jwtv *JWTValidator) ValidateJWT(jwt string) (*claims.LoginClaims, error) {
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
	loginClaims := claims.LoginClaims{}
	err = json.Unmarshal(jsonClaims, &loginClaims)
	if err != nil {
		// This is an unknown flop, by now things shouldn't flop
		log.Warn("Unable to Unmarshal login claims!")
		return nil, err
	}

	// Check if token has expired
	if time.Now().UTC().After(time.Unix(loginClaims.ExpirationTime, 0).UTC()) {
		return nil, errors.New("Token Has Expired!")
	}

	return &loginClaims, nil
}