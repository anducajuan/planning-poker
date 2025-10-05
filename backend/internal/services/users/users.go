package users

import (
	"context"
	"encoding/json"
	"flip-planning-poker/internal/utils"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func SetDB(database *pgxpool.Pool) {
	db = database
}

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	SessionID string `json:"session_id"`
}

func GetUsers(w http.ResponseWriter, r *http.Request) {

	sessionId := r.URL.Query().Get("session_id")

	selectQuery := "SELECT id, name, session_id FROM users"
	var rows pgx.Rows
	var err error

	if sessionId != "" {
		selectQuery += " WHERE session_id = $1"
		rows, err = db.Query(context.Background(), selectQuery, sessionId)
	} else {
		rows, err = db.Query(context.Background(), selectQuery)
	}
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err, "Erro ao buscar usuários")
		return
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.SessionID); err != nil {
			utils.SendError(w, http.StatusInternalServerError, err, "Erro ao buscar dados dos usuários")
		}
		users = append(users, u)
	}

	utils.SendSuccessWithTotal(w, http.StatusOK, users, len(users), "Usuários encontrados com sucesso")

}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var u User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		utils.SendError(w, http.StatusBadRequest, err, "Dados inválidos no corpo da requisição")
		return
	}

	validateUserData(w, &u)

	err := db.QueryRow(
		context.Background(),
		"INSERT INTO users (name, session_id) VALUES ($1, $2) RETURNING id",
		u.Name,
		u.SessionID,
	).Scan(&u.ID)

	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err, "Erro ao criar usuário")
		return
	}

	utils.SendSuccess(w, http.StatusCreated, u, "Usuário criado com sucesso")
}

func verifyIfNameAlreadyExists(name string) bool {
	row := db.QueryRow(context.Background(), "SELECT id FROM users WHERE name=$1", name)
	var id string
	if err := row.Scan(&id); err != nil {
		return false
	}

	return id != ""
}

func validateUserData(w http.ResponseWriter, user *User) {
	if user.Name == "" {
		utils.SendError(w, http.StatusBadRequest, nil, "Nome do usuário é obrigatório")
		return
	}

	if user.SessionID == "" {
		utils.SendError(w, http.StatusBadRequest, nil, "ID da sessão é obrigatório")
		return
	}

	if verifyIfNameAlreadyExists(user.Name) {
		utils.SendError(w, http.StatusBadRequest, nil, "Nome do usuário já existe")
		return
	}
}
