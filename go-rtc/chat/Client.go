package chat

import (
	"context"
	"errors"
	"io"
	"log"
	"time"

	"github.com/coder/websocket"
)

type Client struct {
	Name string
	conn *websocket.Conn
	room *chatRoom
}

func newClient(clientName string, c *websocket.Conn, r *chatRoom) *Client {
	return &Client{Name: clientName, conn: c, room: r}
}

func (c *Client) Disconnect() {
	log.Printf("Disconnecting client: %s", c.Name)
	c.conn.Close(websocket.StatusNormalClosure, "Client disconnecting")
	if c.room != nil {
		c.room.RemoveClient(c)
	}
}

func (c *Client) Write(message string, ctx context.Context) error {
	return c.conn.Write(ctx, websocket.MessageText, []byte(message))
}

func (c *Client) ReadClient(ctx context.Context) (string, error) {
	buf := make([]byte, 1024)
	for {
		typ, reader, err := c.conn.Reader(ctx)
		if err != nil {
			status := websocket.CloseStatus(err)
			switch status {
			case websocket.StatusNormalClosure:
				log.Println("Websocket closed normally")
			case websocket.StatusGoingAway:
				log.Println("Client is going away (browser/tab closed)")
			case websocket.StatusAbnormalClosure:
				log.Println("Websocket closed abnormally")
			default:
				log.Printf("Websocket closed with status: %v (err: %v)", status, err)
			}
			break
		}

		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			log.Println("error reading message body:", err)
			break
		}

		if typ == websocket.MessageText || typ == websocket.MessageBinary {
			log.Printf("message (%d bytes): %s", n, buf[:n])
			c.room.PublishMessage(Message{Sender: c.Name, Message: string(buf[:n]), Timestamp: time.Now().Format("15:04:05")}, ctx)
			return string(buf[:n]), nil
		}
	}
	return "", errors.New("failed to read message for unknown reasons")
}
