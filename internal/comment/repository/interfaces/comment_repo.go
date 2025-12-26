package interfaces

import (
	"context"

	"github.com/Ramsi97/edu-social-backend/internal/comment/domain"
	"github.com/google/uuid"
)

type CommentRepository interface {
	Create(ctx context.Context,comment *domain.Comment) error
	Delete(ctx context.Context, commentID uuid.UUID) error
	GetByPostID(ctx context.Context, postID uuid.UUID) ([]domain.Comment, error)
	GetByID(ctx context.Context, commentID uuid.UUID) (domain.Comment, error)
}