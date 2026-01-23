package interfaces

import (
	"context"

	"github.com/Ramsi97/edu-social-backend/internal/group/domain"
	"github.com/google/uuid"
)

type GroupChatRepo interface {
	CreateGroup(ctx context.Context, group *domain.Group) error
	GetGroup(ctx context.Context, groupName string) (uuid.UUID, error)
	JoinGroup(ctx context.Context, groupID, userID uuid.UUID) error
	LeaveGroup(ctx context.Context, groupID, userID uuid.UUID) error
	SaveMessage(ctx context.Context, msg *domain.Message) error
	IsMember(ctx context.Context, userID, groupID uuid.UUID) (bool, error)
	GetGroupsForUser(ctx context.Context, userID uuid.UUID) ([]*domain.Group, error)
	GetMessages(ctx context.Context, groupID uuid.UUID, limit int) ([]*domain.Message, error)
}