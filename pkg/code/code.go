package code

import (
	"time"

	"github.com/guillaumebchd/polaris/pkg/crypto"
)

// This is used in the authorization code flow

const Lifetime = 60 * time.Second

func GenerateCode() (string, error) {
	code, err := generateUniqueCode()
	if err != nil {
		panic(err)
	}

	err = storeCodeInCache(code, Lifetime)

	return code, err
}

func Check(code string) bool {
	ok := isCodeInCache(code)

	if ok {
		deleteCodeInCache(code)
	}

	return ok
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
