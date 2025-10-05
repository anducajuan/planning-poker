package user

import (
	"encoding/json"
	"flip-planning-poker/internal/model"
	"flip-planning-poker/internal/repository"
	"flip-planning-poker/internal/utils"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	db   *pgxpool.Pool
	repo *repository.UserRepository
}

func NewUserService(database *pgxpool.Pool) *UserService {
	return &UserService{
		db:   database,
		repo: repository.NewUserRepository(database),
	}
}

func (s *UserService) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionId := r.URL.Query().Get("session_id")

	users, err := s.repo.GetUsers(ctx, sessionId)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err, "Erro ao buscar usuários")
		return
	}

	utils.SendSuccessWithTotal(w, http.StatusOK, users, len(users), "Usuários encontrados com sucesso")
}

func (s *UserService) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	var u model.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		utils.SendError(w, http.StatusBadRequest, err, "Dados inválidos no corpo da requisição")
		return
	}

	if err := s.repo.ValidateUserData(ctx, &u); err != nil {
		utils.SendError(w, http.StatusBadRequest, err, err.Error())
		return
	}

	err := s.repo.CreateUser(ctx, &u)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err, "Erro ao criar usuário")
		return
	}

	utils.SendSuccess(w, http.StatusCreated, u, "Usuário criado com sucesso")
}
