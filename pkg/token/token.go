/*
Package Token implements a simple library for creating oauth and openid tokens.

These tokens are JWT token created with the package "github.com/dgrijalva/jwt-go". We only implements the RSA256 signing method at the moment, other methods may come later.

This package is initialize with environment variable, you should call godotenv.load() in your main function if you want to load a .env file, then call Initialize() to load these environment variables.

These environment variables are :
	- TOKEN_PRIVATE_KEY_PATH : the private rsa key's path used to sign the tokens
	- TOKEN_PUBLIC_KEY_PATH : the public rsa key's path used to verify the tokens' signature.
	- TOKEN_ACCESS_LIFETIME : the lifetime of an access token (in seconds)
	- TOKEN_REFRESH_LIFETIME : the lifetime of a refresh token (in seconds)
	- TOKEN_ID_LIFETIME : the lifetime of an id token (in seconds)

You can also directly initialize the variables by code. If you set the variables PublicKeyPath and PrivateKeyPath by hand you should call the method LoadKeys() afterwards.



*/
package token

import (
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

var (
	// The public rsa key's path used to verify the tokens' signature.
	PublicKeyPath string

	// The private rsa key's path used to sign the tokens
	PrivateKeyPath string

	// The lifetime of an access token
	AccessTokenLifetime time.Duration

	// The lifetime of a refresh token
	RefreshTokenLifetime time.Duration

	// The lifetime of an id token
	IdTokenLifetime time.Duration

	// The two variables that will be used to sign and verify our JWTs.
	pubKey  PublicKey
	privKey PrivateKey
)

// Function init that read the environment variables and generate the rsa keys at the end.
func Initialize() error {
	PublicKeyPath = os.Getenv("TOKEN_PUBLIC_KEY_PATH")
	PrivateKeyPath = os.Getenv("TOKEN_PRIVATE_KEY_PATH")

	accessLifetime, err := strconv.Atoi(os.Getenv("TOKEN_ACCESS_LIFETIME"))
	if err == nil {
		AccessTokenLifetime = time.Duration(accessLifetime) * time.Second
	}

	refreshLifetime, err := strconv.Atoi(os.Getenv("TOKEN_REFRESH_LIFETIME"))
	if err == nil {
		RefreshTokenLifetime = time.Duration(refreshLifetime) * time.Second
	}

	idLifetime, err := strconv.Atoi(os.Getenv("TOKEN_ID_LIFETIME"))
	if err == nil {
		IdTokenLifetime = time.Duration(idLifetime) * time.Second
	}

	err = LoadKeys()
	return err
}

func Verify(t string) bool {
	_, ok := verify(t)
	return ok
}

func verify(t string) (*jwt.Token, bool) {
	token, err := extractToken(t)
	if err != nil {
		return nil, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	return token, ok && token.Valid && claims.VerifyExpiresAt(time.Now().Unix(), true)
}

func extractToken(t string) (*jwt.Token, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return pubKey.Key, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
