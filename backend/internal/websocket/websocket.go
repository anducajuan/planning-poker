package websocket

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebsocketService struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan []byte
	Mutex     *sync.Mutex
}

func NewWebsocketService() *WebsocketService {
	return &WebsocketService{
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan []byte, 256),
		Mutex:     &sync.Mutex{},
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *WebsocketService) WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	defer conn.Close()

	s.Mutex.Lock()
	s.Clients[conn] = true
	s.Mutex.Unlock()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			s.Mutex.Lock()
			delete(s.Clients, conn)
			s.Mutex.Unlock()
			break
		}
		s.Broadcast <- message
	}
}

func (s *WebsocketService) HandleMessages() {
	for {

		message := <-s.Broadcast

		s.Mutex.Lock()
		for client := range s.Clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				client.Close()
				delete(s.Clients, client)
			}
		}
		s.Mutex.Unlock()
	}
}
