package code

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var rdb *redis.Client
var ctx = context.TODO()

func init() {
	godotenv.Load()

	redis_url := os.Getenv("REDIS_URL")
	redis_port := os.Getenv("REDIS_PORT")
	redis_password := os.Getenv("REDIS_PASSWORD")

	rdb = redis.NewClient(&redis.Options{
		Addr:     redis_url + ":" + redis_port,
		Password: redis_password,
		DB:       0,
	})
}

func storeCodeInCache(code string, data CodeData, lifetime time.Duration) error {
	v, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, code, v, lifetime).Err()
}

func getCodeInCache(code string) (CodeData, error) {
	value, err := rdb.Get(ctx, code).Bytes()
	if err != nil {
		return CodeData{}, err
	}

	var data CodeData
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
