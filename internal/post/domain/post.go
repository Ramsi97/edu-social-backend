package domain

import (
	"time"

	"github.com/google/uuid"
)


type Post struct{
	ID uuid.UUID `json:"id"`
	AuthorID uuid.UUID `json:"authod_id"`
	Content string `json:"content"`
	MediaUrl string `json:"media_url"`
	CreatedAt string `json:"created_at"`
}

type PostUseCase interface {
	CreatePost(post *Post) error
	GetFeed(limit int, lastSeenTime *time.Time) ([]Post, error)
}