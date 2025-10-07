package main

import (
	"flip-planning-poker/internal/config"
	"flip-planning-poker/internal/database"
	"flip-planning-poker/internal/middleware"
	"flip-planning-poker/internal/services/session"
	"flip-planning-poker/internal/services/story"
	"flip-planning-poker/internal/services/user"
	"flip-planning-poker/internal/services/vote"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	database.InitDB(cfg)

	sessionService := session.NewSessionService(database.GetDB())

	userService := user.NewUserService(database.GetDB())
	storyService := story.NewStoryService(database.GetDB())
	voteService := vote.NewVoteService(database.GetDB())

	router := mux.NewRouter()

	router.Use(middleware.CORS)
	router.Use(middleware.Logger)

	router.HandleFunc("/sessions", sessionService.GetSessions).Methods("GET", "OPTIONS")
	router.HandleFunc("/sessions", sessionService.CreateSession).Methods("POST", "OPTIONS")
	router.HandleFunc("/sessions/{id}", sessionService.DeleteSession).Methods("DELETE", "OPTIONS")

	router.HandleFunc("/users", userService.CreateUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/users", userService.GetUsers).Methods("GET", "OPTIONS")

	router.HandleFunc("/stories", storyService.GetSessionStories).Methods("GET", "OPTIONS")
	router.HandleFunc("/stories", storyService.CreateStory).Methods("POST", "OPTIONS")

	router.HandleFunc("/votes", voteService.CreateVote).Methods("POST", "OPTIONS")
	router.HandleFunc("/votes", voteService.FindVotes).Methods("GET", "OPTIONS")

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Rota não encontrada: %s %s", r.Method, r.URL.Path)
		http.Error(w, "rota não encontrada", http.StatusNotFound)
	})

	log.Printf("Servidor rodando em :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
