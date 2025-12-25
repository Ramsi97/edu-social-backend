package interfaces

import (
	"context"

	"github.com/google/uuid"
)

type LikeRepository interface {
	Create(ctx context.Context, userID, postID uuid.UUID) error
	Delete(ctx context.Context, userID, postID uuid.UUID) error
	Exists(ctx context.Context, userID, postID uuid.UUID) (bool, error)
	// GetCountByPostID(ctx context.Context, postID string) (int, error)
}
