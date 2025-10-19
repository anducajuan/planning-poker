package services

import (
	"context"
	"flip-planning-poker/internal/models"
	"flip-planning-poker/internal/repositories"
	"flip-planning-poker/internal/websocket"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StoryService struct {
	db        *pgxpool.Pool
	repo      *repositories.StoryRepository
	wsService *websocket.WebsocketService
}

func NewStoryService(database *pgxpool.Pool, websocketService *websocket.WebsocketService) *StoryService {
	return &StoryService{
		db:        database,
		repo:      repositories.NewStoryRepository(database),
		wsService: websocketService,
	}
}

func (s *StoryService) ListStories(ctx context.Context, sessionId string) ([]models.Story, error) {
	if sessionId == "" {
		return nil, repositories.ErrInvalidID
	}

	return s.repo.FindStoryBySessionId(sessionId)
}

func (s *StoryService) CreateStory(ctx context.Context, story *models.Story) (*models.Story, error) {
	if err := s.repo.CreateStory(story); err != nil {
		return nil, err
	}

	if story.Status == "ACTUAL" {
		s.repo.SetStoriesToOld(story.SessionID, story.ID)
	}

	return story, nil
}

func (s *StoryService) RevealStory(ctx context.Context, storyId int) error {
	vote_repository := repositories.NewVoteRepository(s.db)

	if err := vote_repository.UpdateStatusPerStory(ctx, storyId); err != nil {
		return err
	}
	story, err := s.repo.GetStoryById(ctx, storyId)
	if err != nil {
		return err
	}
	log.Println("Enviando evento de Story Revelada para a sessão: ", story.SessionID)
	err = s.wsService.SendSessionMessage(story.SessionID, websocket.WSMessage{
		Event: websocket.STORY_REVEALED,
		Data:  nil,
	})
	if err != nil {
		return err
	}
	return nil
}
