package chat

import (
	"bufio"
	"fmt"
	"log"
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
	log.Println("Starting tps server")
	c.connections = make(map[string][]net.Conn)
	c.create_connection = make(chan conn)
	c.get_connection = make(chan user_chan)
	for {
		select {
		case conn := <-c.create_connection:
			log.Println("creating connection")
			c.connections[conn.username] = append(c.connections[conn.username], conn.socket)
			log.Println(c.connections)
		case dis := <-c.delete_connection:

			log.Println("creating disconnections")
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
			log.Println("fetching users")
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
	log.Println("incoming,", message)
	switch message[0] {
	case 0:
		chat.parse_message(conn, message[1:len(message)-1])
	case 1:
		chat.parse_request(conn, message[1:])
	case 2:
		key := message[1 : len(message)-1]
		err, username := chat.connection_request(key, conn)
		if err != nil {
			write_line(conn, []byte(fmt.Sprint("Could not create a connection: ", err)))
			return
		}
		write_line(conn, []byte(username))
	default:
		write_line(conn, []byte("Bad request"))
	}
}

func (chat *Chat) parse_message(conn net.Conn, message []byte) {
	request := strings.Trim(string(message), "\n")
	m := strings.Split(request, string([]byte{255})) // we're are going to use 0-127 (ascii) so 255 should be safe as a separator
	log.Println("parsing message:", m)
	if len(m) < 4 { // type, key, conversation id, message itself
		write_line(conn, []byte("Cannot parse message, got wrong length"))
		return
	}

	key := m[0]
	conversation := m[1]
	t := []byte(m[2]) // typpppee
	text := []byte(m[3])

	err, user := models.FetchUsername(string(key))
	if err != nil {
		write_line(conn, []byte("Could not find user with such key"))
		return
	}

	datetime := time.Now().Format("2006-01-02 15:04")
	err = models.CreateMessage(t[0], user, conversation, datetime, text)
	if err != nil {
		write_line(conn, []byte("could not send message, we are checking"))
		return
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
	request := strings.Trim(string(message), "\n")
	data := strings.Split(request, string([]byte{255}))
	log.Println("data:", data)
	if len(data) < 3 {
		write_line(conn, []byte("Cannot parse request, got wrong length"))
		return
	}
	key := data[0]
	err, username := models.FetchUsername(string(key))
	if err != nil {
		write_line(conn, []byte("Bad user"))
		return
	}
	log.Println("proceding with the request")
	request_type := data[1]

	switch request_type {
	case "get":
		receiver := data[2]
		var builder strings.Builder
		switch receiver {
		case "message":
			if len(data) < 4 {
				write_line(conn, []byte("request for messages got no conversation"))
				return
			}
			log.Println(data)
			_, messages := models.FetchMessages(data[3])
			log.Println("got:", messages)
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
		case "conversation":
			err, conversations := models.FetchConversationsByUsername(username)
			if err != nil {
				log.Println("cannot get convos", err)
			}
			for _, conv := range conversations {
				builder.WriteString(conv.Id)
				builder.WriteByte(255)
				builder.WriteString(conv.Users)
				builder.WriteByte(254)
			}
		case "users":
			_, users := models.FetchUsers()
			for _, user := range users {
				builder.WriteString(user)
				builder.WriteByte(254)
			}
		}
		log.Println(builder.String())
		write_line(conn, []byte(builder.String()))

	case "create":
		users := data[2]
		id := uuid.NewString()
		models.CreateConversation(id, fmt.Sprintf("%s|%s", username, users))
		write_line(conn, []byte(id))
	}

}

func (chat *Chat) connection_request(key []byte, conn net.Conn) (error, string) {
	err, username := models.FetchUsername(string(key))
	if err != nil {
		return err, ""
	}
	chat.connect(username, conn)
	return nil, username
}

func write_line(conn net.Conn, parts ...[]byte) {
	var builder []byte
	for i, part := range parts {
		builder = append(builder, part...)
		if i != len(parts)-1 {
			builder = append(builder, 255)
		}
	}

	builder = append(builder, '\n')
	log.Println("i don't know what is happening at this point", builder)
	_, err := conn.Write(builder)
	if err != nil {
		fmt.Println("Error writing to client:", err)
	}
}
