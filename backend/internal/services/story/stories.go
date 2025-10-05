package story

import (
	"encoding/json"
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

func CreateStory(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var s model.Story

	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		utils.SendError(w, http.StatusBadRequest, err, "Dados inválidos no corpo da requisição")
		return
	}
	storyRepository := repository.NewStoryRepository(db)
	err := storyRepository.CreateStory(&s)

	if s.Status == "ACTUAL" {
		utils.Logger("Alterando tipos de outras Stories para 'OLD'")
		storyRepository.SetStoriesToOld(s.SessionID, s.ID)
	}

	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err, "Erro ao criar story")
		return
	}
	utils.SendSuccess(w, http.StatusOK, s, "Story criada com sucesso")
}
