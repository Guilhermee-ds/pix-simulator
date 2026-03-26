package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func Connet() *sql.DB {
	db, err := sql.Open("postgres", "postgres://pix:pix@postgres:5432/pix?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}
