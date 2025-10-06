package repository

import (
	"context"
	"flip-planning-poker/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VoteRepository struct {
	db *pgxpool.Pool
}

func NewVoteRepository(db *pgxpool.Pool) *VoteRepository {
	return &VoteRepository{db: db}
}

func (r *VoteRepository) CreateVote(ctx context.Context, v *model.Vote) error {
	if v.Status == "" {
		v.Status = "HIDDEN"
	}
	insertStatement := "insert into votes (vote, user_id, session_id, story_id, status) values($1, $2, $3, $4,$5 ) returning id"

	err := r.db.QueryRow(ctx, insertStatement, v.Vote, v.UserID, v.SessionID, v.StoryID, v.Status).Scan(&v.ID)

	return err

}
