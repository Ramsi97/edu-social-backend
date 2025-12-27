package socket

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
)

// StartSocketServer starts the Socket.IO server
func StartSocketServer() (*socketio.Server, error) {
	// For v1.7.0, use this initialization
	server := socketio.NewServer(nil)
	if server == nil {
		log.Fatal("Failed to create socket.io server")
	}

	// On client connect
	server.OnConnect("/", func(s socketio.Conn) error {
		log.Println("Client connected:", s.ID())
		return nil
	})

	// Join a room
	server.OnEvent("/", "join_room", func(s socketio.Conn, roomID string) {
		s.Join(roomID)
		log.Println("Client", s.ID(), "joined room", roomID)
	})

	// Handle chat messages
	server.OnEvent("/", "send_message", func(s socketio.Conn, msg map[string]string) {
		roomID := msg["room_id"]
		server.BroadcastToRoom("/", roomID, "new_message", msg)
		log.Println("Message sent to room", roomID)
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("Socket error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("Client disconnected:", s.ID(), reason)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatal("socketio listen error:", err)
		}
	}()
	log.Println("Socket.IO server started")
	return server, nil
}