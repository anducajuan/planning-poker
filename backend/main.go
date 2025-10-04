package main

import (
	"flip-planning-poker/config"
	"flip-planning-poker/middleware"
	"flip-planning-poker/services/sessions"
	"flip-planning-poker/services/users"
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
	sessions.SetDB(db)
	sessions.SetBroadcast(broadcast)

	users.SetDB(db)

	router := mux.NewRouter()

	// Middlewares
	router.Use(middleware.CORS)
	router.Use(middleware.Logger)

	// REST
	router.HandleFunc("/sessions", sessions.GetSessions).Methods("GET")
	router.HandleFunc("/sessions", sessions.CreateSession).Methods("POST")
	router.HandleFunc("/sessions/{id}", sessions.DeleteSession).Methods("DELETE")
	router.HandleFunc("/users", users.CreateUser).Methods("POST")
	router.HandleFunc("/users", users.GetUsers)
	// WebSocket
	router.HandleFunc("/ws", handleConnections)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Rota não encontrada: %s %s", r.Method, r.URL.Path)
		http.Error(w, "rota não encontrada", http.StatusNotFound)
	})

	log.Printf("Servidor rodando em :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
