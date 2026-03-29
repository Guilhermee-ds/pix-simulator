package queue

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Job struct {
	ID       string
	Sender   string
	Receiver string
	Amount   float64
}

func Push(rdb *redis.Client, key string, job Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return rdb.LPush(ctx, key, data).Err()
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
