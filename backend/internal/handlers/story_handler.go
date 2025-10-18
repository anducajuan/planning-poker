package handlers

import (
	"encoding/json"
	"flip-planning-poker/internal/models"
	"flip-planning-poker/internal/services"
	"flip-planning-poker/internal/utils/response"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type StoryHandler struct {
	service *services.StoryService
}

func NewStoryHandler(service *services.StoryService) *StoryHandler {
	return &StoryHandler{
		service: service,
	}
}

func (h *StoryHandler) GetPathPrefix() string {
	return "/stories"
}

func (h *StoryHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("", h.ListStories).Methods("GET", "OPTIONS")
	r.HandleFunc("", h.CreateStory).Methods("POST", "OPTIONS")
	r.HandleFunc("/{id}/reveal", h.RevealStory).Methods("POST", "OPTIONS")
}

func (h *StoryHandler) RevealStory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		response.SendError(w, http.StatusBadRequest, nil, "o id da story é obrigatório")
	}

	var storyId int
	var err error

	storyId, err = strconv.Atoi(idStr)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, err, "Necessário informar um id de sessão")
		return
	}

	err = h.service.RevealStory(ctx, storyId)
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, err, "Erro ao revelar story")
		return
	}
	response.SendSuccess(w, http.StatusCreated, nil, "Story revelada com sucesso")

}

func (h *StoryHandler) ListStories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionId := r.URL.Query().Get("session_id")

	if sessionId == "" {
		response.SendError(w, http.StatusBadRequest, nil, "Necessário informar um id de sessão válido")
		return
	}

	stories, err := h.service.ListStories(ctx, sessionId)
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, err, "Erro ao buscar stories")
		return
	}

	response.SendJSONResponse(w, http.StatusOK, stories)
}

func (h *StoryHandler) CreateStory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	var story models.Story
	if err := json.NewDecoder(r.Body).Decode(&story); err != nil {
		response.SendError(w, http.StatusBadRequest, err, "Dados inválidos no corpo da requisição")
		return
	}

	createdStory, err := h.service.CreateStory(ctx, &story)
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, err, "Erro ao criar story")
		return
	}

	response.SendJSONResponse(w, http.StatusCreated, createdStory)
}
