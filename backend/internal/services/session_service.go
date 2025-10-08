package services

import (
	"context"
	"flip-planning-poker/internal/models"
	"flip-planning-poker/internal/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionService struct {
	db   *pgxpool.Pool
	repo *repositories.SessionRepository
}

func NewSessionService(database *pgxpool.Pool) *SessionService {
	return &SessionService{
		db:   database,
		repo: repositories.NewSessionRepository(database),
	}
}

func (s *SessionService) ListSessions(ctx context.Context, query any) ([]models.Session, error) {
	return s.repo.GetAll(ctx, query)
}

func (s *SessionService) CreateSession(ctx context.Context, session *models.Session) (*models.Session, error) {
	if err := s.repo.ValidateSessionData(ctx, session); err != nil {
		return nil, err
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *SessionService) DeleteSession(ctx context.Context, id string) error {
	if id == "" {
		return repositories.ErrInvalidID
	}

	rowsAffected, err := s.repo.DeleteSession(ctx, id)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return repositories.ErrNotFound
	}

	return nil
}
