package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func getSessions(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.Query(context.Background(), "SELECT id, content FROM sessions")
	defer rows.Close()

	var msgs []Session
	for rows.Next() {
		var m Session
		rows.Scan(&m.ID, &m.Name)
		msgs = append(msgs, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msgs)
}

func createSession(w http.ResponseWriter, r *http.Request) {
	var m Session
	json.NewDecoder(r.Body).Decode(&m)

	err := db.QueryRow(
		context.Background(),
		"INSERT INTO session (content) VALUES ($1) RETURNING id",
		m.Name,
	).Scan(&m.ID)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 🔥 Notifica websocket
	broadcast <- m

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

func deleteSession(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_, err := db.Exec(context.Background(), "DELETE FROM sessions WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
