package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/coder/websocket"
)

func main() {
	fmt.Println("Started client..")
	done := make(chan struct{})
	go ConnectAndPingServer(done)
	<-done
	fmt.Println("Exiting program")
}

func ConnectAndPingServer(done chan struct{}) {
	fmt.Println("Connecting Client")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, "ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatalf("Failed to connext: %v", err)
		close(done)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "closing")
	fmt.Println("Client is connected to server")
	for i := range 10 {
		err = conn.Write(ctx, websocket.MessageText, []byte(fmt.Sprintf("Ping %v from client.", i)))
		if err != nil {
			log.Fatalf("Failed to write message: %v", err)
			close(done)
			return
		}
		fmt.Println("Pinged server")
		time.Sleep(1 * time.Second)
	}
	close(done)
}
