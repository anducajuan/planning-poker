package main

import (
	"flip-planning-poker/config"
	"flip-planning-poker/middleware"
	"flip-planning-poker/services/session"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Carregar configurações
	cfg := config.LoadConfig()

	initDB(cfg)
	defer db.Close()

	// Configurar dependências dos serviços
	session.SetDB(db)
	session.SetBroadcast(broadcast)

	router := mux.NewRouter()

	// Middlewares
	router.Use(middleware.CORS)
	router.Use(middleware.Logger)

	// REST
	router.HandleFunc("/sessions", session.GetSessions).Methods("GET")
	router.HandleFunc("/sessions", session.CreateSession).Methods("POST")
	router.HandleFunc("/sessions/{id}", session.DeleteSession).Methods("DELETE")

	// WebSocket
	router.HandleFunc("/ws", handleConnections)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Rota não encontrada: %s %s", r.Method, r.URL.Path)
		http.Error(w, "rota não encontrada", http.StatusNotFound)
	})

	log.Printf("Servidor rodando em :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
