package basic

import (
	"crypto/rand"
	"encoding/base64"
	apiUserPasswordGenerator "github.com/iot-my-world/brain/user/api/password/generator"
	"github.com/satori/go.uuid"
)

type generator struct {
}

func New() apiUserPasswordGenerator.Generator {
	return &generator{}
}

func (g *generator) Generate(request *apiUserPasswordGenerator.GenerateRequest) (*apiUserPasswordGenerator.GenerateResponse, error) {
	keyBytes := make([]byte, 0)
	c := request.CryptoBytesLength
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	u, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	for _, value := range u.Bytes() {
		keyBytes = append(keyBytes, value)
	}
	for _, value := range b {
		keyBytes = append(keyBytes, value)
	}

	return &apiUserPasswordGenerator.GenerateResponse{
		Password: base64.StdEncoding.EncodeToString(keyBytes),
	}, nil
}
