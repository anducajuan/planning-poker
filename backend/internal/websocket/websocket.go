package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Hub mantém o conjunto de clientes ativos e transmite mensagens para os clientes.
type Hub struct {
	// Clientes registrados.
	clients map[*Client]bool

	// Canal de broadcast de mensagens.
	broadcast chan interface{}

	// Registrar requisições dos clientes.
	register chan *Client

	// Cancelar registro de requisições dos clientes.
	unregister chan *Client

	// Mutex para proteger o mapa de clientes
	mu sync.RWMutex
}

// Client é um intermediário entre o websocket connection e o hub.
type Client struct {
	hub *Hub

	// A conexão websocket.
	conn *websocket.Conn

	// Canal de mensagens de saída.
	send chan []byte
}

// NewHub cria um novo Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan interface{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run inicia o hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Cliente conectado. Total: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("Cliente desconectado. Total: %d", len(h.clients))

		case message := <-h.broadcast:
			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("Erro ao serializar mensagem: %v", err)
				continue
			}

			h.mu.RLock()
			clientsToRemove := []*Client{}
			for client := range h.clients {
				select {
				case client.send <- data:
					// Mensagem enviada com sucesso
				default:
					// Cliente não conseguiu receber, marcar para remoção
					clientsToRemove = append(clientsToRemove, client)
				}
			}
			h.mu.RUnlock()

			// Remover clientes desconectados
			if len(clientsToRemove) > 0 {
				h.mu.Lock()
				for _, client := range clientsToRemove {
					if _, ok := h.clients[client]; ok {
						close(client.send)
						delete(h.clients, client)
					}
				}
				h.mu.Unlock()
				log.Printf("Removidos %d clientes desconectados", len(clientsToRemove))
			}
		}
	}
}

// GetBroadcastChannel retorna o canal de broadcast
func (h *Hub) GetBroadcastChannel() chan interface{} {
	return h.broadcast
}

// readPump bombeia mensagens do websocket connection para o hub.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Erro no websocket: %v", err)
			}
			break
		}
	}
}

// writePump bombeia mensagens do hub para o websocket connection.
func (c *Client) writePump() {
	defer c.conn.Close()

	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Erro ao escrever mensagem: %v", err)
			return
		}
	}
}

// Instância global para compatibilidade (será removida posteriormente)
var globalService *WebSocketService

func init() {
	globalService = NewWebSocketService()
}

// GetBroadcastChannel retorna o canal de broadcast global (para compatibilidade)
func GetBroadcastChannel() chan interface{} {
	if globalService != nil {
		return globalService.hub.broadcast
	}
	return nil
}

// GetGlobalService retorna a instância global do WebSocketService
func GetGlobalService() *WebSocketService {
	return globalService
}

// HandleConnections lida com as conexões WebSocket
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	if globalService == nil {
		log.Printf("WebSocket service não inicializado")
		http.Error(w, "WebSocket service não disponível", http.StatusInternalServerError)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Erro ao fazer upgrade da conexão: %v", err)
		return
	}

	client := &Client{
		hub:  globalService.hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	client.hub.register <- client

	// Permitir coleta de memória de referências ao caller fazendo todo o trabalho em novas goroutines.
	go client.writePump()
	go client.readPump()
}
