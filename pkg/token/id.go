package token

import (
	"errors"
	"reflect"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
)

type IdToken struct {
	Issuer            string `mapstructure:"iss"`
	Subject           string `mapstructure:"sub"`
	Audience          string `mapstructure:"aud"`
	Name              string `mapstructure:"name"`
	PreferredUsername string `mapstructure:"preferred_username"`
	Profile           string `mapstructure:"profile"`
	Picture           string `mapstructure:"picture"`
	Email             string `mapstructure:"email"`
	Gender            string `mapstructure:"gender"`
	PhoneNumber       string `mapstructure:"phone_number"`
	Address           string `mapstructure:"address"`
}

// Encode a given id_token in a string format using the private key.
func (it *IdToken) Encode() (string, error) {
	Claims := it.createClaims()
	Claims["type"] = "id"
	Claims["exp"] = time.Now().Unix() + int64(AccessTokenLifetime)
	Claims["iat"] = time.Now().Unix()

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, Claims)

	token, err := t.SignedString(privKey.Key)
	return token, err
}

// Decode a given access_token in a IdToken struct. TODO
func (it *IdToken) Decode(token string) error {
	t, ok := verify(token)
	if !ok {
		return errors.New("Invalid Token")
	}

	claims := t.Claims.(jwt.MapClaims)

	if claims["type"] != "id" {
		return errors.New("Not an id_token")
	}

	mapstructure.Decode(claims, it)

	return nil
}

func (it *IdToken) createClaims() jwt.MapClaims {
	Claims := jwt.MapClaims{}

	v := reflect.ValueOf(it).Elem()
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsZero() {
			Claims[typeOfS.Field(i).Tag.Get("mapstructure")] = v.Field(i).Interface()
		}
	}

	return Claims
}
