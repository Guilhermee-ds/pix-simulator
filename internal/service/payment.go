package service

import (
	"database/sql"
	"pix-simulator/internal/queue"
)

type Service struct {
	DB *sql.DB
}

func New(db *sql.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) Process(job queue.Job) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO transactions (id, end_to_end_id, sender_account, receiver_account, amount, status)
		VALUES ($1, $2, $3, $4, $5, 'pending')
	`, job.ID, job.ID, job.Sender, job.Receiver, job.Amount)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		UPDATE accounts SET balance = balance - $1 WHERE id = $2
	`, job.Amount, job.Sender)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		UPDATE accounts SET balance = balance + $1 WHERE id = $2
	`, job.Amount, job.Receiver)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Exec("UPDATE transactions SET status='done' WHERE end_to_end_id=$1", job.ID)

	return tx.Commit()
}
