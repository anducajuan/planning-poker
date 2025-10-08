package services

import (
	"context"
	"flip-planning-poker/internal/models"
	"flip-planning-poker/internal/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	db   *pgxpool.Pool
	repo *repositories.UserRepository
}

func NewUserService(database *pgxpool.Pool) *UserService {
	return &UserService{
		db:   database,
		repo: repositories.NewUserRepository(database),
	}
}

func (s *UserService) ListUsers(ctx context.Context, sessionId string) ([]models.User, error) {
	return s.repo.GetUsers(ctx, sessionId)
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if err := s.repo.ValidateUserData(ctx, user); err != nil {
		return nil, err
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
