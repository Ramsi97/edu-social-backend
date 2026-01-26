package http

import (
	"net/http"

	"github.com/Ramsi97/edu-social-backend/internal/like/domain"
	"github.com/Ramsi97/edu-social-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type likeHandler struct {
	useCase domain.LikeUseCase
}

func NewLikeHandler(rg *gin.RouterGroup, uc domain.LikeUseCase) {
	handler := &likeHandler{
		useCase: uc,
	}

	rg.POST("/togglelike", handler.Togglelike)
}

func (l *likeHandler) Togglelike(ctx *gin.Context) {

	var req domain.LikeRequest
	
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}

	userIDStr := ctx.GetString("user_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, "Invalid user ID", "")
		return
	}

	postID, err := uuid.Parse(req.PostID)

	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid post ID", "")
		return
	}
	liked, err := l.useCase.ToggleUseCase(ctx, userID, postID)

	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Internal server Error", err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Successfully toggled", gin.H{
		"liked": liked,
	})

}
