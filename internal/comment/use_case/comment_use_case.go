package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Ramsi97/edu-social-backend/internal/comment/domain"
	"github.com/Ramsi97/edu-social-backend/internal/comment/repository/interfaces"
	"github.com/google/uuid"
)

type commentUseCase struct {
	repo interfaces.CommentRepository
}

func NewCommentUseCase(repo interfaces.CommentRepository) domain.CommentUseCase {
	return &commentUseCase{
		repo: repo,
	}
}

// Create implements domain.CommentUseCase.
func (c *commentUseCase) Create(ctx context.Context, userID string, postID string, content string) error {
	
	if content == ""{
		return errors.New("comment content cannot be empty")
	}

	uID, err := uuid.Parse(userID)

	if err != nil {
		return errors.New("invalid user id")
	}
	pID, err := uuid.Parse(postID)
	if err != nil {
		return errors.New("invalid user id")
	}

	comment := domain.Comment{
		ID: uuid.New(),
		UserID: uID,
		PostID: pID,
		Content: content,
		CreatedAT: time.Now(),
	}

	return c.repo.Create(ctx, &comment)
}

func (c *commentUseCase) Delete(ctx context.Context,userID, commentID string) error {
	
	uID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	cID, err := uuid.Parse(commentID)
	if err != nil {
		return errors.New("invalid comment id")
	}

	existingComment, err := c.repo.GetByID(ctx, cID)
	if err != nil {
		return err
	}

	if existingComment.UserID != uID {
		return errors.New("you are not authorized to delete this comment")
	}

	return c.repo.Delete(ctx, cID)
}

func (c *commentUseCase) GetByPostID(ctx context.Context, postID string) ([]domain.Comment, error) {
	pID, err := uuid.Parse(postID)
	if err != nil {
		return nil, errors.New("invalid post id")
	}

	comments, err := c.repo.GetByPostID(ctx, pID)
	if err != nil {
		return nil, err
	}

	return comments, nil
}
