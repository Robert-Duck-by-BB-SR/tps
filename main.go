package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	cht "github.com/Robert-Duck-by-BB-SR/tps/internal/chat"
	"github.com/Robert-Duck-by-BB-SR/tps/internal/database"
	"github.com/Robert-Duck-by-BB-SR/tps/internal/hash"
	"github.com/Robert-Duck-by-BB-SR/tps/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func assert(cond bool, expect string) {
	if cond {
		panic("shit happened: " + expect)
	}
}

func main() {
	database.DB = sqlx.MustConnect("sqlite3", "testing.db")

	username := flag.String("create", "", "Create a user with given username")
	flag.Parse()
	if *username == "" {

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
	} else {
		file, err := os.Open(".env")
		if err != nil {
			panic(err)
		}
		bytes, err := io.ReadAll(file)
		if err != nil {
			panic(err)
		}
		secret := strings.Split(string(bytes), "=")[1]
		id := uuid.NewString()
		hashed := hash.Encode([]byte(*username+id), []byte(secret))
		key := sha256.Sum256(hashed)
		log.Println(key)
		err = models.CreateUser(id, *username, key)
		if err != nil {
			panic(err)
		}
		fmt.Println(hex.EncodeToString(key[:]))
	}
}
