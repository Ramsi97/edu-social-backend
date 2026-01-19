package http

import (
	"net/http"

	"github.com/Ramsi97/edu-social-backend/internal/chat/domain"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	usecase domain.ChatUseCase
}

func NewChatHandler(rg *gin.RouterGroup, uc domain.ChatUseCase) {
	handler := &ChatHandler{usecase: uc}

	rg.POST("/send", handler.SendMessage)
	rg.GET("/history/:room_id", handler.GetMessages)
}

func (h *ChatHandler) SendMessage(ctx *gin.Context) {
	var msg domain.Message
	if err := ctx.ShouldBindJSON(&msg); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.SendMessage(ctx, &msg); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "message sent"})
}

func (h *ChatHandler) GetMessages(ctx *gin.Context) {
	roomID := ctx.Param("room_id")
	messages, err := h.usecase.GetMessages(ctx, roomID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, messages)
}
	