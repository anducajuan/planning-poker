package handlers

import (
	"encoding/json"
	"flip-planning-poker/internal/models"
	"flip-planning-poker/internal/services"
	"flip-planning-poker/internal/utils/response"
	"net/http"

	"github.com/gorilla/mux"
)

type VoteHandler struct {
	service *services.VoteService
}

func NewVoteHandler(service *services.VoteService) *VoteHandler {
	return &VoteHandler{
		service: service,
	}
}

func (h *VoteHandler) GetPathPrefix() string {
	return "/votes"
}

func (h *VoteHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("", h.ListVotes).Methods("GET", "OPTIONS")
	r.HandleFunc("", h.CreateVote).Methods("POST", "OPTIONS")
}

func (h *VoteHandler) ListVotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	storyId := r.URL.Query().Get("story_id")

	votes, err := h.service.ListVotes(ctx, storyId)
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, err, "Erro ao buscar votos")
		return
	}

	response.SendJSONResponse(w, http.StatusOK, votes)
}

func (h *VoteHandler) CreateVote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	var vote models.Vote
	if err := json.NewDecoder(r.Body).Decode(&vote); err != nil {
		response.SendError(w, http.StatusBadRequest, err, "Dados inválidos no corpo da requisição")
		return
	}

	createdVote, err := h.service.CreateVote(ctx, &vote)
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, err, "Erro ao criar voto")
		return
	}

	response.SendJSONResponse(w, http.StatusCreated, createdVote)
}
