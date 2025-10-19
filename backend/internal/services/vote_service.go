package services

import (
	"context"
	"flip-planning-poker/internal/models"
	"flip-planning-poker/internal/repositories"
	"flip-planning-poker/internal/websocket"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VoteService struct {
	db        *pgxpool.Pool
	repo      *repositories.VoteRepository
	wsService *websocket.WebsocketService
}

func NewVoteService(database *pgxpool.Pool, websocketService *websocket.WebsocketService) *VoteService {
	return &VoteService{
		db:        database,
		repo:      repositories.NewVoteRepository(database),
		wsService: websocketService,
	}
}

func (s *VoteService) Create(ctx context.Context, vote *models.Vote) (*models.Vote, error) {
	if vote.Status == "" {
		vote.Status = "HIDDEN"
	}
	if !vote.Status.IsValid() {
		return nil, fmt.Errorf("status inválido: %s", vote.Status)
	}

	existUserVote, err := s.verifyExistingVoteForStory(ctx, vote.UserID, vote.StoryID)
	if err != nil {
		return nil, err
	}

	if existUserVote {
		return nil, fmt.Errorf("já existe um vote para o usuário %d na story %d", vote.UserID, vote.StoryID)
	}

	if err := s.repo.CreateVote(ctx, vote); err != nil {
		return nil, err
	}
	log.Println("Enviando evento de criação de voto via websocket para a sessão: ", vote.SessionID)
	s.wsService.SendSessionMessage(vote.SessionID, websocket.WSMessage{
		Event: websocket.VOTE_CREATED_WS_EVENT,
		Data: struct {
			Vote *models.Vote `json:"vote"`
		}{
			Vote: vote,
		},
	})

	return vote, nil
}

func (s *VoteService) verifyExistingVoteForStory(ctx context.Context, userId int, storyId int) (bool, error) {

	var votes []models.Vote
	if err := s.repo.FindVotes(ctx, &votes, &repositories.VoteQuery{
		UserId:  userId,
		StoryId: storyId,
	}); err != nil {
		return false, err
	}
	if len(votes) > 0 {
		return true, nil
	}
	return false, nil
}

func (s *VoteService) List(ctx context.Context, storyId int) ([]models.Vote, error) {
	query := repositories.VoteQuery{
		StoryId: storyId,
	}

	var votes []models.Vote
	if err := s.repo.FindVotes(ctx, &votes, &query); err != nil {
		return nil, err
	}

	return votes, nil
}

func (s *VoteService) Patch(ctx context.Context, id int, patch *repositories.VotePatch) error {

	vote, err := s.repo.UpdateVote(ctx, id, patch)
	if err != nil {
		return err
	}

	log.Println("Enviando evento de voto alterado para a sessão: ", vote.SessionID)
	err = s.wsService.SendSessionMessage(vote.SessionID, websocket.WSMessage{
		Event: websocket.VOTE_CHANGED_WS_EVENT,
		Data: struct {
			Vote *models.Vote `json:"vote"`
		}{
			Vote: vote,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *VoteService) Get(ctx context.Context, id int) (*models.Vote, error) {
	return s.repo.GetVoteByID(ctx, id)
}
