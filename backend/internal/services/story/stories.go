package story

import (
	"flip-planning-poker/internal/model"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func SetDB(database *pgxpool.Pool) {
	db = database
}

// func GetSessionStories(w http.ResponseWriter, r *http.Request) {
// 	sessionId := r.URL.Query().Get("session_id")

// }

// func CreateUser(w http.ResponseWriter, r *http.Request) {

// }

func validateStoryData(w http.ResponseWriter, user *model.Story) {

}

// func getStories
