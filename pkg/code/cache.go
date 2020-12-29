package code

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

func Connexion() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     RedisUrl + ":" + RedisPort,
		Password: RedisPassword,
		DB:       0,
	})
}

func storeCodeInCache(code string, data Data, lifetime time.Duration) error {
	v, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, code, v, lifetime).Err()
}

func getCodeInCache(code string) (Data, error) {
	value, err := rdb.Get(ctx, code).Bytes()
	if err != nil {
		return Data{}, err
	}

	var data Data
	err = json.Unmarshal(value, &data)

	return data, err
}

func isCodeInCache(code string) bool {
	_, err := rdb.Get(ctx, code).Result()
	return err == nil
}

func deleteCodeInCache(code string) error {
	_, err := rdb.Del(ctx, code).Result()
	return err
}
