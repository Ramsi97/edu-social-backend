package http

import (
	"net/http"
	"time"

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

	r.POST("/create", handler.CreateGroup)
	r.POST("/join/:name", handler.JoinGroup)
	r.POST("/leave/:name", handler.LeaveGroup)
	r.GET("/messages/:group_id", handler.GetMessages)
	r.GET("/", handler.GetGroups)

}

func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	groupID, err := h.usecase.CreateGroup(c.Request.Context(), userID, req.Name)
	if err != nil {
		response.Error(c, http.StatusConflict, "", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "", gin.H{"group_id": groupID})
}

func (h *GroupHandler) JoinGroup(c *gin.Context) {
	groupName := c.Param("name")
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	if err := h.usecase.JoinGroup(c.Request.Context(), groupName, userID); err != nil {
		response.Error(c, http.StatusBadRequest, "", err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *GroupHandler) LeaveGroup(c *gin.Context) {
	groupName := c.Param("name")
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "", err.Error())
		return
	}

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

func (h *GroupHandler) GetGroups(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	groups, err := h.usecase.GetGroupsForUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch groups"})
		return
	}

	// 3️⃣ Convert domain objects to response DTOs
	type groupResponse struct {
		ID      uuid.UUID `json:"id"`
		Name    string    `json:"name"`
		OwnerID uuid.UUID `json:"owner_id"`
		Created string    `json:"created_at"`
	}

	res := make([]groupResponse, len(groups))
	for i, g := range groups {
		res[i] = groupResponse{
			ID:      g.ID,
			Name:    g.Name,
			OwnerID: g.OwnerID,
			Created: g.CreatedAt.Format(time.RFC3339),
		}
	}

	// 4️⃣ Send JSON response
	c.JSON(200, gin.H{"groups": res})
}