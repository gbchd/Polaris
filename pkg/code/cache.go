package code

import (
	"context"
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

func storeCodeInCache(code string, lifetime time.Duration) error {
	err := rdb.Set(ctx, code, "", lifetime).Err()
	return err
}

func isCodeInCache(code string) bool {
	_, err := rdb.Get(ctx, code).Result()
	//fmt.Println(err)
	return err == nil
}

func deleteCodeInCache(code string) error {
	_, err := rdb.Del(ctx, code).Result()
	return err
}
