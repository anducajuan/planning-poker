package handlers

import (
	"encoding/json"
	"flip-planning-poker/internal/models"
	"flip-planning-poker/internal/repositories"
	"flip-planning-poker/internal/services"
	"flip-planning-poker/internal/utils/response"
	"net/http"
	"strconv"

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
	r.HandleFunc("/{id}", h.PatchVote).Methods("PATCH", "OPTIONS")
}

func (h *VoteHandler) ListVotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	storyId := r.URL.Query().Get("story_id")

	votes, err := h.service.List(ctx, storyId)
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

	createdVote, err := h.service.Create(ctx, &vote)
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, err, "Erro ao criar voto")
		return
	}

	response.SendJSONResponse(w, http.StatusCreated, createdVote)
}

func (h *VoteHandler) PatchVote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		response.SendError(w, http.StatusBadRequest, nil, "ID do voto é obrigatório")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, err, "ID do voto deve ser um número válido")
		return
	}

	var patch repositories.VotePatch
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		response.SendError(w, http.StatusBadRequest, err, "Dados inválidos no corpo da requisição")
		return
	}

	if err := h.service.Patch(ctx, id, &patch); err != nil {
		response.SendError(w, http.StatusInternalServerError, err, "Erro ao atualizar voto")
		return
	}

	updatedVote, err := h.service.Get(ctx, id)
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, err, "Erro ao buscar voto atualizado")
		return
	}

	response.SendJSONResponse(w, http.StatusOK, updatedVote)
}
