package domain

import "context"

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName string  `json:"last_name"`
	StudentID string `json:"student_id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	JoinedYear string `json:"joined_year"`
	ProfilePicture *string `json:"profile_picture"`
	Gender string `json:"gender"`
	CreatedAt string `json:"created_at"`
}


type AuthUseCase interface {
	Register(ctx context.Context, user *User) error
	LoginWithEmail(ctx context.Context, email string, password string) (string, error)
	LoginWithId(ctx context.Context, studentId string, password string) (string, error)
}