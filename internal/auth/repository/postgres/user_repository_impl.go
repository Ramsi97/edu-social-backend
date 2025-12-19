package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Ramsi97/edu-social-backend/internal/auth/domain"
	"github.com/Ramsi97/edu-social-backend/internal/auth/repository/interfaces"
	"github.com/google/uuid"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) interfaces.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	newID, err := uuid.NewV7()
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	user.ID = newID
	user.CreatedAt = now

	const layout = "2006-01-02"

	joinedYear, err := time.Parse(layout, user.JoinedYear)
	fmt.Println("year"+user.JoinedYear)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO users (
			id, first_name, last_name, student_id, email, 
			password_hash, joined_year, profile_picture, gender, created_at
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err = r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.StudentID,
		user.Email,
		user.Password,
		joinedYear,
		user.ProfilePicture,
		user.Gender,
		user.CreatedAt,
	)
	return err
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT 
			id, first_name, last_name, student_id,
			email, password_hash, joined_year,
			profile_picture, gender, created_at
		FROM users
		WHERE email = $1
	`
	var user domain.User

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.StudentID,
		&user.Email,
		&user.Password,
		&user.JoinedYear,
		&user.ProfilePicture,
		&user.Gender,
		&user.CreatedAt,
	)

	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByStudentId(ctx context.Context, studentID string) (*domain.User, error) {

	query := `
		SELECT
			id, first_name, last_name, student_id,
			email, password_hash, joined_year,
			profile_picture, gender, created_at
		FROM users
		WHERE student_id = $1
	`

	var user domain.User

	err := r.db.QueryRowContext(ctx, query, studentID).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.StudentID,
		&user.Email,
		&user.Password,
		&user.JoinedYear,
		&user.ProfilePicture,
		&user.Gender,
		&user.CreatedAt,
	)

	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
