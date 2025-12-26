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
		post.AuthorID,
		post.Content,
		post.MediaUrl,
		post.CreatedAt,
	)

	return err
}

func (r *postRepo) GetFeed(ctx context.Context, limit int, lastSeenTime *time.Time) ([]domain.Post, error) {
	var rows *sql.Rows
	var err error

	query1 := `
		SELECT 
			p.id,
			p.author_id,
			p.content,
			p.media_url,
			p.created_at,
			COUNT(pl.post_id) AS like_count
		FROM posts p
		LEFT JOIN post_likes pl ON p.id = pl.post_id
		GROUP BY p.id, p.author_id, p.content, p.media_url, p.created_at
		ORDER BY p.created_at DESC
		LIMIT $1;
	`

	query2 := `
		SELECT 
			p.id,
			p.author_id,
			p.content,
			p.media_url,
			p.created_at,
			COUNT(pl.post_id) AS like_count
		FROM posts p
		LEFT JOIN post_likes pl ON p.id = pl.post_id
		WHERE p.created_at < $1
		GROUP BY p.id, p.author_id, p.content, p.media_url, p.created_at
		ORDER BY p.created_at DESC
		LIMIT $2;
	`

	if lastSeenTime == nil {
		rows, err = r.db.QueryContext(ctx, query1, limit)
	} else {
		rows, err = r.db.QueryContext(ctx, query2, *lastSeenTime, limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []domain.Post

	for rows.Next() {
		var p domain.Post
		if err := rows.Scan(
			&p.ID,
			&p.AuthorID,
			&p.Content,
			&p.MediaUrl,
			&p.CreatedAt,
			&p.LikeCount,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}
