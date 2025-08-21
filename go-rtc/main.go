package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-rtc/chat"
	"log"
	"net/http"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")
	chatOrchestrator := NewChatOrcWrapper()

	http.HandleFunc("/chat/create/", withCORS(chatOrchestrator.CreateChatHandler))
	http.HandleFunc("/chat/", withCORS(chatOrchestrator.GetChatHandler))
	http.HandleFunc("/chat/join/", chatOrchestrator.JoinChatHandler)
	http.HandleFunc("/chats", withCORS(chatOrchestrator.GetChatsHandler))

	http.HandleFunc("/chat/test", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Logging all chats!")
		chats, _ := chatOrchestrator.orchestrator.GetChats()
		for _, chatName := range chats {
			log.Printf("Chat: %v", chatName)
		}
	})
	fmt.Println("Server running on :8080")

	http.ListenAndServe(":8080", nil)
}

type ChatOrcWrapper struct {
	orchestrator *chat.ChatOrchestrator
}

func NewChatOrcWrapper() *ChatOrcWrapper {
	return &ChatOrcWrapper{orchestrator: chat.NewChatOrchestrator()}
}

func (cow ChatOrcWrapper) CreateChatHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Failed to read chat name from url", http.StatusBadRequest)
		return
	}
	roomName := parts[3]
	err := cow.orchestrator.CreateChat(roomName)
	if err != nil {
		var chatErr *chat.ApiError
		if errors.As(err, &chatErr) {
			http.Error(w, chatErr.Message, chatErr.StatusCode)
		} else {
			log.Printf("faed t create chat, error: %v", err)
			http.Error(w, "Internal server error", http.StatusBadRequest)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}
func withCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4321")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed to actual handler
		handler(w, r)
	}
}

func (cow ChatOrcWrapper) GetChatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
	}
	chats, err := cow.orchestrator.GetChats()
	if err != nil {
		log.Fatalf("Failed to get chats, error: %v", err)
		http.Error(w, "Failed to get chats", http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(chatsResponse{Chats: chats})
	if err != nil {
		log.Fatalf("Failed to marshall response, error: %v", err)
		http.Error(w, "Failed to get chats", http.StatusInternalServerError)
	}

	w.Write(b)
}

type chatsResponse struct {
	Chats []string `json:"chats"`
}

func (cow ChatOrcWrapper) GetChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid Method", http.StatusBadRequest)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	chatInfo, err := cow.orchestrator.GetChatInformation(parts[2])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := json.Marshal(chatInfo)
	if err != nil {
		log.Fatal("Failed to marshal json in request")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (cow ChatOrcWrapper) JoinChatHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		log.Fatalf("error invalid url")
		return
	}
	roomName := parts[3]
	clientName := r.URL.Query().Get("clientName")

	client, err := cow.orchestrator.JoinChat(roomName, clientName, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to join server, error: %v", err.Error())
		return
	}
	// Make sure client's and thereby rooms are closed
	defer client.Disconnect()

	// Keep client connection open
	log.Printf("Listening on connection with %v", client.Name)

	ctx := context.Background()

	// Listen for disonnection
	for {
		message, err := client.ReadClient(ctx)
		if err != nil {
			log.Printf("Client failed to read, error: %v", err)
			return
		}

		log.Printf("Received message from client: '%v'", message)
	}
}
