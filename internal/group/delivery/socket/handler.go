package socket

import (
	"context"
	"log"

	"github.com/Ramsi97/edu-social-backend/internal/group/domain"
	"github.com/google/uuid"
	"github.com/zishang520/socket.io/v2/socket"
	"github.com/Ramsi97/edu-social-backend/pkg/auth"
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
			
			senderIDRaw := client.Data() 
    		senderID := senderIDRaw.(uuid.UUID) 	
			groupID, _ := uuid.Parse(groupIDStr)
			message := &domain.Message{
				GroupID: groupID,
				Content: content,
				AuthorID: senderID,
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

func (h *socketHandler) RegisterMiddleWare() {
	h.io.Use(func (s *socket.Socket, next func (*socket.ExtendedError))  {
		
		authData := s.Handshake().Auth
		authMap, ok := authData.(map[string]any)

		if !ok {
			next(socket.NewExtendedError("Authentication failed", nil))
			return
		}

		token, tokenOK := authMap["token"].(string)

		if !tokenOK || token == "" {
			next(socket.NewExtendedError("Token required", nil))
			return
		}

		userID, err := auth.ValidateToken(token)

		if err != nil {
			next(socket.NewExtendedError("Invalid token", nil))
			return
		}

		s.SetData(userID)

		next(nil)
	})
}