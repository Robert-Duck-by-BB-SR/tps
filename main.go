package main

import (
	"fmt"
	"net"

	cht "github.com/Robert-Duck-by-BB-SR/tps/internal/chat"
	"github.com/Robert-Duck-by-BB-SR/tps/internal/database"
	"github.com/jmoiron/sqlx"
)

func assert(cond bool, expect string) {
	if cond {
		panic("shit happened: " + expect)
	}
}

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
		go chat.HandleConnection(conn)

	}
}
