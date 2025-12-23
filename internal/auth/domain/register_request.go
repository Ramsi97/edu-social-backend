package domain

import "mime/multipart"

type RegisterRequest struct {
	FirstName string                 `form:"first_name"`
	LastName  string                 `form:"last_name"`
	StudentID string                 `form:"student_id"`
	Email     string                 `form:"email"`
	Password  string                 `form:"password"`
	JoinedYear string                `form:"joined_year"`
	Gender    string                 `form:"gender"`
	ProfilePictureFile *multipart.FileHeader `form:"profile_picture"` 
}
