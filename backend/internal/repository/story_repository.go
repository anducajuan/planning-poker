package repository

import (
	"context"
	"flip-planning-poker/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StoryRepository struct {
	db *pgxpool.Pool
}

func NewStoryRepository(db *pgxpool.Pool) *StoryRepository {
	return &StoryRepository{db: db}
}

func (r *StoryRepository) FindStoryBySessionId(sessionId string) ([]model.Story, error) {
	rows, err := r.db.Query(context.Background(), "SELECT id, name, status, session_id from stories where session_id = $1", sessionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stories := []model.Story{}

	for rows.Next() {
		var s model.Story

		if err := rows.Scan(&s.ID, &s.Name, &s.Status, &s.SessionID); err != nil {
			return nil, err
		}
		stories = append(stories, s)
	}
	return stories, nil
}
