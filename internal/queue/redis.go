package queue

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Job struct {
	ID string
}

func Push(rdb *redis.Client, key string, job Job) {
	data, _ := json.Marshal(job)
	rdb.LPush(ctx, key, data)
}

func Pop(rdb *redis.Client, key string) (Job, error) {
	val, err := rdb.RPop(ctx, key).Result()
	if err != nil {
		return Job{}, err
	}

	var job Job
	json.Unmarshal([]byte(val), &job)
	return job, nil
}
