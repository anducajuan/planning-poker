package repository

import (
	"context"
	"errors"
	"flip-planning-poker/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUsers(ctx context.Context, sessionID string) ([]model.User, error) {
	selectQuery := "SELECT id, name, session_id FROM users"
	var args []interface{}

	if sessionID != "" {
		selectQuery += " WHERE session_id = $1"
		args = append(args, sessionID)
	}

	rows, err := r.db.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []model.User{}

	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.SessionID); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	err := r.db.QueryRow(
		ctx,
		"INSERT INTO users (name, session_id) VALUES ($1, $2) RETURNING id",
		user.Name,
		user.SessionID,
	).Scan(&user.ID)

	return err
}

func (r *UserRepository) VerifyIfNameAlreadyExists(ctx context.Context, name string, session string) (bool, error) {
	var id int
	err := r.db.QueryRow(ctx, "SELECT id FROM users WHERE name=$1 AND session_id=$2", name, session).Scan(&id)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (r *UserRepository) ValidateUserData(ctx context.Context, user *model.User) error {
	if user.Name == "" {
		return errors.New("nome do usuário é obrigatório")
	}

	if user.SessionID == "" {
		return errors.New("ID da sessão é obrigatório")
	}

	exists, err := r.VerifyIfNameAlreadyExists(ctx, user.Name, user.SessionID)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("nome do usuário já existe")
	}

	return nil
}
