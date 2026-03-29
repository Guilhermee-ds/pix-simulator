package main

import (
	"encoding/json"
	"net/http"

	"pix-simulator/internal/idempotency"
	"pix-simulator/internal/queue"

	"github.com/google/uuid"
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
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		txID := "tx_" + uuid.New().String()

		if err := queue.Push(rdb, "pix", queue.Job{
			ID:       txID,
			Sender:   req.Sender,
			Receiver: req.Receiver,
			Amount:   req.Amount,
		}); err != nil {
			http.Error(w, "queue error", http.StatusInternalServerError)
			return
		}

		idempotency.Save(rdb, key)
		w.Write([]byte("queued"))
	})

	http.ListenAndServe(":8080", nil)
}
