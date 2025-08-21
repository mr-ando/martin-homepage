package chat

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/coder/websocket"
)

type Message struct {
	Sender    string `json:"sender"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type chatRoom struct {
	Name         string
	Clients      []*Client
	messages     []Message
	orchestrator *ChatOrchestrator

	deleteTimer *time.Timer
	mu          sync.RWMutex
}

type ChatInformation struct {
	ChatName     string    `json:"chatName"`
	Messages     []Message `json:"messages"`
	Participants []string  `json:"participants"`
}

func newRoom(chatName string, orchestrator *ChatOrchestrator) *chatRoom {
	room := &chatRoom{
		Name:         chatName,
		Clients:      []*Client{},
		orchestrator: orchestrator,
		messages:     []Message{},
		mu:           sync.RWMutex{},
	}
	room.scheduleDelete()
	return room
}

func (r *chatRoom) RemoveClient(clientToRemove *Client) {
	r.Clients = slices.DeleteFunc(r.Clients, func(c *Client) bool { return c == clientToRemove })
	if len(r.Clients) == 0 {
		r.scheduleDelete()
	}
}

func (r *chatRoom) scheduleDelete() {
	r.cancelDelete()

	r.deleteTimer = time.AfterFunc(30*time.Second, func() {
		r.mu.Lock()
		defer r.mu.Unlock()

		if len(r.Clients) == 0 {
			r.orchestrator.deleteRoom(r)
		}
	})
}

func (r *chatRoom) cancelDelete() {
	if r.deleteTimer != nil {
		r.deleteTimer.Stop()
		r.deleteTimer = nil
	}
}

func (room *chatRoom) Connect(clientName string, w http.ResponseWriter, r *http.Request) (*Client, error) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{
			"localhost:4321",
			"http://localhost:4321/",
		},
	})
	if err != nil {
		log.Fatalf("Failed to connect to chatroom, error: %v", err)
		return nil, err
	}
	client := newClient(clientName, c, room)
	room.Clients = append(room.Clients, client)
	return client, nil
}

func (room *chatRoom) GetChatInformation() *ChatInformation {
	clients := make([]string, len(room.Clients))

	for i, c := range room.Clients {
		clients[i] = c.Name
	}
	return &ChatInformation{ChatName: room.Name, Messages: room.messages, Participants: clients}
}

func (room *chatRoom) PublishMessage(m Message, ctx context.Context) {
	jsonMessage, err := json.Marshal(m)
	if err != nil {
		log.Fatalf("Failed to marshal json message, error: %v", err)
	}
	for _, client := range room.Clients {
		go func(c *Client) {
			err := client.conn.Write(ctx, websocket.MessageText, jsonMessage)
			if err != nil {
				log.Printf("Failed to write to client: %v\n", err)
			}
		}(client)
	}
	room.messages = append(room.messages, m)
}
