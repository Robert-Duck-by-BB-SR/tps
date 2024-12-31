package models

import (
	"log"

	"github.com/Robert-Duck-by-BB-SR/tps/internal/database"
)

type Conversation struct {
	Id, Users string
}

func FetchConversationUsers(id string) (error, string) {
	var users string
	if err := database.DB.Get(&users, "select users from conversation where id=?", id); err != nil {
		log.Println("could not fetch users from conversation: ", err)
		return err, ""
	}
	return nil, users
}

func FetchConversationsByUsername(username string) (error, []Conversation) {
	var conversations []Conversation
	if err := database.DB.Select(&conversations, "select * from conversation where users like '%?%'", username); err != nil {
		log.Println("could not fetch users from conversation: ", err)
		return err, []Conversation{}
	}
	return nil, conversations
}

func CreateConversation(id, users string) error {
	if _, err := database.DB.Exec("insert into conversation values (?, ?)", id, users); err != nil {
		return err
	}
	return nil
}
