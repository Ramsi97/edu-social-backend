package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `json:"id"`
	Author  UserSummary `json:"author"`
	Content   string    `json:"content"`
	MediaUrl  string    `json:"media_url"`
	LikeCount int       `json:"like_count"`
	CommentCount int `json:"comment_count"`
	CreatedAt time.Time `json:"created_at"`
	LikedByMe bool `json:"liked_by_me"`
}

type UserSummary struct {
    ID            uuid.UUID `json:"id"`
    FirstName     string    `json:"first_name"`
    LastName      string    `json:"last_name"`
    ProfilePicture string   `json:"profile_picture"`
    JoinedYear    time.Time       `json:"joined_year"`
}

type PostUseCase interface {
	CreatePost(ctx context.Context, post *Post) error
	GetFeed(ctx context.Context, limit int, lastSeenTime *time.Time, authorID uuid.UUID) ([]Post, error)
}
