package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID uuid.UUID `json:"id"`
	GroupID uuid.UUID `json:"group_id"`
	AuthorID uuid.UUID `json:"sender_id"`
	Content string `json:"content"`
	MediaURL string `json:"media_url"`
	CreatedAt time.Time `json:"created_at"`
}

type Group struct {
	ID uuid.UUID `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	OwnerID uuid.UUID `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
}

type GroupMember struct {
	GroupID  uuid.UUID
	UserID   uuid.UUID
	Role     string
	JoinedAt time.Time
}

const (
	GroupRoleOwner = "owner"
	GroupRoleAdmin = "admin"
	GroupRoleUser  = "member"
)


var (
	ErrGroupAlreadyExists = errors.New("group already exists")
	ErrGroupNotFound = errors.New("group not found")
	ErrNotMember     = errors.New("user is not a member of the group")
	ErrAlreadyMember = errors.New("user is already a member")
)
type GroupChatUseCase interface {
    CreateGroup(ctx context.Context, ownerID uuid.UUID, groupName string) (uuid.UUID, error)
    JoinGroup(ctx context.Context, groupName string, userID uuid.UUID) error
    LeaveGroup(ctx context.Context, groupName string, userID uuid.UUID) error
    SendMessage(ctx context.Context, msg *Message) error
    GetMessages(ctx context.Context, groupID uuid.UUID, limit int) ([]*Message, error)
}

