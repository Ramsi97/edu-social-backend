package interfaces

import (
	"context"
	"github.com/Ramsi97/edu-social-backend/internal/auth/domain"
)
type AuthUseCase interface {
	Register(ctx context.Context, user *domain.User) error
	LoginWithEmail(ctx context.Context, email, password string) (string, error)
	LoginWithId(ctx context.Context, studentId, password string) (string, error)
}