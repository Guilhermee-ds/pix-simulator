package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func Connet() *sql.DB {
	db, err := sql.Open("postgres", "postgres://pix:pix@postgres:5432/pix?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(50)           // máximo de conexões simultâneas
	db.SetMaxIdleConns(25)           // conexões ociosas mantidas abertas
	db.SetConnMaxLifetime(time.Hour) // tempo máximo de vida de uma conexão

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}
