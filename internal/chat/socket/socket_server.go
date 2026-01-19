package socket

import (
	"context"
	"log"

	"github.com/Ramsi97/edu-social-backend/internal/chat/domain"
	"github.com/Ramsi97/edu-social-backend/pkg/auth"
	"github.com/google/uuid"
	"github.com/zishang520/socket.io/v2/socket"
)

type SocketHandler struct {
	io          *socket.Server
	chatUsecase domain.ChatUseCase
}

func NewSocketHandler(
	io *socket.Server,
	chatUC domain.ChatUseCase,
) *SocketHandler {
	return &SocketHandler{
		io:          io,
		chatUsecase: chatUC,
	}
}

type SendMessageDTO struct {
	RoomID  string `json:"room_id"`
	Content string `json:"content"`
}
func (h *SocketHandler) RegisterEvents() {

	h.io.On("connection", func(args ...any) {

		client, ok := args[0].(*socket.Socket)
		if !ok {
			return
		}

		log.Println("Client connected:", client.Id())

		// -------------------------
		// JOIN GROUP
		// -------------------------
		client.On("join_group", func(data ...any) {

			groupIDStr, ok := data[0].(string)
			if !ok {
				client.Emit("error", "invalid group id")
				return
			}

			client.Join(socket.Room(groupIDStr))
			log.Printf("User %s joined group %s", client.Id(), groupIDStr)
		})

		// -------------------------
		// SEND MESSAGE
		// -------------------------
		client.On("send_message", func(data ...any) {

			payload, ok := data[0].(map[string]any)
			if !ok {
				client.Emit("error", "invalid payload")
				return
			}

			roomIDStr, _ := payload["group_id"].(string)
			content, _ := payload["content"].(string)

			senderIDStr, ok := client.Data().(string)
			if !ok {
				client.Emit("error", "unauthorized")
				return
			}

			roomID, err := uuid.Parse(roomIDStr)
			if err != nil {
				client.Emit("error", "invalid group id")
				return
			}

			senderID, err := uuid.Parse(senderIDStr)
			if err != nil {
				client.Emit("error", "invalid sender id")
				return
			}

			msg := &domain.Message{
				RoomID:  roomID,
				SenderID: senderID,
				Content:  content,
			}

			if err := h.chatUsecase.SendMessage(context.Background(), msg); err != nil {
				client.Emit("error", err.Error())
				return
			}

			h.io.
				To(socket.Room(roomIDStr)).
				Emit("new_message", msg)
		})

		// -------------------------
		// DISCONNECT
		// -------------------------
		client.On("disconnect", func(...any) {
			log.Println("User disconnected:", client.Id())
		})
	})
}



func (h *SocketHandler) RegisterMiddleWare() {
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