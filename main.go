package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

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
		ParseIncoming(chat, conn, []byte(message))
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
		ParseMessage(chat, conn, message[1:])
	case 1:
		ParseRequest(message[1:])
	case 2:
		err := ConnectionRequest(chat, message[1:], conn)
		if err != nil {
			WriteLine(conn, []byte(fmt.Sprint("Could not create a connection: ", err)))
		}
	default:
		WriteLine(conn, []byte("Bad request"))
	}
}

func ParseMessage(chat *cht.Chat, conn net.Conn, message []byte) {
	request := string(message)
	m := strings.Split(request, string([]byte{255})) // we're are going to use 0-127 (ascii) so 255 should be safe as a separator
	if len(m) < 4 {                                  // type, key, conversation id, message itself
		WriteLine(conn, []byte("Cannot parse message, got wrong length"))
	}

	key := m[0]
	conversation := m[1]
	t := m[2] // typpppee
	text := m[3]

	err, user := models.FetchUsername(key)
	if err != nil {
		WriteLine(conn, []byte("Could not find user with such key"))
	}

	datetime := time.Now().Format("2006-01-02 15:04")
	err = models.CreateMessage(t, user, conversation, datetime, text)
	if err != nil {
		WriteLine(conn, []byte("could not send message, we are checking"))
	}

	_, users := models.FetchConversationUsers(conversation)
	usernames := strings.Split(users, "|")
	for _, username := range usernames {
		conns := chat.GetConnections(username)
		for _, c := range conns {
			WriteLine(c, []byte(fmt.Sprint(
				t, string([]byte{255}),
				user, string([]byte{255}),
				conversation, string([]byte{255}),
				datetime, string([]byte{255}),
				text, "\n",
			)))
		}
	}

}

func ParseRequest(conn net.Conn, message []byte) {
	request := string(message)
	data := strings.Split(request, string([]byte{255}))
	if len(data) < 3 {
		WriteLine(conn, []byte("Cannot parse request, got wrong length"))
	}
	request_type := data[0]
	receiver := data[1]
	key := data[2]

	switch request_type {
	case "get":
		switch receiver {
		case "message":
		case "conversation":
		case "users":
		}

	case "create":
	}

}

func ConnectionRequest(chat *cht.Chat, key []byte, conn net.Conn) error {
	err, username := models.FetchUsername(string(key))
	if err != nil {
		return err
	}
	return chat.Connect(username, conn)
}
