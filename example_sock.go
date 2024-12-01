package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

func main() {
	// Define the WebSocket server URL (replace with your server's address)
	serverAddr := "ws://localhost:6969/websocket"

	// Parse the URL
	u, err := url.Parse(serverAddr)
	if err != nil {
		log.Fatalf("Invalid server address: %v", err)
	}

	// Dial the WebSocket server
	log.Printf("Connecting to %s...", serverAddr)
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	log.Println("Connected to server")

	// Set up interrupt handling for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Communication goroutine
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}
			log.Printf("Received: %s", message)
		}
	}()

	// Send a test message
	testMessage := "Message"
	err = conn.WriteMessage(websocket.TextMessage, []byte(testMessage))
	if err != nil {
		log.Fatalf("Write error: %v", err)
	}
	log.Printf("Sent: %s", testMessage)

	// Wait for interrupt signal to close the connection
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("Interrupt received, closing connection")
			// Send a close message to the server
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Close error:", err)
			}
			return
		}
	}
}
