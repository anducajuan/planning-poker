package services

import (
	"context"
	"flip-planning-poker/internal/models"
	"flip-planning-poker/internal/repositories"

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

func (s *VoteService) CreateVote(ctx context.Context, vote *models.Vote) (*models.Vote, error) {
	if err := s.repo.CreateVote(ctx, vote); err != nil {
		return nil, err
	}

	return vote, nil
}

func (s *VoteService) ListVotes(ctx context.Context, storyId string) ([]models.Vote, error) {
	query := repositories.VoteQuery{
		StoryId: storyId,
	}

	var votes []models.Vote
	if err := s.repo.FindVotes(ctx, &votes, &query); err != nil {
		return nil, err
	}

	return votes, nil
}
