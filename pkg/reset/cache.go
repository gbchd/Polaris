package reset

import (
	"time"

	"github.com/go-redis/redis/v8"
)

func Connexion() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     RedisUrl + ":" + RedisPort,
		Password: RedisPassword,
		DB:       1,
	})
}

func storeCodeInCache(code string, email string, lifetime time.Duration) error {
	return rdb.Set(ctx, code, email, lifetime).Err()
}

func getCodeInCache(code string) (string, error) {
	value, err := rdb.Get(ctx, code).Bytes()
	return string(value), err
}

func isCodeInCache(code string) bool {
	_, err := rdb.Get(ctx, code).Result()
	return err == nil
}

func deleteCodeInCache(code string) error {
	_, err := rdb.Del(ctx, code).Result()
	return err
}
