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
	h.io.On("connection", func(clients ...any) {
		if len(clients) == 0 {
			return
		}
		client, ok := clients[0].(*socket.Socket)
		if !ok {
			return
		}
		log.Printf("User connected: %s", client.Id())
		client.On("join_group", func(data ...any) {
			if len(data) == 0 || data[0] == nil {
				return
			}
			groupIDstr, ok := data[0].(string)
			if !ok {
				log.Printf("Error: join_group expected string, got %T", data[0])
				return
			}
			client.Join(socket.Room(groupIDstr))
			log.Printf("User %s joined group %s", client.Id(), groupIDstr)
		})
		client.On("send_message", func(data ...any) {
			if len(data) == 0 || data[0] == nil {
				return
			}
			msgData, ok := data[0].(map[string]any)
			if !ok {
				log.Printf("Error: send_message expected map, got %T", data[0])
				return
			}
			// Use safe retrieval with "ok" idiom to prevent nil-to-string panics
			groupIDStr, groupOK := msgData["group_id"].(string)
			content, contentOK := msgData["content"].(string)
			if !groupOK || !contentOK {
				log.Printf("Error: group_id or content missing/nil in payload")
				return
			}
			// Get sender ID from SetData(userID) in middleware
			senderIDRaw := client.Data()
			if senderIDRaw == nil {
				log.Printf("Error: sender ID not found in socket data")
				return
			}
			senderID, ok := senderIDRaw.(uuid.UUID)
			if !ok {
				log.Printf("Error: senderID in socket data is not uuid.UUID, got %T", senderIDRaw)
				return
			}
			groupID, err := uuid.Parse(groupIDStr)
			if err != nil {
				log.Printf("Error parsing group ID: %v", err)
				return
			}
			message := &domain.Message{
				GroupID:  groupID,
				Content:  content,
				AuthorID: senderID,
			}
			err = h.chatUsecase.SendMessage(context.Background(), message)
			if err != nil {
				client.Emit("error", err.Error())
				return
			}
			// Broadcast to room
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

		userIDStr, err := auth.ValidateToken(token)

		if err != nil {
			next(socket.NewExtendedError("Invalid token", nil))
			return
		}

		userUUID, err := uuid.Parse(userIDStr)
		if err != nil {
			next(socket.NewExtendedError("Invalid user ID", nil))
			return
		}

		s.SetData(userUUID)

		next(nil)
	})
}