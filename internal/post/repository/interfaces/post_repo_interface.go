package interfaces

import (
	"time"

	"github.com/Ramsi97/edu-social-backend/internal/post/domain"
)

type PostRepository interface{
	CreatePost(post *domain.Post) error
	GetFeed(limit int, lastSeenTime *time.Time) ([]domain.Post, error)
}