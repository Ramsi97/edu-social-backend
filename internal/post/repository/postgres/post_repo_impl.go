package postgres

import (
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

func (r *postRepo) CreatePost(post *domain.Post) error {

	post.ID = uuid.New()
	post.CreatedAt = time.Now()

	query := `
        INSERT INTO posts (id, author_id, content, media_url, created_at)
        VALUES ($1, $2, $3, $4, $5)
    `
	result, err := r.db.Exec(query,
		post.ID,
		post.AuthorID,
		post.Content,
		post.MediaUrl,
		post.CreatedAt,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *postRepo) GetFeed(limit int, lastSeenTime *time.Time) ([]domain.Post, error) {
	var rows *sql.Rows
	var err error

	if lastSeenTime == nil {
		rows, err = r.db.Query(
			`SELECT * FROM posts
			ORDERED BY created_at DESC
			LIMIT $1`, limit,
		)
	} else {
		rows, err = r.db.Query(`
			SELECT *
			FROM posts
			WHERE created_at < $1
			ORDER BY created_at DESC
			LIMIT $2`, *lastSeenTime, limit,
		)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var posts []domain.Post

	for rows.Next() {
		var p domain.Post
		if err = rows.Scan(&p.ID, &p.AuthorID, &p.Content, &p.MediaUrl, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}
