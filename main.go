package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	cht "github.com/Robert-Duck-by-BB-SR/tps/internal/chat"
	"github.com/Robert-Duck-by-BB-SR/tps/internal/database"
	"github.com/Robert-Duck-by-BB-SR/tps/internal/models"
	"github.com/jmoiron/sqlx"
)

func assert(cond bool, expect string) {
	if cond {
		panic("shit happened: " + expect)
	}
}

var Connections = make(map[string][]net.Conn)

func main() {
	// TODO: separate user creation from running server possibly with cmd args

	database.DB = sqlx.MustConnect("sqlite3", "testing.db")

	listener, err := net.Listen("tcp4", ":6969")
	assert(err != nil, "cannot listen")
	defer listener.Close()

	chat := cht.Chat{}
	go chat.Start()
	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Cannot accept a connection")
			continue
		}
		go HandleConnection(&chat, conn)

	}
}

func HandleConnection(chat *cht.Chat, conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Client disconnected")
			break
		}
		go ParseIncoming(chat, conn, []byte(message))
	}
}

func WriteLine(conn net.Conn, message []byte) {
	// Echo message back to the client
	_, err := conn.Write(message)
	if err != nil {
		fmt.Println("Error writing to client:", err)
	}
}

func ParseIncoming(chat *cht.Chat, conn net.Conn, message []byte) {
	switch message[0] {
	case 0:
		ParseMessage(message[1:])
	case 1:
		ParseRequest(message[1:])
	case 2:
		ConnectionRequest(chat, message[1:], conn)
	default:
		WriteLine(conn, []byte("Bad request"))
	}
}

func ParseMessage(message []byte) error {
	request := string(message)
	m := strings.Split(request, string([]byte{255}))
	if len(m) < 2 {
		return fmt.Errorf("Cannot parse message, got wrong length")
	}

	// key := m[0]
	// text := m[1]

	return nil
}

func ParseRequest(message []byte) {
	request := string(message)
	_ = strings.Split(request, string([]byte{255}))
	// TODO: use the above kv pairs to figure out what are we requesting
}

func ConnectionRequest(chat *cht.Chat, key []byte, conn net.Conn) error {
	err, username := models.FetchUsername(string(key))
	if err != nil {
		return err
	}
	return chat.Connect(username, conn)
}
