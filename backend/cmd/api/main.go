package main

import (
	"flip-planning-poker/internal/config"
	"flip-planning-poker/internal/database"
	"flip-planning-poker/internal/router"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	database.InitDB(cfg)
	defer database.Close()

	apiRouter := &router.ApiRouter{}

	r, err := apiRouter.NewRouter(database.GetDB())
	if err != nil {
		log.Fatalf("erro ao iniciar router: %v", err)
	}

	log.Printf("Servidor rodando em :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
