package usecase

import (
	"context"

	"github.com/Ramsi97/edu-social-backend/internal/chat/domain"
)

type chatUseCase struct {
	repo domain.ChatRepository
}

func NewChatUseCase(r domain.ChatRepository) domain.ChatUseCase {
	return &chatUseCase{repo: r}
}

func (u *chatUseCase) SendMessage(ctx context.Context, msg *domain.Message) error {
	if msg.Content == "" {
		return &domain.ChatError{Message: "message cannot be empty"}
	}
	return u.repo.SaveMessage(ctx, *msg)
}

func (u *chatUseCase) GetMessages(ctx context.Context, roomID string) ([]domain.Message, error) {
	return u.repo.GetChatHistory(ctx, roomID)
}
