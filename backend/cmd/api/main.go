package main

import (
	"flip-planning-poker/internal/config"
	"flip-planning-poker/internal/database"
	"flip-planning-poker/internal/middleware"
	"flip-planning-poker/internal/services/session"
	"flip-planning-poker/internal/services/story"
	"flip-planning-poker/internal/services/user"
	"flip-planning-poker/internal/websocket"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	database.InitDB(cfg)

	session.SetDB(database.GetDB())
	session.SetBroadcast(websocket.GetBroadcastChannel())
	user.SetDB(database.GetDB())
	story.SetDB(database.GetDB())

	router := mux.NewRouter()

	router.Use(middleware.CORS)
	router.Use(middleware.Logger)

	router.HandleFunc("/sessions", session.GetSessions).Methods("GET", "OPTIONS")
	router.HandleFunc("/sessions", session.CreateSession).Methods("POST", "OPTIONS")
	router.HandleFunc("/sessions/{id}", session.DeleteSession).Methods("DELETE", "OPTIONS")

	router.HandleFunc("/users", user.CreateUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/users", user.GetUsers).Methods("GET", "OPTIONS")

	router.HandleFunc("/stories", story.GetSessionStories).Methods("GET", "OPTIONS")
	router.HandleFunc("/stories", story.CreateSession).Methods("POST", "OPTIONS")

	router.HandleFunc("/ws", websocket.HandleConnections)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Rota não encontrada: %s %s", r.Method, r.URL.Path)
		http.Error(w, "rota não encontrada", http.StatusNotFound)
	})

	log.Printf("Servidor rodando em :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
