package postgres

import (
	"context"
	"database/sql"

	"github.com/Ramsi97/edu-social-backend/internal/like/repository/interfaces"
	"github.com/google/uuid"
)

type likeRepository struct {
	db *sql.DB
}

func NewLikeRepository(db *sql.DB) interfaces.LikeRepository {
	return &likeRepository{
		db: db,
	}
}

func (l *likeRepository) Create(ctx context.Context, userID, postID uuid.UUID) error {

	query := `
		INSERT INTO posts_likes (user_id, post_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`
	_, err := l.db.ExecContext(ctx, query, userID, postID)
	return err
}

func (l *likeRepository) Delete(ctx context.Context, userID, postID uuid.UUID) error {

	query := `
		DELETE FROM posts_likes
		WHERE user_id = $1 AND post_id = $2	
	`

	_, err := l.db.ExecContext(ctx, query, userID, postID)
	return err
}

func (l *likeRepository) Exists(ctx context.Context, userID, postID uuid.UUID) (bool, error) {

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM posts_likes
			WHERE user_id = $1 AND post_id = $2
		)
	`

	var exists bool
	err := l.db.QueryRowContext(ctx, query, userID, postID).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}
