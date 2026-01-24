package http

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Ramsi97/edu-social-backend/internal/auth/domain"
	"github.com/Ramsi97/edu-social-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type authHandler struct {
	usecase domain.AuthUseCase
}

func NewAuthHandler(rg *gin.RouterGroup, uc domain.AuthUseCase) {
	handler := &authHandler{
		usecase: uc,
	}

	rg.POST("/register", handler.Register)
	rg.POST("/login", handler.Login)
}
func (h *authHandler) Register(ctx *gin.Context) {
	var req *domain.RegisterRequest

	
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(req)

	err := h.usecase.Register(ctx, req)
	if err != nil {
		log.Fatal(err)
		response.Error(ctx, http.StatusInternalServerError, "Server Error", err.Error())
		return
	}

	response.Success(ctx, http.StatusCreated, "user regitered succefully", gin.H{})
}



func (h *authHandler) Login(ctx *gin.Context) {
	var req domain.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(req)
	if req.Email != nil && *req.Email != "" {
		fmt.Println("email: " + *req.Email)
		user, token, err := h.usecase.LoginWithEmail(ctx.Request.Context(), *req.Email, req.Password)
		if err != nil {
			response.Error(ctx, http.StatusUnauthorized, "invalid email or password", err.Error())
			return
		}
		userResponse := domain.UserResponse{
			ID:         user.ID.String(),
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			StudentID:  user.StudentID,
			Email:      user.Email,
			JoinedYear: user.JoinedYear,
			ProfilePicture: user.ProfilePicture,
			Gender:     user.Gender,
			CreatedAt:  user.CreatedAt.Format(time.RFC3339),
		}




		response.Success(ctx, http.StatusOK, "Login Successful", domain.LoginResponse{
			Token: token,
			User: userResponse,
		})
		return
	} else if req.StudentID != nil && *req.StudentID != "" {
		token, err := h.usecase.LoginWithId(ctx.Request.Context(), *req.StudentID, req.Password)
		if err != nil {
			response.Error(ctx, http.StatusUnauthorized, "invalid Id or password", err.Error())
			return
		}

		response.Success(ctx, http.StatusOK, "Login Successful", domain.LoginResponse{
			Token: token,
		})
		return
	}
	response.Error(ctx, http.StatusBadRequest, "Email or Student ID is required", "")

}