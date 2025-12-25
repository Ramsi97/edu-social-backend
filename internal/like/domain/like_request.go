package domain


type LikeRequest struct {
	UserID string `json:"user_id"`
	PostID string `json:"post_id"`
}