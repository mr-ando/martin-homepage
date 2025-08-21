package chat

import (
	"errors"
	"log"
	"net/http"
	"sort"
)

type ApiError struct {
	Message    string
	StatusCode int
}

func (e *ApiError) Error() string {
	return e.Message
}

type ChatOrchestrator struct {
	rooms map[string]*chatRoom
}

func NewBadRequestError(message string) *ApiError {
	return &ApiError{Message: message, StatusCode: 400}
}

func (ch *ChatOrchestrator) deleteRoom(room *chatRoom) {
	delete(ch.rooms, room.Name)
}

func NewChatOrchestrator() *ChatOrchestrator {
	return &ChatOrchestrator{rooms: map[string]*chatRoom{}}
}

func (ch *ChatOrchestrator) CreateChat(chatRoomName string) error {
	log.Printf("Creating room '%v'", chatRoomName)
	if ch.rooms[chatRoomName] != nil {
		return NewBadRequestError("chat name is already taken")
	}

	ch.rooms[chatRoomName] = newRoom(chatRoomName, ch)
	log.Printf("Rooms available: %v", len(ch.rooms))
	return nil
}

func (ch *ChatOrchestrator) GetChats() ([]string, error) {
	roomNames := make([]string, 0, len(ch.rooms))
	log.Println("Reading rooms!")
	for _, room := range ch.rooms {
		log.Printf("read room: %v", room.Name)
		roomNames = append(roomNames, room.Name)
	}
	sort.Strings(roomNames)
	return roomNames, nil
}

func (ch *ChatOrchestrator) JoinChat(chatRoomName string, clientName string, w http.ResponseWriter, r *http.Request) (*Client, error) {
	log.Printf("Joining room: '%v'", chatRoomName)
	chats, _ := ch.GetChats()
	log.Printf("current room 1: '%v'", chats[0])
	if ch.rooms[chatRoomName] != nil {
		return ch.rooms[chatRoomName].Connect(clientName, w, r)
	}
	return nil, errors.New("failed to join, chat room does not exist")
}

func (ch *ChatOrchestrator) GetChatInformation(chatName string) (*ChatInformation, error) {
	for _, room := range ch.rooms {
		if room.Name == chatName {
			return room.GetChatInformation(), nil
		}
	}
	return nil, errors.New("Chat does not exist")
}
