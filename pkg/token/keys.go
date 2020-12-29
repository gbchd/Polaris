package token

import (
	"crypto/rsa"
	"io/ioutil"

	"github.com/dgrijalva/jwt-go"
)

// A struct that holds the information about a public RSA key.
type PublicKey struct {
	Raw []byte
	Key *rsa.PublicKey
}

// A struct that holds the information about a private RSA key.
type PrivateKey struct {
	Raw []byte
	Key *rsa.PrivateKey
}

// This method is used to read the public and private keys' path and generate these keys in PublicKey and PrivateKey struct format.
// This method should be called everytime the PublicKeyPath and PrivateKeyPath variables are modified.
func LoadKeys() error {
	privKeyRaw, err := ioutil.ReadFile(PrivateKeyPath)
	if err != nil {
		return err
	}

	privKeyValue, err := jwt.ParseRSAPrivateKeyFromPEM(privKeyRaw)
	if err != nil {
		return err
	}

	pubKeyRaw, err := ioutil.ReadFile(PublicKeyPath)
	if err != nil {
		return err
	}

	pubKeyValue, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyRaw)
	if err != nil {
		return err
	}

	privKey.Raw = privKeyRaw
	privKey.Key = privKeyValue

	pubKey.Raw = pubKeyRaw
	pubKey.Key = pubKeyValue

	return nil
}
