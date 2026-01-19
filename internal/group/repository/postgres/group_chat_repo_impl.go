package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Ramsi97/edu-social-backend/internal/group/domain"
	"github.com/Ramsi97/edu-social-backend/internal/group/repository/interfaces"
	"github.com/google/uuid"
)

type groupChatRepo struct {
	db *sql.DB
}

func NewGroupChatRepo(db *sql.DB) interfaces.GroupChatRepo {
	return &groupChatRepo{
		db: db,
	}
}

func (r *groupChatRepo) CreateGroup(ctx context.Context, group *domain.Group) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO groups (id, name, created_at) VALUES ($1, $2, NOW())
	`, group.ID, group.Name)
	return err
}

// GetGroup retrieves group ID by name
func (r *groupChatRepo) GetGroup(ctx context.Context, groupName string) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, `
		SELECT id FROM groups WHERE name = $1
	`, groupName).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("group not found: %w", err)
	}
	return id, nil
}

// JoinGroup adds a user to a group
func (r *groupChatRepo) JoinGroup(ctx context.Context, groupID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO group_members (group_id, user_id, joined_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT DO NOTHING
	`, groupID, userID)
	return err
}

// LeaveGroup removes a user from a group
func (r *groupChatRepo) LeaveGroup(ctx context.Context, groupID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM group_members WHERE group_id=$1 AND user_id=$2
	`, groupID, userID)
	return err
}

// SaveMessage inserts a new message
func (r *groupChatRepo) SaveMessage(ctx context.Context, msg *domain.Message) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO messages (id, group_id, sender_id, content, created_at)
		VALUES ($1, $2, $3, $4, NOW())
	`, msg.ID, msg.GroupID, msg.SenderID, msg.Content)
	return err
}

// IsMember checks if user is a member of a group
func (r *groupChatRepo) IsMember(ctx context.Context, userID, groupID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM group_members WHERE user_id=$1 AND group_id=$2
		)
	`, userID, groupID).Scan(&exists)
	return exists, err
}

// GetMessages retrieves the latest messages for a group
func (r *groupChatRepo) GetMessages(ctx context.Context, groupID uuid.UUID, limit int) ([]*domain.Message, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, group_id, sender_id, content, created_at
		FROM messages
		WHERE group_id=$1
		ORDER BY created_at DESC
		LIMIT $2
	`, groupID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []*domain.Message{}
	for rows.Next() {
		var msg domain.Message
		if err := rows.Scan(&msg.ID, &msg.GroupID, &msg.SenderID, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	return messages, nil
}
