package models

import (
	"log"

	"github.com/Robert-Duck-by-BB-SR/tps/internal/database"
	"github.com/google/uuid"
)

type Message struct {
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
