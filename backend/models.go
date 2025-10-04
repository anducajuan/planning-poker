package main

type Session struct {
	ID   string `json:"id"`   // UUID
	Name string `json:"name"` // nome da sessão
}

type User struct {
	ID        int    `json:"id"`         // SERIAL
	Name      string `json:"name"`       // nome do usuário
	SessionID string `json:"session_id"` // FK -> sessions.id
}

type Story struct {
	ID        int    `json:"id"`         // SERIAL
	Name      string `json:"name"`       // nome da história
	Status    string `json:"status"`     // "ACTUAL" ou "OLD"
	SessionID string `json:"session_id"` // FK -> sessions.id
}

type Vote struct {
	ID        int    `json:"id"`         // SERIAL
	Vote      int    `json:"vote"`       // valor do voto
	UserID    int    `json:"user_id"`    // FK -> users.id
	SessionID string `json:"session_id"` // FK -> sessions.id
	StoryID   int    `json:"story_id"`   // FK -> stories.id
}
