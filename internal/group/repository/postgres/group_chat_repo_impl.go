package postgres

import (
	"context"
	"database/sql"
	"errors"
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


func (r *groupChatRepo) CreateGroup(
	ctx context.Context,
	group *domain.Group,
) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Insert group
	res, err := tx.ExecContext(ctx, `
		INSERT INTO groups (id, name, owner_id, created_at)
		VALUES ($1, $2, $3, $4)
	`,
		group.ID,
		group.Name,
		group.OwnerID,
		group.CreatedAt,
	)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domain.ErrGroupAlreadyExists
	}

	// 2. Insert owner as member
	_, err = tx.ExecContext(ctx, `
		INSERT INTO group_members (group_id, user_id, role)
		VALUES ($1, $2, $3, $4)
	`,
		group.ID,
		group.OwnerID,
		domain.GroupRoleOwner,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *groupChatRepo) GetGroup(ctx context.Context, groupName string) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, `
		SELECT id FROM groups WHERE name = $1
	`, groupName).Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, domain.ErrGroupNotFound
		}
		return uuid.Nil, fmt.Errorf("failed to query group: %w", err)
	}

	return id, nil
}


// JoinGroup adds a user to a group
func (r *groupChatRepo) JoinGroup(ctx context.Context, groupID, userID uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO group_members (group_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, groupID, userID)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domain.ErrAlreadyMember
	}
	return nil
}

// LeaveGroup removes a user from a group
func (r *groupChatRepo) LeaveGroup(ctx context.Context, groupID, userID uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `
		DELETE FROM group_members WHERE group_id=$1 AND user_id=$2
	`, groupID, userID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domain.ErrNotMember
	}
	return nil
}

// SaveMessage inserts a new message
func (r *groupChatRepo) SaveMessage(ctx context.Context, msg *domain.Message) error {
	// optional: validate content before inserting
	if msg.Content == "" && msg.MediaURL == "" {
		return errors.New("either content or media_url must be provided")
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO group_msgs (id, group_id, author_id, content, media_url, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`,
		msg.ID,
		msg.GroupID,
		msg.AuthorID,
		msg.Content,
		msg.MediaURL,
	)
	return err
}

// IsMember checks if user is a member of a group
func (r *groupChatRepo) IsMember(ctx context.Context, userID, groupID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM group_members WHERE user_id=$1 AND room_id=$2
		)
	`, userID, groupID).Scan(&exists)
	return exists, err
}

func (r *groupChatRepo) GetMessages(ctx context.Context, groupID uuid.UUID, limit int) ([]*domain.Message, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, group_id, author_id, content, media_url, created_at
		FROM group_posts
		WHERE group_id=$1
		ORDER BY created_at DESC
		LIMIT $2
	`, groupID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*domain.Message{}
	for rows.Next() {
		var p domain.Message
		if err := rows.Scan(&p.ID, &p.GroupID, &p.AuthorID, &p.Content, &p.MediaURL, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}
	return posts, nil
}
