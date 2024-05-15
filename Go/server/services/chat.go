package services

import (
	"context"
	"sync"
	"time"
)

const (
	// Maximum number of chat messages retained in memory
	CHAT_MESSAGE_COUNT_MAX = 1_000
	// Maximum number of characters in a chat message
	CHAT_MESSAGE_SIZE_MAX = 500
	// Maximum number of public rooms available
	PUBLIC_ROOM_MAX = 1000
	// Maximum number of private rooms available
	PRIVATE_ROOM_MAX = 1000
)

type ChatService struct {
	PublicRooms  []Room
	PrivateRooms []Room
}

func HandleChatService() *ChatService {
	return &ChatService{
		PublicRooms:  []Room{},
		PrivateRooms: []Room{},
	}
}

type Room struct {
	Public   bool
	Users    map[string]string // uuid -> username
	Messages []Message
	Input    chan Message
	Output   chan Message
	mtx      sync.RWMutex
	context  context.Context
	cancel   context.CancelFunc
}

type Message struct {
	User      string
	Timestamp time.Time
	Content   string
}

func CreateRoom(public bool) *Room {
	context, cancel := context.WithCancel(context.Background())

	room := Room{
		Public:   public,
		Users:    make(map[string]string),
		Messages: []Message{},
		Input:    make(chan Message, 100),
		Output:   make(chan Message, 100),
		context:  context,
		cancel:   cancel,
	}

	go room.handleInput()

	return &room
}

func (room *Room) handleInput() {
	for message := range room.Input {
		select {
		case <-room.context.Done():

		default:
			{
				room.mtx.Lock()
				if len(room.Messages) >= CHAT_MESSAGE_COUNT_MAX {
					remove := len(room.Messages) - CHAT_MESSAGE_COUNT_MAX + 1
					room.Messages = append(room.Messages[remove:], room.Messages[remove+1:]...)
				}
				room.Messages = append(room.Messages, message)
				room.mtx.Unlock()
			}
		}
	}
}
