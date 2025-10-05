package main

import (
	"flip-planning-poker/internal/config"
	"flip-planning-poker/internal/database"
	middleware "flip-planning-poker/internal/middlewares"
	"flip-planning-poker/internal/services/sessions"
	"flip-planning-poker/internal/services/users"
	"flip-planning-poker/internal/websocket"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	database.InitDB(cfg)

	sessions.SetDB(database.GetDB())
	sessions.SetBroadcast(websocket.GetBroadcastChannel())
	users.SetDB(database.GetDB())

	router := mux.NewRouter()

	router.Use(middleware.CORS)
	router.Use(middleware.Logger)

	router.HandleFunc("/sessions", sessions.GetSessions).Methods("GET", "OPTIONS")
	router.HandleFunc("/sessions", sessions.CreateSession).Methods("POST", "OPTIONS")
	router.HandleFunc("/sessions/{id}", sessions.DeleteSession).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/users", users.CreateUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/users", users.GetUsers).Methods("GET", "OPTIONS")

	router.HandleFunc("/ws", websocket.HandleConnections)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Rota não encontrada: %s %s", r.Method, r.URL.Path)
		http.Error(w, "rota não encontrada", http.StatusNotFound)
	})

	log.Printf("Servidor rodando em :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
