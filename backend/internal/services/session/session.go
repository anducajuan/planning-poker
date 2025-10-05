package session

import (
	"encoding/json"
	"flip-planning-poker/internal/model"
	"flip-planning-poker/internal/repository"
	"flip-planning-poker/internal/utils"
	"flip-planning-poker/internal/websocket"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionService struct {
	db        *pgxpool.Pool
	repo      *repository.SessionRepository
	wsService websocket.Broadcaster
}

func NewSessionService(database *pgxpool.Pool) *SessionService {
	return &SessionService{
		db:   database,
		repo: repository.NewSessionRepository(database),
	}
}

func (s *SessionService) SetWebSocketService(wsService websocket.Broadcaster) {
	s.wsService = wsService
}

func (s *SessionService) GetSessions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessions, err := s.repo.GetSessions(ctx)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err, "Erro ao buscar sessões")
		return
	}

	utils.SendSuccessWithTotal(w, http.StatusOK, sessions, len(sessions), "Sessões recuperadas com sucesso")
}
func (s *SessionService) CreateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	var session model.Session

	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		utils.SendError(w, http.StatusBadRequest, err, "Dados inválidos no corpo da requisição")
		return
	}

	if err := s.repo.ValidateSessionData(ctx, &session); err != nil {
		utils.SendError(w, http.StatusBadRequest, err, err.Error())
		return
	}

	err := s.repo.CreateSession(ctx, &session)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err, "Erro ao criar sessão")
		return
	}

	if s.wsService != nil {
		sessionData := websocket.SessionData{
			ID:   session.ID,
			Name: session.Name,
		}
		s.wsService.NotifySessionCreated(sessionData)
	}

	utils.SendSuccess(w, http.StatusCreated, session, "Sessão criada com sucesso")
}

func (s *SessionService) DeleteSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	if id == "" {
		utils.SendError(w, http.StatusBadRequest, nil, "ID da sessão é obrigatório")
		return
	}

	rowsAffected, err := s.repo.DeleteSession(ctx, id)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err, "Erro ao deletar sessão")
		return
	}

	if rowsAffected == 0 {
		utils.SendError(w, http.StatusNotFound, nil, "Sessão não encontrada")
		return
	}

	if s.wsService != nil {
		s.wsService.NotifySessionDeleted(id)
	}

	utils.SendSuccess(w, http.StatusOK, nil, "Sessão deletada com sucesso")
}
