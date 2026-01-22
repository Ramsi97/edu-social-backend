package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `json:"id"`
	AuthorID  uuid.UUID `json:"author_id"`
	Content   string    `json:"content"`
	MediaUrl  string    `json:"media_url"`
	LikeCount int       `json:"like_count"`
	CreatedAt time.Time `json:"created_at"`
}

type PostUseCase interface {
	CreatePost(ctx context.Context, post *Post) error
	GetFeed(ctx context.Context, limit int, lastSeenTime *time.Time) ([]Post, error)
}
