package services

import (
	"context"
	"flip-planning-poker/internal/models"
	"flip-planning-poker/internal/repositories"
	"flip-planning-poker/internal/websocket"
	"fmt"
	"log"
	"strconv"

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
	if err := s.repo.CreateStory(ctx, story); err != nil {
		return nil, err
	}

	if story.Status == "ACTUAL" {
		s.repo.SetStoriesToOld(story.SessionID, story.ID)
	}
	err := s.wsService.SendSessionMessage(story.SessionID, websocket.WSMessage{
		Event: websocket.STORY_CREATED_WS_EVENT,
		Data: struct {
			Story *models.Story `json:"story"`
		}{
			Story: story,
		},
	})
	if err != nil {
		return nil, err
	}

	return story, nil
}

func (s *StoryService) RevealStory(ctx context.Context, storyId int) error {
	vote_repository := repositories.NewVoteRepository(s.db)

	story, err := s.repo.GetStoryById(ctx, storyId)
	if err != nil {
		return err
	}
	if story.Status != "ACTUAL" {
		return fmt.Errorf("story não é a atual")
	}
	if err := vote_repository.UpdateStatusPerStory(ctx, storyId); err != nil {
		return err
	}
	err = s.CalculateEstimativeAverage(ctx, story)
	if err != nil {
		return err
	}
	log.Println("Enviando evento de Story Revelada para a sessão: ", story.SessionID)
	err = s.wsService.SendSessionMessage(story.SessionID, websocket.WSMessage{
		Event: websocket.STORY_REVEALED_WS_EVENT,
		Data: struct {
			Story *models.Story `json:"story"`
		}{
			Story: story,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *StoryService) CalculateEstimativeAverage(ctx context.Context, story *models.Story) error {
	var storyVotes []models.Vote
	voteRepository := repositories.NewVoteRepository(s.db)
	err := voteRepository.FindVotes(ctx, &storyVotes, &repositories.VoteQuery{
		StoryId: story.ID,
	})
	if err != nil {
		return err
	}
	var voutesSum int
	var votesCount int
	for _, vote := range storyVotes {
		voteValue := vote.Vote
		voteInt, err := strconv.Atoi(voteValue)
		if err != nil {
			continue
		}
		voutesSum += voteInt
		votesCount += 1
	}
	var estimationAverage string
	if votesCount == 0 {
		return fmt.Errorf("nenhum voto encontrado")
	}
	if voutesSum == 0 && votesCount > 0 {
		estimationAverage = "0"
	} else {
		estimationAverage = fmt.Sprintf("%d", voutesSum/votesCount)
	}

	err = s.repo.UpdateStory(ctx, story)
	if err != nil {
		return err
	}

	story.EstimationAverage = estimationAverage
	return nil
}
