package https

import (
	"fmt"
	"net/http"

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

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &domain.User{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		Password:       req.Password,
		StudentID:      req.StudentID,
		JoinedYear:     req.JoinedYear,
		ProfilePicture: req.ProfilePicture,
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
