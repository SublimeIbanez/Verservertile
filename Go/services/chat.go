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
	PublicRooms  []ChatRoom
	PrivateRooms []ChatRoom
}

func HandleChatService() *ChatService {
	return &ChatService{
		PublicRooms:  []ChatRoom{},
		PrivateRooms: []ChatRoom{},
	}
}

type ChatRoom struct {
	Public   bool
	Users    map[string]string // uuid -> username
	Messages []Message
	Input    chan Message
	Output   chan []Message
	mtx      sync.RWMutex
	context  context.Context
	cancel   context.CancelFunc
}

type Message struct {
	User      string
	Timestamp time.Time
	Content   string
}

func CreateRoom(public bool) *ChatRoom {
	context, cancel := context.WithCancel(context.Background())

	room := ChatRoom{
		Public:   public,
		Users:    make(map[string]string),
		Messages: []Message{},
		Input:    make(chan Message, 100),
		Output:   make(chan []Message),
		context:  context,
		cancel:   cancel,
	}

	go room.handleInput()

	return &room
}

func (room *ChatRoom) handleInput() {
	for message := range room.Input {
		select {
		case <-room.context.Done():

		default:
			{
				// Intake the message andupdate the messages array
				room.mtx.Lock()
				if len(room.Messages) >= CHAT_MESSAGE_COUNT_MAX {
					remove := len(room.Messages) - CHAT_MESSAGE_COUNT_MAX + 1
					room.Messages = append(room.Messages[remove:], room.Messages[remove+1:]...)
				}
				room.Messages = append(room.Messages, message)
				room.mtx.Unlock()

				// Send the updated array to the output
				room.mtx.RLock()
				room.Output <- room.Messages
				room.mtx.RUnlock()
			}
		}
	}
}
