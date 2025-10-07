package vote

import (
	"encoding/json"
	"flip-planning-poker/internal/model"
	"flip-planning-poker/internal/repository"
	"flip-planning-poker/internal/utils"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VoteService struct {
	db   *pgxpool.Pool
	repo *repository.VoteRepository
}

func NewVoteService(database *pgxpool.Pool) *VoteService {
	return &VoteService{
		db:   database,
		repo: repository.NewVoteRepository(database),
	}
}

func (s *VoteService) CreateVote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var v model.Vote

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		utils.SendError(w, http.StatusBadRequest, err, "Erro ao ler dados da requisição")
		return
	}

	err := s.repo.CreateVote(ctx, &v)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err, "Erro ao criar voto")
		return
	}

	utils.SendSuccess(w, http.StatusCreated, v, "Voto criado com sucesso")
}

func (s *VoteService) FindVotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	storyId := r.URL.Query().Get("story_id")

	query := repository.VoteQuery{
		StoryId: storyId,
	}

	var votes []model.Vote

	if err := s.repo.FindVotes(ctx, &votes, &query); err != nil {
		utils.SendError(w, http.StatusInternalServerError, err, "Erro ao buscar votos")
		return
	}
	utils.SendSuccessWithTotal(w, http.StatusOK, votes, len(votes), "Busca por votos realizada com sucesso")

}
