package users

import (
	"encoding/json"
	"flip-planning-poker/utils"
	"net/http"

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

func CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var u User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		utils.SendError(w, http.StatusBadRequest, err, "Dados inválidos no corpo da requisição")
		return
	}

}
