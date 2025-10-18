package router

import (
	"flip-planning-poker/internal/handlers"
	middleware "flip-planning-poker/internal/middlewares"
	"flip-planning-poker/internal/services"
	"flip-planning-poker/internal/websocket"
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
	wsRouter := r.router.NewRoute().Subrouter()
	wsService := websocket.NewWebsocketService()

	wsRouter.HandleFunc("/ws", wsService.WsHandler)
	go wsService.HandleMessages()

	r.router.Use(middleware.CORS)
	r.router.Use(middleware.Logger)

	// Inicializa Serviços
	sessionService := services.NewSessionService(r.db)
	userService := services.NewUserService(r.db)
	storyService := services.NewStoryService(r.db)
	voteService := services.NewVoteService(r.db)

	// Lista com todos os Handlers que serão registrados
	handlersList := []Handler{
		handlers.NewSessionHandler(sessionService),
		handlers.NewUserHandler(userService),
		handlers.NewUserHandler(userService),
		handlers.NewStoryHandler(storyService),
		handlers.NewVoteHandler(voteService),
	}

	// Registra Handlers
	r.registerHandlers(r.router, handlersList)

	r.router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Rota não encontrada: %s %s", r.Method, r.URL.Path)
		http.Error(w, "rota não encontrada", http.StatusNotFound)
	})

	return r.router, nil
}
