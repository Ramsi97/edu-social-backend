package postgres

import (
	"database/sql"
	"time"

	"github.com/Ramsi97/edu-social-backend/internal/post/repository/interfaces"
	"github.com/Ramsi97/edu-social-backend/internal/post/domain"
)


type postRepo struct{
	db *sql.DB
}

func NewPostRepository(db *sql.DB) interfaces.PostRepository{
	return &postRepo{
		db: db,
	}
}

func (r *postRepo) CreatePost(post *domain.Post) error {
	return nil
}

func (r *postRepo) GetFeed(limit int, lastSeenTime *time.Time) ([]domain.Post, error){
	var rows *sql.Rows
	var err error

	if lastSeenTime == nil {
		rows, err = r.db.Query(
			`SELECT * FROM posts
			ORDERED BY created_at DESC
			LIMIT $1`, limit, 
		)
	}else{
		rows, err = r.db.Query(`
			SELECT id, author_id, content, created_at
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

	for rows.Next(){
		var p domain.Post
		if err = rows.Scan(&p.ID, &p.AuthorID, &p.Content, &p.MediaUrl, &p.CreatedAt); err != nil{
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}
