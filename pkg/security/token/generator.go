package token

import (
	"crypto/rsa"
	"encoding/json"
	"github.com/iot-my-world/brain/internal/log"
	claims2 "github.com/iot-my-world/brain/pkg/security/claims"
	wrapped2 "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	"gopkg.in/square/go-jose.v2"
)

type JWTGenerator struct {
	privateKey *rsa.PrivateKey
	signer     jose.Signer
}

func NewJWTGenerator(privateRSAKey *rsa.PrivateKey) JWTGenerator {
	//Create a new signer using RSASSA-PSS (SHA512) with the given private key.
	joseSigner, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.PS512, Key: privateRSAKey}, nil)
	if err != nil {
		log.Fatal(err)
	}

	return JWTGenerator{
		privateKey: privateRSAKey,
		signer:     joseSigner,
	}
}

func (g JWTGenerator) GenerateToken(claimsToWrap claims2.Claims) (string, error) {
	wrapped, err := wrapped2.Wrap(claimsToWrap)
	if err != nil {
		return "", err
	}
	return getSignedJWT(wrapped, g.signer)
}

func getSignedJWT(claims interface{}, signer jose.Signer) (string, error) {
	//Marshall the claims data to a json string
	claimsPayload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	//Sign the marshalled payload
	signedObj, err := signer.Sign(claimsPayload)
	if err != nil {
		return "", err
	}

	//Serialise the signed object
	signedJWT, err := signedObj.CompactSerialize()
	if err != nil {
		return "", err
	}

	return signedJWT, nil
}
