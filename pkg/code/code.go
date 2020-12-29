/*
Package Code implements a simple library for creating random code.

These codes are meant to be used in the authorization code flow of OAuth2.0. We use a custom library name crypto to generate the code.

The code are stored in a redis cache.
*/
package code

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/guillaumebchd/polaris/pkg/crypto"
)

var (
	Lifetime      time.Duration
	RedisUrl      string
	RedisPort     string
	RedisPassword string

	rdb *redis.Client
	ctx = context.TODO()
)

func Initialize() error {
	RedisUrl = os.Getenv("CODE_REDIS_URL")
	RedisPort = os.Getenv("CODE_REDIS_PORT")
	RedisPassword = os.Getenv("CODE_REDIS_PASSWORD")

	codeLifetime, err := strconv.Atoi(os.Getenv("TOKEN_ACCESS_LIFETIME"))
	if err == nil {
		Lifetime = time.Duration(codeLifetime) * time.Second
	}

	Connexion()
	return err
}

type Data struct {
	Email           string `json:"email"`
	ChallengeMethod string `json:"code_challenge_method"`
	Challenge       string `json:"code_challenge"`
}

func Generate(data Data) (string, error) {
	code, err := generateUniqueCode()
	if err != nil {
		panic(err)
	}

	err = storeCodeInCache(code, data, Lifetime)
	return code, err
}

func Get(code string) (Data, error) {
	data, err := getCodeInCache(code)
	if err == nil {
		deleteCodeInCache(code)
	}
	return data, err
}

func generateUniqueCode() (string, error) {
	code, err := crypto.RandomHex(16)
	if err != nil {
		return "", err
	}

	for isCodeInCache(code) {
		code, err = crypto.RandomHex(16)
		if err != nil {
			return "", err
		}
	}
	return code, err
}
