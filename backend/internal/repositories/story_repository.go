package repositories

import (
	"context"
	"errors"
	"flip-planning-poker/internal/models"
	"flip-planning-poker/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StoryRepository struct {
	db *pgxpool.Pool
}

func NewStoryRepository(db *pgxpool.Pool) *StoryRepository {
	return &StoryRepository{db: db}
}

func (r *StoryRepository) FindStoryBySessionId(sessionId string) ([]models.Story, error) {
	rows, err := r.db.Query(context.Background(), "SELECT id, name, status, session_id, estimation_average from stories where session_id = $1", sessionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stories := []models.Story{}

	for rows.Next() {
		var s models.Story

		if err := rows.Scan(&s.ID, &s.Name, &s.Status, &s.SessionID, &s.EstimationAverage); err != nil {
			return nil, err
		}
		stories = append(stories, s)
	}
	return stories, nil
}

func (r *StoryRepository) CreateStory(ctx context.Context, story *models.Story) error {
	validateErrs := validateStoryData(story)
	if len(validateErrs) > 0 {
		return validateErrs[0]
	}
	err := r.db.QueryRow(ctx, `
		INSERT INTO 
			stories (name, session_id, status) 
			VALUES ($1, $2, $3) 
			returning id`,
		story.Name, story.SessionID, story.Status).Scan(&story.ID)

	return err
}

func (r *StoryRepository) GetStoryById(ctx context.Context, storyId int) (*models.Story, error) {
	var story models.Story
	query := "select s.id, s.name, s.status, s.session_id, coalesce(s.estimation_average, '') from stories s where s.id = $1"
	err := r.db.QueryRow(ctx, query, storyId).Scan(&story.ID, &story.Name, &story.Status, &story.SessionID, &story.EstimationAverage)
	if err != nil {
		return nil, err
	}
	return &story, nil
}

func validateStoryData(story *models.Story) []error {
	var errs []error
	possibleStatus := []string{
		"ACTUAL",
		"OLD",
	}
	if story.Name == "" {
		errs = append(errs, errors.New("name cannot be an empty string"))
		return errs
	}
	if story.SessionID == "" {
		errs = append(errs, errors.New("session_id cannot be an empty string"))
		return errs
	}
	if story.Status == "" {
		errs = append(errs, errors.New("status must be ACTUAL or OLD"))
		return errs
	}
	if !utils.ContainsString(possibleStatus, story.Status) {
		errs = append(errs, errors.New("status must be ACTUAL or OLD"))
		return errs
	}
	return errs
}

func (r *StoryRepository) SetStoriesToOld(sessionID string, storyToKeepID int) error {

	updateStatement := `
		update stories
		set status = 'OLD'
		where session_id = $1
		and id <> $2 
		and status <> 'OLD'
	`

	_, err := r.db.Exec(context.Background(), updateStatement, sessionID, storyToKeepID)
	if err != nil {
		utils.Logger("Erro ao atualizar stories", err)
		return err
	}

	return nil
}
