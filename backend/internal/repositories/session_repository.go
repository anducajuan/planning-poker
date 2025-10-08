package repositories

import (
	"context"
	"errors"
	"flip-planning-poker/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrInvalidID = errors.New("ID inválido")
	ErrNotFound  = errors.New("registro não encontrado")
)

type SessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) GetAll(ctx context.Context, query any) ([]models.Session, error) {
	rows, err := r.db.Query(ctx, "SELECT id, name FROM sessions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := []models.Session{}

	for rows.Next() {
		var s models.Session
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}

	return sessions, nil
}

func (r *SessionRepository) CreateSession(ctx context.Context, session *models.Session) error {
	err := r.db.QueryRow(
		ctx,
		"INSERT INTO sessions (name) VALUES ($1) RETURNING id",
		session.Name,
	).Scan(&session.ID)

	return err
}

func (r *SessionRepository) DeleteSession(ctx context.Context, id string) (int64, error) {
	result, err := r.db.Exec(ctx, "DELETE FROM sessions WHERE id=$1", id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *SessionRepository) VerifyIfNameAlreadyExists(ctx context.Context, name string) (bool, error) {
	var id string
	err := r.db.QueryRow(ctx, "SELECT id FROM sessions WHERE name=$1", name).Scan(&id)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (r *SessionRepository) ValidateSessionData(ctx context.Context, session *models.Session) error {
	if session.Name == "" {
		return errors.New("nome da sessão não pode estar vazio")
	}

	exists, err := r.VerifyIfNameAlreadyExists(ctx, session.Name)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("nome da sessão já existe")
	}

	return nil
}
