package domain

type User struct {
	ID        string `json:"id"`
	First_name string `json:"first_name"`
	Last_name string  `json:"last_name"`
	StudentId string `json:"studentId"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Joined_year string `json:"joined_year"`
	Profile_picture string `json:"profile_picture"`
	Gender string `json:"gender"`
	Created_at string `json:"created_at"`
}
