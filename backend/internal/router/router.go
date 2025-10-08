package router

import (
	"flip-planning-poker/internal/handlers"
	middleware "flip-planning-poker/internal/middlewares"
	"flip-planning-poker/internal/services"
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

type Handler interface {
	RegisterRoutes(r *mux.Router)
	GetPathPrefix() string
}

func (r *ApiRouter) registerHandlers(router *mux.Router, handlers []Handler) {
	if len(handlers) > 0 {
		for _, handler := range handlers {
			pathPrefix := handler.GetPathPrefix()
			subRouter := router.PathPrefix(pathPrefix).Subrouter()
			handler.RegisterRoutes(subRouter)
		}
	}
}

func (r *ApiRouter) NewRouter(database *pgxpool.Pool) (*mux.Router, error) {
	r.router = mux.NewRouter()

	if database == nil {
		log.Fatal("banco de dados não definido!")
	}
	r.setDB(database)

	r.router.Use(middleware.CORS)
	r.router.Use(middleware.Logger)

	// Lista com todos os Handlers que serão registrados
	handlersList := []Handler{}

	// Inicializa Serviços
	sessionService := services.NewSessionService(r.db)
	handlersList = append(handlersList, handlers.NewSessionHandler(sessionService))

	userService := services.NewUserService(r.db)
	handlersList = append(handlersList, handlers.NewUserHandler(userService))

	storyService := services.NewStoryService(r.db)
	handlersList = append(handlersList, handlers.NewStoryHandler(storyService))

	voteService := services.NewVoteService(r.db)
	handlersList = append(handlersList, handlers.NewVoteHandler(voteService))

	// Registra Handlers
	r.registerHandlers(r.router, handlersList)

	r.router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Rota não encontrada: %s %s", r.Method, r.URL.Path)
		http.Error(w, "rota não encontrada", http.StatusNotFound)
	})

	return r.router, nil
}
