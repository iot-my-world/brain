package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/go-errors/errors"
	"github.com/iot-my-world/brain/internal/log"
	"io/ioutil"
	"os"
)

func FetchPrivateKey(dir string) *rsa.PrivateKey {
	privateKeyBitCount := 4096
	privateKeyFilePath := dir + "privateKey.pem"
	publicKeyFilePath := dir + "publicKey.pem"

	pvtKey, err := fetchRSAPrivateKeyFromFile(privateKeyFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Info("Private Key not found at: '" + privateKeyFilePath + "' Generating a new Key Pair.")
			pvtKey, err = generateRSAKeyPair(privateKeyBitCount)
			if err = saveRSAPrivateKeyToFile(privateKeyFilePath, pvtKey); err != nil {
				log.Fatal("Failed to save private key to file!", err)
			}
		} else {
			log.Fatal(err)
		}
	}

	//Update Public Key File
	if err = saveRSAPublicKeyToFile(publicKeyFilePath, &pvtKey.PublicKey); err != nil {
		log.Fatal("Failed to save public key to file!", err)
	}

	return pvtKey
}

func generateRSAKeyPair(privateKeyBitCount int) (*rsa.PrivateKey, error) {
	pvtKey, err := rsa.GenerateKey(rand.Reader, privateKeyBitCount)
	if err != nil {
		return nil, err
	}
	return pvtKey, nil
}

func saveRSAPrivateKeyToFile(privateKeyFilePath string, privateKey *rsa.PrivateKey) error {
	outputFile, err := os.Create(privateKeyFilePath)
	defer outputFile.Close()
	if err != nil {
		return err
	}

	var privateKeyBlock = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	return pem.Encode(outputFile, privateKeyBlock)
}

func saveRSAPublicKeyToFile(publicKeyFilePath string, publicKey *rsa.PublicKey) error {
	outputFile, err := os.Create(publicKeyFilePath)
	defer outputFile.Close()
	if err != nil {
		return err
	}

	publicKeyByteData, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}

	var publicKeyBlock = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyByteData,
	}

	return pem.Encode(outputFile, publicKeyBlock)
}

func fetchRSAPrivateKeyFromFile(privateKeyFilePath string) (*rsa.PrivateKey, error) {
	//Check that file exists and is readable
	if _, err := os.Stat(privateKeyFilePath); err != nil {
		return nil, err
	}

	//Read the file
	fileBytes, err := ioutil.ReadFile(privateKeyFilePath)
	if err != nil {
		return nil, err
	}

	//Return the parsed key
	return parseRSAPrivateKeyFromString(string(fileBytes))
}

func fetchRSAPublicKeyFromFile(publicKeyFilePath string) (*rsa.PublicKey, error) {
	//Check that the file exists and is readable
	if _, err := os.Stat(publicKeyFilePath); err != nil {
		return nil, err
	}

	//Read the file
	filebytes, err := ioutil.ReadFile(publicKeyFilePath)
	if err != nil {
		return nil, err
	}

	//Return the Parsed key
	return parseRSAPublicKeyFromString(string(filebytes))
}

func parseRSAPrivateKeyFromString(rsaPrivateKeyString string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(rsaPrivateKeyString))
	if block == nil {
		return nil, errors.New("Failed to parse private key string!")
	}

	pvtKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pvtKey, nil
}

func parseRSAPublicKeyFromString(rsaPublicKeyString string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(rsaPublicKeyString))
	if block == nil {
		return nil, errors.New("Failed to parse public key string!")
	}

	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pubKey, nil
}
