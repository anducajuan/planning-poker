package vote

import (
	"encoding/json"
	"flip-planning-poker/internal/model"
	"flip-planning-poker/internal/repository"
	"flip-planning-poker/internal/utils"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool
var voteRepository *repository.VoteRepository

func Init(database *pgxpool.Pool) {
	db = database
	voteRepository = repository.NewVoteRepository(db)
}

func CreateVote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var v model.Vote

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		utils.SendError(w, http.StatusBadRequest, err, "Erro ao ler dados da requisição")
	}

	err := voteRepository.CreateVote(ctx, &v)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err, "Erro ao criar voto")
		return
	}

	utils.SendSuccess(w, http.StatusCreated, v, "Voto criado com sucesso")
}
