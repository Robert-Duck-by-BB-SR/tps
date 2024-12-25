package models

import (
	"log"

	"github.com/Robert-Duck-by-BB-SR/tps/internal/database"
)

type Conversation struct{}

func FetchConversationUsers(id string) (error, string) {
	var users string
	if err := database.DB.Get(&users, "select users from conversation where id=?", id); err != nil {
		log.Println("could not fetch users from conversation: ", err)
		return err, ""
	}
	return nil, users
}
