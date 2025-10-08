package repositories

import (
	"context"
	"flip-planning-poker/internal/models"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VoteRepository struct {
	db *pgxpool.Pool
}

func NewVoteRepository(db *pgxpool.Pool) *VoteRepository {
	return &VoteRepository{db: db}
}

func (r *VoteRepository) CreateVote(ctx context.Context, v *models.Vote) error {
	if v.Status == "" {
		v.Status = "HIDDEN"
	}
	insertStatement := "insert into votes (vote, user_id, session_id, story_id, status) values($1, $2, $3, $4,$5 ) returning id"

	err := r.db.QueryRow(ctx, insertStatement, v.Vote, v.UserID, v.SessionID, v.StoryID, v.Status).Scan(&v.ID)

	return err
}

type VoteQuery struct {
	StoryId string
	Status  string
}

func (r *VoteRepository) FindVotes(ctx context.Context, v *[]models.Vote, q *VoteQuery) error {
	baseQuery := `
	SELECT v.id, v.vote, v.user_id, v.session_id, v.story_id, v.status 
	FROM votes v 
	`
	var conditions []string
	var args []any

	argIndex := 1

	if q.Status != "" {
		conditions = append(conditions, fmt.Sprintf("v.status = $%d", argIndex))
		args = append(args, q.Status)
		argIndex++
	}
	if q.StoryId != "" {
		conditions = append(conditions, fmt.Sprintf("v.story_id = $%d", argIndex))
		args = append(args, q.StoryId)
		argIndex++
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := r.db.Query(ctx, baseQuery, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var vote models.Vote
		if err := rows.Scan(&vote.ID, &vote.Vote, &vote.UserID, &vote.SessionID, &vote.StoryID, &vote.Status); err != nil {
			return err
		}
		*v = append(*v, vote)
	}

	return rows.Err()
}
