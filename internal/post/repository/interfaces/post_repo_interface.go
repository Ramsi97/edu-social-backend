package interfaces

import (
	"context"
	"time"

	"github.com/Ramsi97/edu-social-backend/internal/post/domain"
)

type PostRepository interface{
	CreatePost(ctx context.Context, post *domain.Post) error
	GetFeed(ctx context.Context, limit int, lastSeenTime *time.Time) ([]domain.Post, error)
}