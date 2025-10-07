package router

import (
	"errors"
	"flip-planning-poker/internal/middleware"
	"flip-planning-poker/internal/services/session"
	"flip-planning-poker/internal/services/story"
	"flip-planning-poker/internal/services/user"
	"flip-planning-poker/internal/services/vote"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ApiRouter struct {
	router *mux.Router
	db     *pgxpool.Pool
}

func (r *ApiRouter) setDB(database *pgxpool.Pool) {
	r.db = database
}

func (r *ApiRouter) NewRouter(database *pgxpool.Pool) (*mux.Router, error) {

	if database == nil {
		return nil, errors.New("banco de dados não definido")
	}
	r.setDB(database)

	sessionService := session.NewSessionService(r.db)
	userService := user.NewUserService(r.db)
	storyService := story.NewStoryService(r.db)
	voteService := vote.NewVoteService(r.db)

	r.router = mux.NewRouter()

	r.router.Use(middleware.CORS)
	r.router.Use(middleware.Logger)

	r.router.HandleFunc("/sessions", sessionService.GetSessions).Methods("GET", "OPTIONS")
	r.router.HandleFunc("/sessions", sessionService.CreateSession).Methods("POST", "OPTIONS")
	r.router.HandleFunc("/sessions/{id}", sessionService.DeleteSession).Methods("DELETE", "OPTIONS")

	r.router.HandleFunc("/users", userService.CreateUser).Methods("POST", "OPTIONS")
	r.router.HandleFunc("/users", userService.GetUsers).Methods("GET", "OPTIONS")

	r.router.HandleFunc("/stories", storyService.GetSessionStories).Methods("GET", "OPTIONS")
	r.router.HandleFunc("/stories", storyService.CreateStory).Methods("POST", "OPTIONS")

	r.router.HandleFunc("/votes", voteService.CreateVote).Methods("POST", "OPTIONS")
	r.router.HandleFunc("/votes", voteService.FindVotes).Methods("GET", "OPTIONS")

	r.router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Rota não encontrada: %s %s", r.Method, r.URL.Path)
		http.Error(w, "rota não encontrada", http.StatusNotFound)
	})

	return r.router, nil
}
