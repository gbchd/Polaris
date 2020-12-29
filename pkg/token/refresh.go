package token

import (
	"crypto/sha256"
	"errors"
	"reflect"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
)

type RefreshToken struct {
	Issuer   string `mapstructure:"iss"`
	Subject  string `mapstructure:"sub"`
	Audience string `mapstructure:"aud"`
}

/*
Encode a given refresh_token in a string format using the private key.

Need an access_token already encoded as input.
*/
func (rt *RefreshToken) Encode(access_token string) (string, error) {
	h := sha256.New()
	h.Write([]byte(access_token))

	Claims := rt.createClaims()
	Claims["type"] = "refresh"
	Claims["at_hash"] = h.Sum(nil)
	Claims["exp"] = time.Now().Unix() + int64(RefreshTokenLifetime)
	Claims["iat"] = time.Now().Unix()

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, Claims)

	token, err := t.SignedString(privKey.Key)
	return token, err
}

// Decode a given refresh_token in a RefreshToken struct.
func (rt *RefreshToken) Decode(token string) error {
	t, ok := verify(token)
	if !ok {
		return errors.New("Invalid Token")
	}

	claims := t.Claims.(jwt.MapClaims)

	if claims["type"] != "refresh" {
		return errors.New("Not a refresh_token")
	}

	mapstructure.Decode(claims, rt)

	return nil
}

func (rt *RefreshToken) createClaims() jwt.MapClaims {
	Claims := jwt.MapClaims{}

	v := reflect.ValueOf(rt).Elem()
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsZero() {
			Claims[typeOfS.Field(i).Tag.Get("mapstructure")] = v.Field(i).Interface()
		}
	}

	return Claims
}
