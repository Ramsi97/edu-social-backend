package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Ramsi97/edu-social-backend/internal/post/domain"
	"github.com/Ramsi97/edu-social-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostHandler struct {
	usecase domain.PostUseCase
}

func NewPostHandler(rg *gin.RouterGroup, uc domain.PostUseCase) {
	handler := PostHandler{
		usecase: uc,
	}

	rg.POST("/createpost", handler.CreatePost)
	rg.GET("/getfeed", handler.GetFeed)
}

func (p *PostHandler) CreatePost(ctx *gin.Context) {
	var req domain.Post

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uuidString := ctx.GetString("user_id")
	authorID, err := uuid.Parse(uuidString)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post := &domain.Post{
		ID:       req.ID,
		AuthorID: authorID,
		Content:  req.Content,
		MediaUrl: req.MediaUrl,
	}

	err = p.usecase.CreatePost(post)

	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	response.Success(ctx, http.StatusCreated, "Successfully Posted", nil)
}

func (p *PostHandler) GetFeed(ctx *gin.Context) {
	limitstr := ctx.DefaultQuery("limit", "20")
	lastSeenStr := ctx.DefaultQuery("lastSeenTime", time.Now().Format(time.RFC3339))

	limit, err := strconv.Atoi(limitstr)

	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid limit ", err.Error())
		return
	}

	var lastSeenTime *time.Time

	if lastSeenStr != "" {
		parsedTime, err := time.Parse(time.RFC3339, lastSeenStr)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "Invalid lastSeenTime", err.Error())
			return
		}
		lastSeenTime = &parsedTime
	}

	posts, err := p.usecase.GetFeed(limit, lastSeenTime)

	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to get feed", err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Feed fetched", posts)

}
