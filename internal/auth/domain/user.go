package domain

import (
	"context"
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string `json:"first_name"`
	LastName string  `json:"last_name"`
	StudentID string `json:"student_id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	JoinedYear string `json:"joined_year"`
	ProfilePicture *string `json:"profile_picture"`
	Gender string `json:"gender"`
	CreatedAt time.Time `json:"created_at"`
}


type AuthUseCase interface {
	Register(ctx context.Context, user *User) error
	LoginWithEmail(ctx context.Context, email string, password string) (string, error)
	LoginWithId(ctx context.Context, studentId string, password string) (string, error)
}
