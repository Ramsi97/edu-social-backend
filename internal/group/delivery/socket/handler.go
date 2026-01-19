package socket

import (
	"context"
	"log"

	"github.com/Ramsi97/edu-social-backend/internal/group/domain"
	"github.com/google/uuid"
	"github.com/zishang520/socket.io/v2/socket"
)



type socketHandler struct {
	chatUsecase domain.GroupChatUseCase
	io *socket.Server
}

func NewSocketHandler(io *socket.Server, uc domain.GroupChatUseCase) *socketHandler {
	return &socketHandler{
		chatUsecase: uc,
		io: io,
	}
}

func (h *socketHandler) RegisterEvents() {

	h.io.On("connection", func (clients ...any)  {
		client := clients[0].(*socket.Socket)
		log.Printf("User connected: %s", client.Id())

		client.On("join_group", func(data ...any) {
			groupIDstr := data[0].(string)
			client.Join(socket.Room(groupIDstr))
			log.Printf("User %s joined group %s", client.Id(), groupIDstr)
		})

		client.On("send_message", func(data ...any) {
			msgData := data[0].(map[string]any)
			groupIDStr := msgData["group_id"].(string)
			content := msgData["content"].(string)
			senderIDStr := msgData["sender_id"].(string)
			
			groupID, _ := uuid.Parse(groupIDStr)
			senderID, _ := uuid.Parse(senderIDStr)
			message := &domain.Message{
				GroupID: groupID,
				Content: content,
				SenderID: senderID,
			}
			err := h.chatUsecase.SendMessage(context.Background(), message)
			if err != nil {
				client.Emit("error", err.Error())
				return
			}
			h.io.To(socket.Room(groupIDStr)).Emit("new_message", content)
		})

		client.On("disconnect", func(...any) {
			log.Println("User disconnected")
		})
	})
}