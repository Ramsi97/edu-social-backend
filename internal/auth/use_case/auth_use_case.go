package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Ramsi97/edu-social-backend/internal/auth/domain"
	repoInterface "github.com/Ramsi97/edu-social-backend/internal/auth/repository/interfaces"
	cldInterface "github.com/Ramsi97/edu-social-backend/internal/shared/interfaces"
	"github.com/Ramsi97/edu-social-backend/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type authUseCase struct {
	userRepo repoInterface.UserRepository
	cld cldInterface.MediaStorage
}

func NewAuthUseCase(repo repoInterface.UserRepository, cld cldInterface.MediaStorage) domain.AuthUseCase {
	return &authUseCase{
		userRepo: repo,
		cld: cld,
	}
}

func (a *authUseCase) LoginWithEmail(ctx context.Context, email string, password string) (string, error) {
	user, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil{
		return "", errors.New("invalid credentials")
	}

	token, err := auth.GenerateToken(user.ID, time.Hour*144)
	if(err != nil){
		return "", errors.New("please, try again")
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

	token, err := auth.GenerateToken(user.ID, time.Hour*144)
	if(err != nil){
		return "", errors.New("please, try again")
	}
	return token, nil
}


func (a *authUseCase) Register(ctx context.Context, req *domain.RegisterRequest) error{
	existingUser, _ := a.userRepo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return errors.New("user already exists")
	}

	existingUser, _ = a.userRepo.FindByStudentId(ctx, req.StudentID)
	if existingUser != nil {
		return errors.New("user already exists")
	}

	hashedByte, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil{
		return err
	}
	
	var profileURL *string

	if req.ProfilePictureFile != nil {
		url, err := a.cld.UploadToCloudinary(ctx, req.ProfilePictureFile)
		if err != nil {
			return err
		}
		profileURL = &url
	}


	user := &domain.User{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		Password:       req.Password,
		StudentID:      req.StudentID,
		JoinedYear:     req.JoinedYear,
		ProfilePicture: profileURL, 
		Gender:         req.Gender,
	}

	user.Password = string(hashedByte)
	return a.userRepo.Create(ctx, user)
}

