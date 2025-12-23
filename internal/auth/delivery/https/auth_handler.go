package https

import (
	"fmt"
	"net/http"
	"context"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/Ramsi97/edu-social-backend/internal/auth/domain"
	"github.com/Ramsi97/edu-social-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	usecase domain.AuthUseCase
}

func NewAuthHandler(rg *gin.RouterGroup, uc domain.AuthUseCase) {
	handler := &AuthHandler{
		usecase: uc,
	}

	rg.POST("/register", handler.Register)
	rg.POST("/login", handler.Login)
}
func (h *AuthHandler) Register(ctx *gin.Context) {
	var req domain.RegisterRequest

	
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var profileURL *string

	if req.ProfilePictureFile != nil {
		url, err := UploadToCloudinary(ctx, req.ProfilePictureFile)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload profile picture"})
			return
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

	err := h.usecase.Register(ctx.Request.Context(), user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})
}



func (h *AuthHandler) Login(ctx *gin.Context) {
	var req domain.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Email != nil && *req.Email != "" {
		fmt.Println("email: " + *req.Email)
		token, err := h.usecase.LoginWithEmail(ctx.Request.Context(), *req.Email, req.Password)
		if err != nil {
			response.Error(ctx, http.StatusUnauthorized, "invalid email or password", err.Error())
			return
		}

		response.Success(ctx, http.StatusOK, "Login Successful", domain.LoginResponse{
			Token: token,
		})
		return
	} else if req.StudentID != nil && *req.StudentID != "" {
		token, err := h.usecase.LoginWithId(ctx.Request.Context(), *req.StudentID, req.Password)
		if err != nil {
			response.Error(ctx, http.StatusUnauthorized, "invalid Id ot password", err.Error())
			return
		}

		response.Success(ctx, http.StatusOK, "Login Successful", domain.LoginResponse{
			Token: token,
		})
		return
	}
	response.Error(ctx, http.StatusBadRequest, "Email or Student ID is required", "")

}
func UploadToCloudinary(ctx context.Context, file *multipart.FileHeader) (string, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return "", err
	}

	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	uploadResult, err := cld.Upload.Upload(ctx, f, uploader.UploadParams{
		Folder: "edu_social/profile_pics",
	})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}