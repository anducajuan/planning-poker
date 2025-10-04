package session

import (
	"context"
	"encoding/json"
	"flip-planning-poker/utils"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Variáveis globais que serão injetadas
var db *pgxpool.Pool
var broadcast chan interface{}

type Session struct {
	ID   string `json:"id"`   // UUID
	Name string `json:"name"` // nome da sessão
}

type Response struct {
	Data  interface{} `json:"data"`
	Total int         `json:"total"`
}

// SetDB injeta a conexão do banco de dados
func SetDB(database *pgxpool.Pool) {
	db = database
}

// SetBroadcast injeta o canal de broadcast
func SetBroadcast(broadcastChan chan interface{}) {
	broadcast = broadcastChan
}

func GetSessions(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(context.Background(), "SELECT id, name FROM sessions")
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err, "Erro ao buscar sessões")
		return
	}
	defer rows.Close()

	sessions := []Session{}

	for rows.Next() {
		var s Session
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			utils.SendError(w, http.StatusInternalServerError, err, "Erro ao processar dados das sessões")
			return
		}
		sessions = append(sessions, s)
	}

	utils.SendSuccessWithTotal(w, http.StatusOK, sessions, len(sessions), "Sessões recuperadas com sucesso")
}

func CreateSession(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var s Session

	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		utils.SendError(w, http.StatusBadRequest, err, "Dados inválidos no corpo da requisição")
		return
	}

	if s.Name == "" {
		utils.SendError(w, http.StatusBadRequest, nil, "Nome da sessão não pode estar vazio")
		return
	}
	if verifyIfNameAlreadyExists(s.Name) {
		utils.SendError(w, http.StatusBadRequest, nil, "Nome da sessão já existe")
		return
	}

	err := db.QueryRow(
		context.Background(),
		"INSERT INTO sessions (name) VALUES ($1) RETURNING id",
		s.Name,
	).Scan(&s.ID)

	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err, "Erro ao criar sessão")
		return
	}

	// Enviar para o canal de broadcast
	if broadcast != nil {
		select {
		case broadcast <- s:
		default:
			// Canal bloqueado, não bloquear a resposta
		}
	}

	utils.SendSuccess(w, http.StatusCreated, s, "Sessão criada com sucesso")
}

func DeleteSession(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id == "" {
		utils.SendError(w, http.StatusBadRequest, nil, "ID da sessão é obrigatório")
		return
	}

	result, err := db.Exec(context.Background(), "DELETE FROM sessions WHERE id=$1", id)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err, "Erro ao deletar sessão")
		return
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		utils.SendError(w, http.StatusNotFound, nil, "Sessão não encontrada")
		return
	}

	utils.SendSuccess(w, http.StatusOK, nil, "Sessão deletada com sucesso")
}

func verifyIfNameAlreadyExists(name string) bool {
	row := db.QueryRow(context.Background(), "SELECT id FROM sessions WHERE name=$1", name)
	var id string
	if err := row.Scan(&id); err != nil {
		return false
	}

	return id != ""
}
