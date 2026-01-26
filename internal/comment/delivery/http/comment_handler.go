package http

import (
	"net/http"

	"github.com/Ramsi97/edu-social-backend/internal/comment/domain"
	"github.com/Ramsi97/edu-social-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type commentHandler struct {
	usecase domain.CommentUseCase
}

func NewCommentHandler(rg *gin.RouterGroup, uc domain.CommentUseCase) {
	handler := &commentHandler{
		usecase: uc,
	}

	rg.POST("/create", handler.Comment)
	rg.DELETE("/delete/:comment_id", handler.Delete)
	rg.GET("/get/:postId", handler.GetByPostID)
}

func (h *commentHandler) Comment(c *gin.Context) {
	var req domain.CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "bad request", err.Error())
		return
	}

	userID := c.GetString("user_id")

	err := h.usecase.Create(c, userID, req.PostID, req.Content)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Server Failure", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Comment Created Successfuly", nil)
}

func (h *commentHandler) Delete(c *gin.Context) {
	commentID := c.Param("comment_id")
	if commentID == "" {
		response.Error(c, http.StatusBadRequest, "bad request", "comment ID required")
		return
	}

	userID := c.GetString("user_id")

	err := h.usecase.Delete(c, userID, commentID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Server Error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "comment deleted successfully", nil)
}

func (h *commentHandler) GetByPostID(c *gin.Context) {
	postID := c.Param("postId")
	if postID == "" {
		response.Error(c, http.StatusBadRequest, "bad request", "postId is required")
		return
	}

	comments, err := h.usecase.GetByPostID(c.Request.Context(), postID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Server Error", err.Error())
		return
	}

	if comments == nil {
		comments = []domain.Comment{}
	}
	
	response.Success(c, http.StatusOK, "", gin.H{"post_id": postID, "comments": comments})
}
