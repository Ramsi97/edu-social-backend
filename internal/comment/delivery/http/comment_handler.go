package http

import (
	"net/http"

	"github.com/Ramsi97/edu-social-backend/internal/comment/domain"
	"github.com/gin-gonic/gin"

)

type commentHandler struct {
	usecase domain.CommentUseCase
}


func NewCommentHandler(rg *gin.RouterGroup,uc domain.CommentUseCase) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")


	err := h.usecase.Create(c, userID, req.PostID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "comment created successfully"})
}

func (h *commentHandler) Delete(c *gin.Context) {
	commentID := c.Param("commentId")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "commentId is required"})
		return
	}

	userID := c.GetString("user_id")

	err := h.usecase.Delete(c, userID, commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment deleted successfully"})
}

func (h *commentHandler) GetByPostID(c *gin.Context) {
	postID := c.Param("postId")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "postId is required"})
		return
	}

	comments, err := h.usecase.GetByPostID(c.Request.Context(), postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if comments == nil {
		comments = []domain.Comment{}
	}

	c.JSON(http.StatusOK, comments)
}