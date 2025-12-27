package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Message represents a chat message
type Message struct {
	ID        uuid.UUID `json:"id"`
	SenderID  uuid.UUID `json:"sender_id"`
	RoomID    uuid.UUID `json:"room_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// ChatRepository defines repository actions
type ChatRepository interface {
	GetChatHistory(ctx context.Context, roomID string) ([]Message, error)
	SaveMessage(ctx context.Context, msg Message) error
}

// ChatUseCase defines the business logic layer
type ChatUseCase interface {
	SendMessage(ctx context.Context, msg *Message) error
	GetMessages(ctx context.Context, roomID string) ([]Message, error)
}

// ChatError is a custom error for chat validation
type ChatError struct {
	Message string
}

func (e *ChatError) Error() string {
	return e.Message
}
