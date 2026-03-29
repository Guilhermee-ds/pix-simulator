package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"pix-simulator/internal/database"
	"pix-simulator/internal/queue"
	"pix-simulator/internal/service"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

func worker(ctx context.Context, rdb *redis.Client, svc *service.Service) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("worker encerrando...")
			return
		default:
			job, err := queue.Pop(rdb, "pix")
			if err != nil {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			err = svc.Process(job)
			if err != nil {
				fmt.Println("error:", err)
			}
		}
	}
}

func main() {
	db := database.Connet()
	svc := service.New(db)

	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	for i := 0; i < 200; i++ {
		go worker(ctx, rdb, svc)
	}

	<-ctx.Done()
	fmt.Println("encerrando workers...")
}
