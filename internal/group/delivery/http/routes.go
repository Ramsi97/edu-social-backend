package http

import "github.com/gin-gonic/gin"

func RegisterGroupRoutes(r *gin.RouterGroup, h *GroupHandler) {
	r.POST("/groups", h.CreateGroup)
	r.POST("/groups/:name/join", h.JoinGroup)
	r.POST("/groups/:name/leave", h.LeaveGroup)
	r.GET("/groups/:group_id/messages", h.GetMessages)
}

