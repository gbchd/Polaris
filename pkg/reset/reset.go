package reset

import (
	"context"
	"fmt"
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
	RedisUrl = os.Getenv("CACHE_REDIS_URL")
	RedisPort = os.Getenv("CACHE_REDIS_PORT")
	RedisPassword = os.Getenv("CACHE_REDIS_PASSWORD")

	resetLifetime, err := strconv.Atoi(os.Getenv("RESET_LINK_KLIFETIME"))
	if err == nil {
		Lifetime = time.Duration(resetLifetime) * time.Second
	}

	Connexion()
	return err
}

func Generate(email string) (string, error) {
	code, err := generateUniqueCode()
	if err != nil {
		fmt.Println(err)
	}

	err = storeCodeInCache(code, email, Lifetime)
	return code, err
}

func Get(code string) (string, error) {
	email, err := getCodeInCache(code)
	if err == nil {
		deleteCodeInCache(code)
	}
	return email, err
}

func Exist(code string) bool {
	return isCodeInCache(code)
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
