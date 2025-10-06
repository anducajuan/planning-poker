package websocket

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketService gerencia as operações de WebSocket
type WebSocketService struct {
	hub      *Hub
	mu       sync.RWMutex
	upgrader websocket.Upgrader
}

// NewWebSocketService cria um novo serviço de WebSocket
func NewWebSocketService() *WebSocketService {
	hub := NewHub()
	service := &WebSocketService{
		hub: hub,
		upgrader: websocket.Upgrader{
			CheckOrigin:     func(r *http.Request) bool { return true },
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}

	// Iniciar o hub em uma goroutine
	go hub.Run()

	return service
}

// GetHub retorna o hub para uso direto (se necessário)
func (ws *WebSocketService) GetHub() *Hub {
	return ws.hub
}

// BroadcastToAll envia uma mensagem para todos os clientes conectados
func (ws *WebSocketService) BroadcastToAll(message *Message) error {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	select {
	case ws.hub.broadcast <- message:
		return nil
	default:
		log.Printf("Canal de broadcast bloqueado, mensagem descartada")
		return nil
	}
}

// BroadcastToSession envia uma mensagem para todos os clientes de uma sessão específica
func (ws *WebSocketService) BroadcastToSession(sessionID string, message *Message) error {
	// Por enquanto, broadcast para todos (pode ser melhorado para filtrar por sessão)
	message.SessionID = sessionID
	return ws.BroadcastToAll(message)
}

// NotifySessionCreated notifica sobre criação de sessão
func (ws *WebSocketService) NotifySessionCreated(sessionData SessionData) error {
	message := NewMessage(SessionCreated, sessionData).WithSession(sessionData.ID)
	return ws.BroadcastToAll(message)
}

// NotifySessionDeleted notifica sobre deleção de sessão
func (ws *WebSocketService) NotifySessionDeleted(sessionID string) error {
	message := NewMessage(SessionDeleted, map[string]string{"session_id": sessionID}).WithSession(sessionID)
	return ws.BroadcastToAll(message)
}

// NotifyUserJoined notifica sobre usuário que entrou na sessão
func (ws *WebSocketService) NotifyUserJoined(userData UserData) error {
	message := NewMessage(UserJoined, userData).WithSession(userData.SessionID).WithUser(userData.ID)
	return ws.BroadcastToSession(userData.SessionID, message)
}

// NotifyStoryCreated notifica sobre criação de story
func (ws *WebSocketService) NotifyStoryCreated(storyData StoryData) error {
	message := NewMessage(StoryCreated, storyData).WithSession(storyData.SessionID)
	return ws.BroadcastToSession(storyData.SessionID, message)
}

// NotifyVoteCreated notifica sobre criação de voto
func (ws *WebSocketService) NotifyVoteCreated(voteData VoteData) error {
	message := NewMessage(VoteCreated, voteData).WithSession(voteData.SessionID).WithUser(voteData.UserID)
	return ws.BroadcastToSession(voteData.SessionID, message)
}

// NotifyError envia uma mensagem de erro
func (ws *WebSocketService) NotifyError(sessionID string, errorMsg string) error {
	message := NewMessage(Error, map[string]string{"message": errorMsg}).WithSession(sessionID)
	return ws.BroadcastToSession(sessionID, message)
}

// GetConnectedClientsCount retorna o número de clientes conectados
func (ws *WebSocketService) GetConnectedClientsCount() int {
	ws.hub.mu.RLock()
	defer ws.hub.mu.RUnlock()
	return len(ws.hub.clients)
}

// Shutdown para um shutdown graceful
func (ws *WebSocketService) Shutdown(ctx context.Context) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	log.Println("Iniciando shutdown do WebSocket service...")

	// Fechar todas as conexões de clientes
	ws.hub.mu.Lock()
	for client := range ws.hub.clients {
		close(client.send)
		client.conn.Close()
		delete(ws.hub.clients, client)
	}
	ws.hub.mu.Unlock()

	log.Println("WebSocket service finalizado")
	return nil
}

// Broadcaster interface para facilitar testes
type Broadcaster interface {
	BroadcastToAll(message *Message) error
	BroadcastToSession(sessionID string, message *Message) error
	NotifySessionCreated(sessionData SessionData) error
	NotifySessionDeleted(sessionID string) error
	NotifyUserJoined(userData UserData) error
	NotifyStoryCreated(storyData StoryData) error
	NotifyVoteCreated(voteData VoteData) error
}

// Garantir que WebSocketService implementa Broadcaster
var _ Broadcaster = (*WebSocketService)(nil)
