package http

import (
	"net/http"

	"github.com/Ramsi97/edu-social-backend/internal/group/domain"
	"github.com/Ramsi97/edu-social-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GroupHandler struct {
	usecase domain.GroupChatUseCase
}

func NewGroupHandler(uc domain.GroupChatUseCase, r *gin.RouterGroup) {
	handler := GroupHandler{usecase: uc}

	r.POST("/groups", handler.CreateGroup)
	r.POST("/groups/:name/join", handler.JoinGroup)
	r.POST("/groups/:name/leave", handler.LeaveGroup)
	r.GET("/groups/:group_id/messages", handler.GetMessages)

}

func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	userIDStr := c.GetString("userID")
	userID, err  := uuid.Parse(userIDStr)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	groupID, err := h.usecase.CreateGroup(c.Request.Context(), userID, req.Name)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"group_id": groupID})
}

func (h *GroupHandler) JoinGroup(c *gin.Context) {
	groupName := c.Param("name")
	userIDStr := c.GetString("userID")
	userID, err  := uuid.Parse(userIDStr)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	if err := h.usecase.JoinGroup(c.Request.Context(), groupName, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *GroupHandler) LeaveGroup(c *gin.Context) {
	groupName := c.Param("name")
	userID := c.MustGet("userID").(uuid.UUID)

	if err := h.usecase.LeaveGroup(c.Request.Context(), groupName, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *GroupHandler) GetMessages(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	msgs, err := h.usecase.GetMessages(c.Request.Context(), groupID, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, msgs)
}
