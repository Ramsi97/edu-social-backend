package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/Ramsi97/edu-social-backend/internal/chat/domain"
	"github.com/google/uuid"
)

type chatRepo struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) domain.ChatRepository {
	return &chatRepo{db: db}
}

// SaveMessage saves a chat message in DB
func (r *chatRepo) SaveMessage(ctx context.Context, msg domain.Message) error {
	if msg.ID == uuid.Nil {
		msg.ID = uuid.New()
	}
	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = time.Now()
	}

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO chat_messages (id, sender_id, room_id, content, created_at) 
		 VALUES ($1, $2, $3, $4, $5)`,
		msg.ID, msg.SenderID, msg.RoomID, msg.Content, msg.CreatedAt,
	)
	return err
}

// GetChatHistory retrieves messages for a room
func (r *chatRepo) GetChatHistory(ctx context.Context, roomID string) ([]domain.Message, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, sender_id, room_id, content, created_at 
		 FROM chat_messages 
		 WHERE room_id=$1
		 ORDER BY created_at ASC`, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []domain.Message{}
	for rows.Next() {
		var msg domain.Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.RoomID, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
