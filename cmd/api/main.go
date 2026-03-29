package main

import (
	"encoding/json"
	"net/http"
	"time"

	"pix-simulator/internal/idempotency"
	"pix-simulator/internal/queue"

	"github.com/redis/go-redis/v9"
)

type Request struct {
	Sender   string  `json:"sender"`
	Receiver string  `json:"receiver"`
	Amount   float64 `json:"amount"`
}

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	http.HandleFunc("/pix", func(w http.ResponseWriter, r *http.Request) {

		key := r.Header.Get("Idempotency-Key")

		if idempotency.Exists(rdb, key) {
			w.Write([]byte("duplicated"))
			return
		}

		var req Request
		json.NewDecoder(r.Body).Decode(&req)

		txID := "tx_" + time.Now().Format("20060102150405")

		queue.Push(rdb, "pix", queue.Job{ID: txID})

		idempotency.Save(rdb, key)

		w.Write([]byte("queued"))
	})

	http.ListenAndServe(":8080", nil)
}
