package chat

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/Robert-Duck-by-BB-SR/tps/internal/models"
	"github.com/google/uuid"
)

type conn struct {
	username string
	socket   net.Conn
}

type user_chan struct {
	username string
	conn     chan []net.Conn
}

type Chat struct {
	connections       map[string][]net.Conn
	create_connection chan conn
	delete_connection chan conn
	get_connection    chan user_chan
}

func (c *Chat) Start() {
	c.connections = make(map[string][]net.Conn)
	c.create_connection = make(chan conn)
	for {
		select {
		case conn := <-c.create_connection:
			c.connections[conn.username] = append(c.connections[conn.username], conn.socket)
		case dis := <-c.delete_connection:

			new_connections := c.connections[dis.username]
			i := -1
			for j, conn := range new_connections {
				if conn == dis.socket {
					i = j
					break
				}
			}
			if i != -1 {
				new_connections = append(new_connections[:i], new_connections[i+1:]...)
			}

			c.connections[dis.username] = new_connections
		case user := <-c.get_connection:
			user.conn <- c.connections[user.username]
		}
	}
}

func (c *Chat) connect(username string, socket net.Conn) {
	c.create_connection <- conn{username: username, socket: socket}
}

func (c *Chat) disconnect(username string, socket net.Conn) {
	c.delete_connection <- conn{username: username, socket: socket}
}

func (c Chat) get_connections(username string) []net.Conn {
	user := user_chan{
		username: username,
		conn:     make(chan []net.Conn),
	}
	c.get_connection <- user
	return <-user.conn
}

func (chat *Chat) HandleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Client disconnected")
			break
		}
		chat.parse_incoming(conn, []byte(message))
	}
}

func (chat *Chat) parse_incoming(conn net.Conn, message []byte) {
	switch message[0] {
	case 0:
		chat.parse_message(conn, message[1:])
	case 1:
		chat.parse_request(conn, message[1:])
	case 2:
		err := chat.connection_request(message[1:], conn)
		if err != nil {
			write_line(conn, []byte(fmt.Sprint("Could not create a connection: ", err)))
		}
	default:
		write_line(conn, []byte("Bad request"))
	}
}

func (chat *Chat) parse_message(conn net.Conn, message []byte) {
	request := string(message)
	m := strings.Split(request, string([]byte{255})) // we're are going to use 0-127 (ascii) so 255 should be safe as a separator
	if len(m) < 4 {                                  // type, key, conversation id, message itself
		write_line(conn, []byte("Cannot parse message, got wrong length"))
	}

	key := m[0]
	conversation := m[1]
	t := m[2] // typpppee
	text := m[3]

	err, user := models.FetchUsername(key)
	if err != nil {
		write_line(conn, []byte("Could not find user with such key"))
	}

	datetime := time.Now().Format("2006-01-02 15:04")
	err = models.CreateMessage(t, user, conversation, datetime, text)
	if err != nil {
		write_line(conn, []byte("could not send message, we are checking"))
	}

	_, users := models.FetchConversationUsers(conversation)
	usernames := strings.Split(users, "|")
	for _, username := range usernames {
		conns := chat.get_connections(username)
		for _, c := range conns {
			write_line(c,
				[]byte(t),
				[]byte(user),
				[]byte(conversation),
				[]byte(datetime),
				[]byte(text),
			)
		}
	}

}

func (chat *Chat) parse_request(conn net.Conn, message []byte) {
	request := string(message)
	data := strings.Split(request, string([]byte{255}))
	if len(data) < 3 {
		write_line(conn, []byte("Cannot parse request, got wrong length"))
		return
	}
	key := data[0]
	err, username := models.FetchUsername(key)
	if err != nil {
		write_line(conn, []byte("Unregistered user"))
		return
	}
	request_type := data[1]

	switch request_type {
	case "get":
		receiver := data[2]
		switch receiver {
		case "message":
			if len(data) < 4 {
				write_line(conn, []byte("request for messages got no conversation"))
				return
			}
			_, messages := models.FetchMessages(username)
			var builder strings.Builder
			for _, messages := range messages {
				builder.WriteByte(messages.Type)
				builder.WriteByte(255)
				builder.WriteString(messages.User)
				builder.WriteByte(255)
				builder.WriteString(messages.Conversation)
				builder.WriteByte(255)
				builder.WriteString(messages.Datetime)
				builder.WriteByte(255)
				builder.WriteString(string(messages.Content))
				builder.WriteByte(254)
			}
			write_line(conn, []byte(builder.String()))
		case "conversation":
			_, conversations := models.FetchConversationsByUsername(username)
			var builder strings.Builder
			for _, conv := range conversations {
				builder.WriteString(conv.Id)
				builder.WriteByte(255)
				builder.WriteString(conv.Users)
				builder.WriteByte(254)
			}
			write_line(conn, []byte(builder.String()))
		case "users":
			_, users := models.FetchUsers()
			var builder strings.Builder
			for _, user := range users {
				builder.WriteString(user)
				builder.WriteByte(254)
			}
			write_line(conn, []byte(builder.String()))
		}

	case "create":
		users := data[2]
		id := uuid.NewString()
		models.CreateConversation(id, fmt.Sprintf("%s|%s", username, users))
	}

}

func (chat *Chat) connection_request(key []byte, conn net.Conn) error {
	err, username := models.FetchUsername(string(key))
	if err != nil {
		return err
	}
	chat.connect(username, conn)
	return nil
}

func write_line(conn net.Conn, parts ...[]byte) {
	for _, part := range parts {
		_, err := conn.Write(part)
		if err != nil {
			fmt.Println("Error writing to client:", err)
		}
		_, err = conn.Write([]byte{255})
		if err != nil {
			fmt.Println("Error writing to client:", err)
		}
	}

	_, err := conn.Write([]byte{'\n'})
	if err != nil {
		fmt.Println("Error writing to client:", err)
	}
}
