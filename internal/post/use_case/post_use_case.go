package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Ramsi97/edu-social-backend/internal/post/domain"
	"github.com/Ramsi97/edu-social-backend/internal/post/repository/interfaces"
	"github.com/google/uuid"
)

type postUseCase struct {
	repo interfaces.PostRepository
}

func NewPostUseCase(r interfaces.PostRepository) domain.PostUseCase {
	return &postUseCase{
		repo: r,
	}
}

func (u *postUseCase) GetFeed(
	ctx context.Context,
	limit int,
	lastSeenTime *time.Time,
	authorID uuid.UUID,
) ([]domain.Post, error) {

	if limit <= 0 {
		limit = 20
	}

	return u.repo.GetFeed(ctx, limit, lastSeenTime, authorID)
}

func (u *postUseCase) CreatePost(ctx context.Context, post *domain.Post) error {
	if post.Content == "" && post.MediaUrl == "" {
		return errors.New("post cannot be empty")
	}

	return u.repo.CreatePost(ctx, post)
}
