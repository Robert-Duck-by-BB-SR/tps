package models

import (
	"log"

	"github.com/Robert-Duck-by-BB-SR/tps/internal/database"
	"github.com/google/uuid"
)

type Message struct {
	Id, User, Conversation, Datetime string
	Type                             byte
	Content                          []byte
}

func CreateMessage(t, user, conversation, datetime, message string) error {
	id := uuid.NewString()
	if _, err := database.DB.Exec(
		"insert into message values (?, ?, ?, ?, ?, ?)",
		id, t, user, conversation, datetime, message); err != nil {
		log.Println("message was not created: ", err)
	}
	return nil
}

func FetchMessages(conversation string) (error, []Message) {
	var messages []Message

	if err := database.DB.Select(&messages, "select * from message where conversation=?", conversation); err != nil {
		log.Println("cannot fetch messages, ", err)
		return err, messages
	}

	return nil, messages
}
