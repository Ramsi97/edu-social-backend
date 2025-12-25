package domain

import (
	"context"

	"github.com/google/uuid"
)

type Like struct {
	UserID uuid.UUID `json:"user_id"`
	PostID uuid.UUID `json:"post_id"`
}

type LikeUseCase interface{
	ToggleUseCase(ctx context.Context, userID, postID uuid.UUID) (bool, error)
}