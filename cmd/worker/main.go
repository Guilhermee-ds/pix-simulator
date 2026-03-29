package worker

import (
	"fmt"
	"pix-simulator/internal/database"
	"pix-simulator/internal/queue"
	"pix-simulator/internal/service"
	"time"

	"github.com/redis/go-redis/v9"
)

func worker(rdb *redis.Client, svc *service.Service) {
	for {
		job, err := queue.Pop(rdb, "pix")
		if err != nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}

		err = svc.Process(job.ID)
		if err != nil {
			fmt.Println("error:", err)
		}
	}
}

func main() {
	db := database.Connet()
	svc := service.New(db)

	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	for i := 0; i < 200; i++ {
		go worker(rdb, svc)
	}
}
