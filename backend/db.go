package main

import (
	"context"
	"flip-planning-poker/config"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func initDB(cfg *config.Config) {
	var err error
	db, err = pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Erro ao conectar no banco: ", err)
	}

	// Testar a conexão
	if err = db.Ping(context.Background()); err != nil {
		log.Fatal("Erro ao fazer ping no banco: ", err)
	}

	log.Println("Conexão com banco de dados estabelecida com sucesso")
}
