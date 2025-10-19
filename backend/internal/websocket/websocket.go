package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebsocketService struct {
	Clients           map[*websocket.Conn]string
	SessionBroadcasts map[string]chan []byte
	Mutex             *sync.Mutex
}

func NewWebsocketService() *WebsocketService {
	return &WebsocketService{
		Clients:           make(map[*websocket.Conn]string),
		Mutex:             &sync.Mutex{},
		SessionBroadcasts: make(map[string]chan []byte),
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *WebsocketService) RemoveConnAndSessionBroadcastIfEmpty(conn *websocket.Conn, sessionId string) {

	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	if conn != nil {
		delete(s.Clients, conn)
	}
	ch := s.SessionBroadcasts[sessionId]
	if ch == nil {
		return
	}
	empty := true
	for _, sid := range s.Clients {
		if sid == sessionId {
			empty = false
			break
		}
	}
	if empty {
		close(ch)
		delete(s.SessionBroadcasts, sessionId)
	}
}

func (s *WebsocketService) WsHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := r.URL.Query().Get("session_id")
	if sessionId == "" {
		http.Error(w, "missing session_id", http.StatusBadRequest)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading:", err)
		return
	}
	defer conn.Close()

	s.Mutex.Lock()
	s.Clients[conn] = sessionId
	ch := s.SessionBroadcasts[sessionId]
	if ch == nil {
		ch = make(chan []byte, 256)
		s.SessionBroadcasts[sessionId] = ch
		go s.HandleSessionMessages(sessionId)
	}
	s.Mutex.Unlock()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			s.RemoveConnAndSessionBroadcastIfEmpty(conn, sessionId)
			break
		}
		log.Println("Mensagem recebida: ", string(message))
	}
}

// Usar apenas sob Mutex.Lock()
func (s *WebsocketService) GetSessionClients(sessionId string) []*websocket.Conn {
	var sessionClients []*websocket.Conn
	for connection, session := range s.Clients {
		if session == sessionId {
			sessionClients = append(sessionClients, connection)
		}
	}
	return sessionClients
}

func (s *WebsocketService) HandleSessionMessages(sessionId string) {
	for {
		s.Mutex.Lock()
		sessionClients := s.GetSessionClients(sessionId)
		broadcast, ok := s.SessionBroadcasts[sessionId]
		s.Mutex.Unlock()
		if !ok {
			return
		}
		if len(sessionClients) == 0 || broadcast == nil {
			s.RemoveConnAndSessionBroadcastIfEmpty(nil, sessionId)
			return
		}
		message, ok := <-broadcast
		if !ok {
			return
		}

		s.Mutex.Lock()
		for _, client := range sessionClients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				client.Close()
				delete(s.Clients, client)
			}
		}
		s.Mutex.Unlock()
		s.RemoveConnAndSessionBroadcastIfEmpty(nil, sessionId)
	}
}

func (s *WebsocketService) SendSessionMessage(sessionId string, message WSMessage) error {
	s.Mutex.Lock()
	ch, ok := s.SessionBroadcasts[sessionId]
	s.Mutex.Unlock()

	if !ok || ch == nil {
		return nil
	}
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}
	select {
	case ch <- jsonMessage:
	default:
		log.Println("Buffer cheio, ignorando mensagem: ", string(jsonMessage))
	}

	return nil
}

type WebsocketEvents string

const (
	STORY_REVEALED WebsocketEvents = "STORY_REVEALED"
)

type WSMessage struct {
	Event WebsocketEvents `json:"event"`
	Data  any             `json:"data"`
}
