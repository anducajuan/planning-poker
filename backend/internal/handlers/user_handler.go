package handlers

import (
	"encoding/json"
	"flip-planning-poker/internal/models"
	"flip-planning-poker/internal/services"
	"flip-planning-poker/internal/utils/response"
	"net/http"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) GetPathPrefix() string {
	return "/users"
}

func (h *UserHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("", h.ListUsers).Methods("GET", "OPTIONS")
	r.HandleFunc("", h.CreateUser).Methods("POST", "OPTIONS")
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionId := r.URL.Query().Get("session_id")

	users, err := h.service.ListUsers(ctx, sessionId)
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, err, "Erro ao buscar usuários")
		return
	}

	response.SendJSONResponse(w, http.StatusOK, users)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		response.SendError(w, http.StatusBadRequest, err, "Dados inválidos no corpo da requisição")
		return
	}

	createdUser, err := h.service.CreateUser(ctx, &user)
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, err, "Erro ao criar usuário")
		return
	}

	response.SendJSONResponse(w, http.StatusCreated, createdUser)
}
