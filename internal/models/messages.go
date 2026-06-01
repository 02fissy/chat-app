package models

import (
	"time"
	"database/sql"
)


type Messages struct {
	MessageID int
	SenderID int
	ConversationID int
	Content string
	CreatedAt time.Time
	ParentID int
} 
type MessageModel struct{
	DB *sql.DB
}

func(m *MessageModel) Insert(senderID, conversationID int, content string, parentID int) error {
	stmt:= `INSERT INTO messages(sender_id, conversation_id, content, parent_id, created_at) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)`
	_, err := m.DB.Exec(stmt, senderID, conversationID, content, parentID)
	if err != nil {
		return err
	}
	return nil
}
func(m *MessageModel) GetMessages(conversationID int) ([]Messages, error) {
	stmt := `SELECT message_id, sender_id, conversation_id, content, created_at, parent_id FROM messages WHERE conversation_id = ?`
	rows, err := m.DB.Query(stmt, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Messages
	for rows.Next() {
		var m Messages
		if err := rows.Scan(&m.MessageID, &m.SenderID, &m.ConversationID, &m.Content, &m.CreatedAt, &m.ParentID); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
func(m *MessageModel) EditMessage(messageID int, content string) error {
	stmt := `UPDATE messages SET content = ? WHERE message_id = ?`
	_, err := m.DB.Exec(stmt, content, messageID)
	if err != nil {
		return err
	}
	return nil
}
func(m *MessageModel) DeleteMessage(messageID int) error {
	stmt := `DELETE FROM messages WHERE message_id = ?`
	_, err := m.DB.Exec(stmt, messageID)
	if err != nil {
		return err
	}
	return nil
}