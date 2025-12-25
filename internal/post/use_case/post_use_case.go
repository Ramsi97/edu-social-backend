package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Ramsi97/edu-social-backend/internal/post/domain"
	"github.com/Ramsi97/edu-social-backend/internal/post/repository/interfaces"
)

type postUseCase struct {
	repo interfaces.PostRepository
}

func NewPostUseCase(r interfaces.PostRepository) domain.PostUseCase {
	return &postUseCase{
		repo: r,
	}
}

func (u *postUseCase) GetFeed(ctx context.Context, limit int, lasttimeSeen *time.Time) ([]domain.Post, error) {
	if limit <= 0 {
		limit = 20
	}
	posts, err := u.repo.GetFeed(ctx, limit, lasttimeSeen)

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (u *postUseCase) CreatePost(ctx context.Context, post *domain.Post) error {
	if post.Content == "" && post.MediaUrl == "" {
		return errors.New("post cannot be empty")
	}

	return u.repo.CreatePost(ctx, post)
}

func (u *postUseCase) CountLikes(ctx context.Context, postID string) (int, error){
	panic("unimplemented")
}