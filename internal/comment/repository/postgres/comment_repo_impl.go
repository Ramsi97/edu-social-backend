package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Ramsi97/edu-social-backend/internal/comment/domain"
	"github.com/Ramsi97/edu-social-backend/internal/comment/repository/interfaces"
	"github.com/google/uuid"
)

type commentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) interfaces.CommentRepository {
	return &commentRepository{
		db: db,
	}
}

func (c *commentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	query := `
		INSERT INTO comments (id, content, user_id, post_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := c.db.ExecContext(ctx, query,
		comment.ID,
		comment.Content,
		comment.User.UserID,
		comment.PostID,
		comment.CreatedAT,
	)
	return err

}

func (c *commentRepository) Delete(ctx context.Context, commentID uuid.UUID) error {
	
	query := `DELETE FROM comments WHERE id = $1`
	
	res, err := c.db.ExecContext(ctx, query, commentID)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domain.ErrCommentNotFound
	}

	return nil
}

func (c *commentRepository) GetByID(ctx context.Context, commentID uuid.UUID) (domain.Comment, error) {
	query := `
		SELECT id, content, user_id, post_id, created_at 
		FROM comments 
		WHERE id = $1
	`
	var comment domain.Comment
	err := c.db.QueryRowContext(ctx, query, commentID).Scan(
		&comment.ID,
		&comment.Content,
		&comment.User.UserID,
		&comment.PostID,
		&comment.CreatedAT,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Comment{}, domain.ErrCommentNotFound
		}
		return domain.Comment{}, err
	}

	return comment, nil
}
func (c *commentRepository) GetByPostID(ctx context.Context, postID uuid.UUID) ([]domain.Comment, error) {
	query := `
		SELECT 
			c.id,
			c.content,
			c.post_id,
			c.created_at,
			u.id AS user_id,
			u.first_name || ' ' || u.last_name AS user_name,
			u.profile_picture
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = $1
		ORDER BY c.created_at ASC
	`

	rows, err := c.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		var user domain.User

		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.PostID,
			&comment.CreatedAT,
			&user.UserID,
			&user.Name,
			&user.ProfilePicture,
		)
		if err != nil {
			return nil, err
		}

		comment.User = user
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
