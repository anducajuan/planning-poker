package story

import (
	"flip-planning-poker/internal/model"
	"flip-planning-poker/internal/repository"
	"flip-planning-poker/internal/utils"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func SetDB(database *pgxpool.Pool) {
	db = database
}

func GetSessionStories(w http.ResponseWriter, r *http.Request) {
	sessionId := r.URL.Query().Get("session_id")
	var err error
	if sessionId == "" {
		utils.SendError(w, http.StatusBadRequest, err, "Necessário informar um id de sessão válido")
		return
	}

	storyRepo := repository.NewStoryRepository(db)

	stories, err := storyRepo.FindStoryBySessionId(sessionId)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err, "Erro ao buscar stories")
		return
	}
	utils.SendSuccessWithTotal(w, http.StatusOK, stories, len(stories), "Busca realizada com sucesso")
}

// func CreateUser(w http.ResponseWriter, r *http.Request) {

// }

func validateStoryData(w http.ResponseWriter, user *model.Story) {

}

// func getStories
