package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID uuid.UUID `json:"id"`
	SenderID uuid.UUID `json:"sender_id"`
	RookID uuid.UUID `json:"room_id"`
	Content string `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatRepository interface {
	GetChatHistory(ctx context.Context, roomID string) ([]Message, error)
	SaveMessage(ctx context.Context, msg Message) error	
}