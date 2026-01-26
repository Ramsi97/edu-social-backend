package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Ramsi97/edu-social-backend/internal/post/domain"
	sharedInterfaces "github.com/Ramsi97/edu-social-backend/internal/shared/interfaces"
	"github.com/Ramsi97/edu-social-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostHandler struct {
	usecase      domain.PostUseCase
	mediaStorage sharedInterfaces.MediaStorage
}

func NewPostHandler(
	rg *gin.RouterGroup,
	uc domain.PostUseCase,
	ms sharedInterfaces.MediaStorage,
) {
	handler := PostHandler{
		usecase:      uc,
		mediaStorage: ms,
	}

	rg.POST("", handler.CreatePost)
	rg.GET("/feed", handler.GetFeed)
}

func (p *PostHandler) CreatePost(ctx *gin.Context) {

	content := ctx.PostForm("content")

	file, err := ctx.FormFile("file")
	if err != nil && err != http.ErrMissingFile {
		response.Error(ctx, http.StatusBadRequest, "Invalid file", err.Error())
		return
	}

	fmt.Println(content)

	uuidString := ctx.GetString("user_id")
	authorID, err := uuid.Parse(uuidString)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	var mediaURL string

	if file != nil {
		mediaURL, err = p.mediaStorage.UploadToCloudinary(ctx, file)
		if err != nil {
			response.Error(ctx, http.StatusInternalServerError, "Failed to upload media", err.Error())
			return
		}
	}

	post := &domain.Post{
		Author: domain.UserSummary{ID: authorID},
		Content:  content,
		MediaUrl: mediaURL,
	}

	err = p.usecase.CreatePost(ctx, post)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Failed to create post", err.Error())
		return
	}

	response.Success(ctx, http.StatusCreated, "Post created successfully", nil)
}

func (p *PostHandler) GetFeed(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "20")
	lastSeenStr := ctx.Query("filter")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid limit", err.Error())
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

	uuidString := ctx.GetString("user_id")
	authorID, err := uuid.Parse(uuidString)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	posts, err := p.usecase.GetFeed(ctx, limit, lastSeenTime, authorID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to fetch feed", err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "Feed fetched successfully", posts)
}
