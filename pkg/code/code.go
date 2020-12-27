package code

import (
	"time"

	"github.com/guillaumebchd/polaris/pkg/crypto"
)

// This is used in the authorization code flow

const Lifetime = 60 * time.Second

type CodeData struct {
	Openid     bool   `json:"openid"`
	PKCEMethod string `json:"pkce_method"`
	PKCEValue  string `json:"pkce_value"`
}

func GenerateCode(data CodeData) (string, error) {
	code, err := generateUniqueCode()
	if err != nil {
		panic(err)
	}

	err = storeCodeInCache(code, data, Lifetime)
	return code, err
}

func Check(code string) (CodeData, error) {
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
