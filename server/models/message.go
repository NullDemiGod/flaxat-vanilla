package models

import (
	"time"
	"flaxat/server/db"
)


type Message struct {
	ID int				`json:"id"`
	SenderID int		`json:"sender_id"`
	ChatID int			`json:"chat_id"`
	Content string		`json:"content"`
	CreatedAt time.Time	`json:"created_at"`
}

func CreateMessage(senderID int, chatID int, content string) (Message, error) {
	query := `INSERT INTO messages (sender_id, chat_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, sender_id, chat_id, content, created_at`

	var newMessage Message
	err := db.DB.QueryRow(query, senderID, chatID, content).
		Scan(&newMessage.ID, &newMessage.SenderID, &newMessage.ChatID, &newMessage.Content, &newMessage.CreatedAt)

	if err != nil {
		return newMessage, err
	}

	_, err = db.DB.Exec(`UPDATE chats SET updated_at = NOW() WHERE id = $1`, chatID)

	return newMessage, err
}

func GetChatMessages(chatID int) ([]Message, error) {
	query := `SELECT id, sender_id, chat_id, content, created_at FROM messages
		WHERE chat_id = $1 ORDER BY created_at ASC`

	rows, err := db.DB.Query(query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messageList []Message
	for rows.Next() {
		var current Message

		err := rows.Scan(&current.ID, &current.SenderID, &current.ChatID, &current.Content, &current.CreatedAt)
		if err != nil {
			return nil, err
		}
		
		messageList = append(messageList, current)
	}
	return messageList, nil
}
