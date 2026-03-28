package idempotency

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func Exists(rdb *redis.Client, key string) bool {
	val, _ := rdb.Exists(ctx, key).Result()
	return val == 1
}

func Save(rdb *redis.Client, key string) {
	rdb.Set(ctx, key, "1", time.Minute*10)
}
