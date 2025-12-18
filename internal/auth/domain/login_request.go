package domain

type LoginRequest struct {
	Email *string `json:"email"`
	StudentID *string `json:"student_id"`
	Password string `json:"password"`
}