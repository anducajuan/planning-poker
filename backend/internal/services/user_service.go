package services

import (
	"context"
	"flip-planning-poker/internal/models"
	"flip-planning-poker/internal/repositories"
	"flip-planning-poker/internal/websocket"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	db        *pgxpool.Pool
	repo      *repositories.UserRepository
	wsService *websocket.WebsocketService
}

func NewUserService(database *pgxpool.Pool, websocketService *websocket.WebsocketService) *UserService {
	return &UserService{
		db:        database,
		repo:      repositories.NewUserRepository(database),
		wsService: websocketService,
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

	log.Println("Enviando evento de entrada de usuário para a sessão: ", user.SessionID)
	s.wsService.SendSessionMessage(user.SessionID, websocket.WSMessage{
		Event: websocket.USER_JOINED_WS_EVENT,
		Data: struct {
			User *models.User `json:"user"`
		}{
			User: user,
		},
	})

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userId int) (*models.User, error) {

	removedUser := &models.User{ID: userId}

	err := s.repo.DeleteUser(ctx, removedUser)
	if err != nil {
		return nil, err
	}

	log.Println("Enviando evento de saída de usuário para a sessão: ", removedUser.SessionID)
	s.wsService.SendSessionMessage(removedUser.SessionID, websocket.WSMessage{
		Event: websocket.USER_LEFT_WS_EVENT,
		Data: struct {
			User *models.User `json:"user"`
		}{
			User: removedUser,
		},
	})

	return removedUser, nil
}
