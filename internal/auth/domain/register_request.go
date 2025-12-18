package domain

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName string  `json:"last_name"`
	StudentID string `json:"student_id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	JoinedYear string `json:"joined_year"`
	ProfilePicture *string `json:"profile_picture"`
	Gender string `json:"gender"`
}