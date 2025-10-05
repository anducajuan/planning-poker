package story

import (
	"encoding/json"
	"flip-planning-poker/internal/model"
	"flip-planning-poker/internal/repository"
	"flip-planning-poker/internal/utils"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StoryService struct {
	db   *pgxpool.Pool
	repo *repository.StoryRepository
}

func NewStoryService(database *pgxpool.Pool) *StoryService {
	return &StoryService{
		db:   database,
		repo: repository.NewStoryRepository(database),
	}
}

func (s *StoryService) GetSessionStories(w http.ResponseWriter, r *http.Request) {
	sessionId := r.URL.Query().Get("session_id")
	if sessionId == "" {
		utils.SendError(w, http.StatusBadRequest, nil, "Necessário informar um id de sessão válido")
		return
	}

	stories, err := s.repo.FindStoryBySessionId(sessionId)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err, "Erro ao buscar stories")
		return
	}
	utils.SendSuccessWithTotal(w, http.StatusOK, stories, len(stories), "Busca realizada com sucesso")
}

func (s *StoryService) CreateStory(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var story model.Story

	if err := json.NewDecoder(r.Body).Decode(&story); err != nil {
		utils.SendError(w, http.StatusBadRequest, err, "Dados inválidos no corpo da requisição")
		return
	}

	err := s.repo.CreateStory(&story)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err, "Erro ao criar story")
		return
	}

	if story.Status == "ACTUAL" {
		utils.Logger("Alterando tipos de outras Stories para 'OLD'")
		s.repo.SetStoriesToOld(story.SessionID, story.ID)
	}

	utils.SendSuccess(w, http.StatusCreated, story, "Story criada com sucesso")
}
