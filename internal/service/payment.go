package service

import "database/sql"

type Service struct {
	DB *sql.DB
}

func New(db *sql.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) Process(id string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	var sender, receiver string
	var amount float64

	err = tx.QueryRow(`
		SELECT sender_account, receiver_account, amount
		FROM transactions WHERE end_to_end_id=$1 FOR UPDATE
	`, id).Scan(&sender, &receiver, &amount)

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id=$2", amount, receiver)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Exec("UPDATE transactions SET status='done' WHERE end_to_end_id=$1", id)
	return tx.Commit()
}
