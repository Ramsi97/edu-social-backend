package usecase

import (
	"context"

	"github.com/Ramsi97/edu-social-backend/internal/like/domain"
	"github.com/Ramsi97/edu-social-backend/internal/like/repository/interfaces"
	"github.com/google/uuid"
)


type likeUseCase struct {
	repo interfaces.LikeRepository
}

func NewLikeUseCase(repo interfaces.LikeRepository) domain.LikeUseCase {
	return &likeUseCase{
		repo: repo,
	}
}

func (u *likeUseCase) ToggleUseCase(ctx context.Context, userID, postID uuid.UUID) (bool, error) {
	exists, err := u.repo.Exists(ctx, userID, postID)
    if err != nil {
        return false, err
    }

    if exists {
		if err := u.repo.Delete(ctx, userID, postID); err != nil {
			return true, err
		}
		return false, nil
	}

	if err := u.repo.Create(ctx, userID, postID); err != nil {
		return false, err
	}

	return true, nil
}