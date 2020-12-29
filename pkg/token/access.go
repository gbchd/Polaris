package token

import (
	"errors"
	"reflect"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
)

/*
	Struct to encode the available claims in an access_token.
*/
type AccessToken struct {
	Issuer   string   `mapstructure:"iss"`
	Subject  string   `mapstructure:"sub"`
	Audience string   `mapstructure:"aud"`
	Scope    []string `mapstructure:"scope"`
}

// Encode a given access_token in a string format using the private key.
func (at *AccessToken) Encode() (string, error) {
	Claims := at.createClaims()
	Claims["type"] = "access"
	Claims["exp"] = time.Now().Unix() + int64(AccessTokenLifetime)
	Claims["iat"] = time.Now().Unix()

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, Claims)

	token, err := t.SignedString(privKey.Key)
	return token, err
}

// Decode a given access_token in a AccessToken struct. Useful to check the scope for example. TODO
func (at *AccessToken) Decode(token string) error {
	t, ok := verify(token)
	if !ok {
		return errors.New("Invalid Token")
	}

	claims := t.Claims.(jwt.MapClaims)

	if claims["type"] != "access" {
		return errors.New("Not an access_token")
	}

	mapstructure.Decode(claims, at)

	return nil
}

func (at *AccessToken) createClaims() jwt.MapClaims {
	Claims := jwt.MapClaims{}

	rv := reflect.ValueOf(at).Elem()
	typeOfS := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		if !rv.Field(i).IsZero() {
			Claims[typeOfS.Field(i).Tag.Get("mapstructure")] = rv.Field(i).Interface()
		}
	}

	return Claims
}
