package websocket

import "time"

// MessageType define os tipos de mensagem do WebSocket
type MessageType string

const (
	// Tipos de mensagem para sessões
	SessionCreated MessageType = "session_created"
	SessionDeleted MessageType = "session_deleted"
	SessionUpdated MessageType = "session_updated"

	// Tipos de mensagem para usuários
	UserJoined MessageType = "user_joined"
	UserLeft   MessageType = "user_left"

	// Tipos de mensagem para stories
	StoryCreated MessageType = "story_created"
	StoryUpdated MessageType = "story_updated"

	// Tipos de mensagem para votos
	VoteCreated MessageType = "vote_created"
	VoteUpdated MessageType = "vote_updated"
	VotingEnded MessageType = "voting_ended"

	// Tipos de mensagem de sistema
	SystemMessage MessageType = "system_message"
	Error         MessageType = "error"
)

// Message representa uma mensagem WebSocket estruturada
type Message struct {
	Type      MessageType `json:"type"`
	Data      interface{} `json:"data"`
	SessionID string      `json:"session_id,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	UserID    int         `json:"user_id,omitempty"`
}

// NewMessage cria uma nova mensagem com timestamp
func NewMessage(msgType MessageType, data interface{}) *Message {
	return &Message{
		Type:      msgType,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// WithSession adiciona o ID da sessão à mensagem
func (m *Message) WithSession(sessionID string) *Message {
	m.SessionID = sessionID
	return m
}

// WithUser adiciona o ID do usuário à mensagem
func (m *Message) WithUser(userID int) *Message {
	m.UserID = userID
	return m
}

type SessionData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserData struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	SessionID string `json:"session_id"`
}

type StoryData struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	SessionID string `json:"session_id"`
}

type VoteData struct {
	ID        int    `json:"id"`
	Vote      int    `json:"vote"`
	UserID    int    `json:"user_id"`
	SessionID string `json:"session_id"`
	StoryID   int    `json:"story_id"`
	Status    string `json:"status"`
}
