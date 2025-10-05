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
			for client := range h.clients {
				select {
				case client.send <- data:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
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

var hub *Hub
var broadcast chan interface{}

func init() {
	hub = NewHub()
	broadcast = hub.GetBroadcastChannel()
	go hub.Run()
}

// GetBroadcastChannel retorna o canal de broadcast global
func GetBroadcastChannel() chan interface{} {
	return broadcast
}

// HandleConnections lida com as conexões WebSocket
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Erro ao fazer upgrade da conexão: %v", err)
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	client.hub.register <- client

	// Permitir coleta de memória de referências ao caller fazendo todo o trabalho em novas goroutines.
	go client.writePump()
	go client.readPump()
}
