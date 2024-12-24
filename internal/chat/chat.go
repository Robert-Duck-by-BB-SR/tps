package chat

import "net"

type conn struct {
	Username string
	Conn     net.Conn
}

type Chat struct {
	connections   map[string][]net.Conn
	connection    chan conn
	disconnection chan conn
	// message     chan string
	// response    chan string
}

func (c *Chat) Start() {
	c.connections = make(map[string][]net.Conn)
	c.connection = make(chan conn)
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

func (c *Chat) Connect(username string, socket net.Conn) error {
	c.connection <- conn{Username: username, Conn: socket}
	return nil
}

func (c *Chat) Disconnect(username string, socket net.Conn) error {
	c.disconnection <- conn{Username: username, Conn: socket}
	return nil
}
