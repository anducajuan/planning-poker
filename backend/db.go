package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func initDB() {
	var err error
	db, err = pgxpool.New(context.Background(), "postgres://admin:admin@localhost:5432/database?sslmode=disable")
	if err != nil {
		log.Fatal("Erro ao conectar no banco: ", err)
	}
}
