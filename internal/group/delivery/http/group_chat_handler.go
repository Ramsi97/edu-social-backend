package http

import (
	"time"

	"github.com/Ramsi97/edu-social-backend/internal/group/domain"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	pingPeriod = (pongWait * 9) / 10
	pongWait   = 60 * time.Second
	writeWait  = 10 * time.Second
)

type wsClient struct {
	id uuid.UUID
	conn *websocket.Conn
	send chan []byte
}

func (c *wsClient) ID() uuid.UUID { return c.id}
func (c *wsClient) Send(msg []byte) { c.send <- msg }
func (c *wsClient) Close() error {return c.conn.Close()}

func(c *wsClient) WritePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func(){
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <- c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			
			w, err := c.conn.NextWriter(websocket.TextMessage)

			if err != nil {
				return
			}
			w.Write(message)
			
			n := len(c.send)
			for i := 0; i<n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}
			w.Close()

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *wsClient) ReadPump(h domain.Message){
		
	
}