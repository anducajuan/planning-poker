package handlers

import (
	"encoding/json"
	"flip-planning-poker/internal/models"
	"flip-planning-poker/internal/services"
	"flip-planning-poker/internal/utils/response"
	"net/http"

	"github.com/gorilla/mux"
)

type SessionHandler struct {
	service *services.SessionService
}

func NewSessionHandler(service *services.SessionService) *SessionHandler {
	return &SessionHandler{
		service: service,
	}
}

func (h *SessionHandler) GetPathPrefix() string {
	return "/sessions"
}

func (h *SessionHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("", h.ListSessions).Methods("GET", "OPTIONS")
	r.HandleFunc("", h.CreateSession).Methods("POST", "OPTIONS")
	r.HandleFunc("/{id}", h.DeleteSession).Methods("DELETE", "OPTIONS")
}

func (h *SessionHandler) ListSessions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	ctx := r.Context()
	defer r.Body.Close()

	sessions, err := h.service.ListSessions(ctx, query)
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, err, "Erro ao buscar sessões")
		return
	}

	response.SendJSONResponse(w, http.StatusOK, sessions)
}

func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	var session models.Session
	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		response.SendError(w, http.StatusBadRequest, err, "Dados inválidos no corpo da requisição")
		return
	}

	createdSession, err := h.service.CreateSession(ctx, &session)
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, err, "Erro ao criar sessão")
		return
	}

	response.SendJSONResponse(w, http.StatusCreated, createdSession)
}

func (h *SessionHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	if id == "" {
		response.SendError(w, http.StatusBadRequest, nil, "ID da sessão é obrigatório")
		return
	}

	err := h.service.DeleteSession(ctx, id)
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, err, "Erro ao deletar sessão")
		return
	}

	response.SendJSONResponse(w, http.StatusOK, map[string]string{"message": "Sessão deletada com sucesso"})
}
