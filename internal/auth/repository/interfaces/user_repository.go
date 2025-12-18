package interfaces

import (
	"context"
	"github.com/Ramsi97/edu-social-backend/internal/auth/domain"
)


type UserRepository interface {
	Create (ctx context.Context, user domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByStudentId(ctx context.Context, studentId string)(*domain.User, error)
}