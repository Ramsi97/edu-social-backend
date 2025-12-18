package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Ramsi97/edu-social-backend/internal/auth/domain"
	"github.com/Ramsi97/edu-social-backend/internal/auth/repository/interfaces"
	"github.com/Ramsi97/edu-social-backend/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type authUseCase struct {
	userRepo interfaces.UserRepository
}

func NewAuthUseCase(repo interfaces.UserRepository) interfaces.AuthUseCase {
	return &authUseCase{userRepo: repo}
}

func (a *authUseCase) LoginWithEmail(ctx context.Context, email string, password string) (string, error) {
	user, err := a.userRepo.FindByStudentId(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil{
		return "", errors.New("invalid credentials")
	}

	token, err := auth.GenerateToken(user.ID, time.Hour*144)
	if(err != nil){
		return "", errors.New("pleas, try again")
	}
	return token, nil
}

func (a *authUseCase) LoginWithId(ctx context.Context, studentId string, password string) (string, error) {
	user, err := a.userRepo.FindByStudentId(ctx, studentId)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil{
		return "", errors.New("invalid credentials")
	}

	return "mock_jwt_token_123", nil
}


func (a *authUseCase) Register(ctx context.Context, user *domain.User) error{
	existingUser, _ := a.userRepo.FindByEmail(ctx, user.Email)
	if existingUser != nil {
		return errors.New("user already exists")
	}

	existingUser, _ = a.userRepo.FindByStudentId(ctx, user.StudentId)
	if existingUser != nil {
		return errors.New("user already exists")
	}

	hashedByte, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil{
		return err
	}
	user.Password = string(hashedByte)
	return a.userRepo.Create(ctx, user)
}

