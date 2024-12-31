package models

import (
	"log"

	"github.com/Robert-Duck-by-BB-SR/tps/internal/database"
)

type User struct {
	Id       string
	Username string
	Key      []byte
}

func FetchUsername(key string) (error, string) {
	var username string

	if err := database.DB.Get(&username, "select username from user where key=?", key); err != nil {
		log.Println("cannot find username: ", err)
		return err, ""
	}

	return nil, username
}

func FetchUsers() (error, []string) {
	var users []string
	if err := database.DB.Get(&users, "select username from user"); err != nil {
		log.Println("cannot fetch users: ", err)
		return err, []string{}
	}
	return nil, users
}
