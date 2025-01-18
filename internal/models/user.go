package models

import (
	"encoding/hex"
	"log"

	"github.com/Robert-Duck-by-BB-SR/tps/internal/database"
)

type User struct {
	Id       string
	Username string
	Key      []byte
}

func FetchUsername(key_hex string) (error, string) {
	log.Println("decoding:", key_hex)
	key, err := hex.DecodeString(key_hex)
	if err != nil {
		log.Println("Cannot decode key, ", err)
		return err, ""
	}

	log.Println("decoded", key)

	var username string

	if err = database.DB.Get(&username, "select username from user where key=?", key); err != nil {
		log.Println("cannot find username: ", err)
		return err, ""
	}

	return nil, username
}

func FetchUsers() (error, []string) {
	var users []string
	if err := database.DB.Select(&users, "select username from user"); err != nil {
		log.Println("cannot fetch users: ", err)
		return err, []string{}
	}
	return nil, users
}

func CreateUser(id, username string, key [32]byte) error {
	if _, err := database.DB.Exec(
		"insert into user values (?, ?, ?)",
		id, username, key[:]); err != nil {
		log.Println("user was not created: ", err)
	}
	return nil
}
