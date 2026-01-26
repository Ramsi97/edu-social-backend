package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/Ramsi97/edu-social-backend/internal/post/domain"
	"github.com/Ramsi97/edu-social-backend/internal/post/repository/interfaces"
	"github.com/google/uuid"
)

type postRepo struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) interfaces.PostRepository {
	return &postRepo{
		db: db,
	}
}

func (r *postRepo) CreatePost(ctx context.Context, post *domain.Post) error {
	post.ID = uuid.New()
	post.CreatedAt = time.Now()

	query := `
		INSERT INTO posts (id, author_id, content, media_url, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		post.ID,
		post.Author.ID,
		post.Content,
		post.MediaUrl,
		post.CreatedAt,
	)

	return err
}

func (r *postRepo) GetFeed(
	ctx context.Context,
	limit int,
	lastSeenTime *time.Time,
	currentUserID uuid.UUID, // current logged-in user
) ([]domain.Post, error) {
	var rows *sql.Rows
	var err error

	// Base query with like count, liked_by_me, and comment count
	baseQuery := `
        SELECT 
            p.id,
            p.content,
            p.media_url,
            p.created_at,

            u.id AS author_id,
            u.first_name,
            u.last_name,
            u.profile_picture,
            u.joined_year,

            COUNT(DISTINCT pl.post_id) AS like_count,
            CASE WHEN ul.user_id IS NOT NULL THEN TRUE ELSE FALSE END AS liked_by_me,
            COUNT(DISTINCT c.id) AS comment_count
        FROM posts p
        JOIN users u ON p.author_id = u.id
        LEFT JOIN posts_likes pl ON p.id = pl.post_id
        LEFT JOIN posts_likes ul ON p.id = ul.post_id AND ul.user_id = $2
        LEFT JOIN comments c ON p.id = c.post_id
    `

	if lastSeenTime == nil {
		query := baseQuery + `
            GROUP BY p.id, u.id, ul.user_id
            ORDER BY p.created_at DESC
            LIMIT $1
        `
		rows, err = r.db.QueryContext(ctx, query, limit, currentUserID)
	} else {
		query := baseQuery + `
            WHERE p.created_at < $3
            GROUP BY p.id, u.id, ul.user_id
            ORDER BY p.created_at DESC
            LIMIT $1
        `
		rows, err = r.db.QueryContext(ctx, query, limit, currentUserID, *lastSeenTime)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []domain.Post{}

	for rows.Next() {
		var p domain.Post
		var author domain.UserSummary

		if err := rows.Scan(
			&p.ID,
			&p.Content,
			&p.MediaUrl,
			&p.CreatedAt,
			&author.ID,
			&author.FirstName,
			&author.LastName,
			&author.ProfilePicture,
			&author.JoinedYear,
			&p.LikeCount,
			&p.LikedByMe,
			&p.CommentCount,
		); err != nil {
			return nil, err
		}

		p.Author = author
		posts = append(posts, p)
	}

	return posts, nil
}
