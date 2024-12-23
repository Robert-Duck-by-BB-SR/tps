package main

import (
	"bufio"
	"fmt"
	"net"
)

func assert(cond bool, expect string) {
	if cond {
		panic("shit happened: " + expect)
	}
}

var Connections = make(map[string][]net.Conn)

func main() {
	listener, err := net.Listen("tcp4", ":6969")
	assert(err != nil, "cannot listen")
	defer listener.Close()

	chat := Chat{}
	go chat.Start()
	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Cannot accept a connection")
			continue
		}
		go HandleConnection(conn)

	}
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Client disconnected")
			break
		}
		go ParseIncoming([]byte(message))
	}
}

func WriteLine(conn net.Conn, message []byte) {
	// Echo message back to the client
	_, err := conn.Write(message)
	if err != nil {
		fmt.Println("Error writing to client:", err)
	}
}

func ParseIncoming(message []byte) {
	switch message[0] {
	case 0:
		ParseMessage(message[1:])
	case 1:
		ParseRequest(message[1:])
	case 2:
		ConnectionRequest(message[1:])
	}
}

func ParseMessage(message []byte) {

}

func ParseRequest(message []byte) {}

func ConnectionRequest(message []byte) {

}

type Connection struct {
	Username string
	Conn     net.Conn
}

type Chat struct {
	connections   map[string][]net.Conn
	connection    chan Connection
	disconnection chan Connection
	// message     chan string
	// response    chan string
}

func (c *Chat) Start() {
	c.connections = make(map[string][]net.Conn)
	c.connection = make(chan Connection)
	for {
		select {
		case conn := <-c.connection:
			c.connections[conn.Username] = append(c.connections[conn.Username], conn.Conn)
		case dis := <-c.disconnection:

			new_connections := c.connections[dis.Username]
			i := -1
			for j, conn := range new_connections {
				if conn == dis.Conn {
					i = j
					break
				}
			}
			if i != -1 {
				new_connections = append(new_connections[:i], new_connections[i+1:]...)
			}

			c.connections[dis.Username] = new_connections
		}
	}
}

func (c *Chat) Connect(key string, socket net.Conn) error {
	// TODO: retrieve username by key from db
	c.connection <- Connection{Username: key, Conn: socket}
	return nil
}

func (c *Chat) Disconnect(key string, socket net.Conn) error {
	// TODO: retrieve username by key from db
	c.connection <- Connection{Username: key, Conn: socket}
	return nil
}
