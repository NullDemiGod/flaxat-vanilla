package models

import (
	"flaxat/server/db"
	"time"
)


type Chat struct {
	ID int						`json:"id"`
	Member1 int					`json:"member_1"`
	Member2 int					`json:"member_2"`
	LastMessage string 			`json:"last_message"`
	LastMessageSenderID int    	`json:"last_message_sender_id"`
	CreatedAt time.Time			`json:"created_at"`
	UpdatedAt time.Time 		`json:"updated_at"`
}

func CreateChat(member1 int, member2 int) (Chat, error) {
	query := `INSERT INTO chats (member_1, member_2)
		VALUES ($1, $2)
		RETURNING id, member_1, member_2, created_at, updated_at
	`

	if member1 > member2 {
		member1, member2 = member2, member1
	}

	var newChat Chat
	err := db.DB.QueryRow(query, member1, member2).
		Scan(&newChat.ID, &newChat.Member1, &newChat.Member2, &newChat.CreatedAt, &newChat.UpdatedAt)

	return newChat, err
}

func GetChatByID(id int) (Chat, error) {
	query := `SELECT id, member_1, member_2, created_at, updated_at FROM chats WHERE id = $1`

	var chat Chat
	err := db.DB.QueryRow(query, id).Scan(&chat.ID, &chat.Member1, &chat.Member2, &chat.CreatedAt, &chat.UpdatedAt)

	return chat, err
}


func GetAllUserChats(userID int) ([]Chat, error) {
	query := `SELECT c.id, c.member_1, c.member_2, c.created_at, c.updated_at,
		COALESCE((SELECT content FROM messages m WHERE m.chat_id = c.id ORDER BY created_at DESC LIMIT 1), '') as last_message,
		COALESCE((SELECT sender_id FROM messages m WHERE m.chat_id = c.id ORDER BY created_at DESC LIMIT 1), 0) as last_message_sender_id
		FROM chats c
		WHERE c.member_1 = $1 OR c.member_2 = $1
		ORDER BY c.updated_at DESC`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatList []Chat
	for rows.Next() {
		var current Chat
		
		err := rows.Scan(
			&current.ID, 
			&current.Member1, 
			&current.Member2, 
			&current.CreatedAt, 
			&current.UpdatedAt, 
			&current.LastMessage,
			&current.LastMessageSenderID,
		)
		
		if err != nil {
			return nil, err
		}
		
		chatList = append(chatList, current)
	}
	
	if chatList == nil {
		return []Chat{}, nil
	}
	
	return chatList, nil
}
