package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	GroupID uuid.UUID `json:"group_id"`
	SenderID uuid.UUID `json:"sender_id"`
	Content string `json:"content"`
}

type Group struct {
	ID uuid.UUID `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	OwnerID uuid.UUID `json:"owner_id"`
	Members map[uuid.UUID]bool `json:"members"`
	CreatedAt time.Time `json:"created_at"`
}

type GroupChatUseCase interface {
    CreateGroup(ctx context.Context, ownerID uuid.UUID, groupName string) (uuid.UUID, error)
    JoinGroup(ctx context.Context, groupName string, userID uuid.UUID) error
    LeaveGroup(ctx context.Context, groupName string, userID uuid.UUID) error
    SendMessage(ctx context.Context, msg *Message) error
    GetMessages(ctx context.Context, groupID uuid.UUID, limit int) ([]*Message, error)
}


type MessagePublisher interface {
    Broadcast(msg *Message) error
}
