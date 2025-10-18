package services

import (
	"context"
	"flip-planning-poker/internal/models"
	"flip-planning-poker/internal/repositories"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VoteService struct {
	db   *pgxpool.Pool
	repo *repositories.VoteRepository
}

func NewVoteService(database *pgxpool.Pool) *VoteService {
	return &VoteService{
		db:   database,
		repo: repositories.NewVoteRepository(database),
	}
}

func (s *VoteService) Create(ctx context.Context, vote *models.Vote) (*models.Vote, error) {
	if vote.Status == "" {
		vote.Status = "HIDDEN"
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
	return s.repo.UpdateVote(ctx, id, patch)
}

func (s *VoteService) Get(ctx context.Context, id int) (*models.Vote, error) {
	return s.repo.GetVoteByID(ctx, id)
}
