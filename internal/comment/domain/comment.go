package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrCommentNotFound = errors.New("comment not found")

type Comment struct {
	ID uuid.UUID `json:"id"`
	Content string `json:"content"`
	UserID uuid.UUID `json:"user_id"`
	PostID uuid.UUID `json:"post_id"`
	CreatedAT time.Time `json:"created_at"`
}

type CommentRequest struct {
	Content string `json:"content"`
	UserID string `json:"user_id"`
	PostID string `json:"post_id"`
}

type CommentUseCase interface {
	Create(ctx context.Context, userID, postID, content string) error
	Delete(ctx context.Context, userID, commentID string) error
	GetByPostID(ctx context.Context, postID string) ([]Comment, error)
}