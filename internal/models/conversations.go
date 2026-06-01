package models

import (
	"database/sql"
)

type Conversations struct {
	ConversationID int
	TypeOf string
	Title string
	ParentID int
}
type ConversationsModel struct{
	DB *sql.DB
}
func(m *ConversationsModel) Insert(typeOf, title string, parentID int) error {
	stmt := `INSERT INTO conversations(type_of, title, parent_id) VALUES (?, ?, ?)`
	_, err := m.DB.Exec(stmt, typeOf, title, parentID)
	if err != nil {
		return err
	}
	return nil
}
func(m *ConversationsModel) GetByID(conversationID int) (*Conversations, error) {
	stmt := `SELECT conversation_id, type_of, title, parent_id FROM conversations WHERE conversation_id = ?`
	row := m.DB.QueryRow(stmt, conversationID)
	var c Conversations
	err := row.Scan(&c.ConversationID, &c.TypeOf, &c.Title, &c.ParentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}
func(m *ConversationsModel) DeleteConversations(conversationID int) error{
	stmt:= `DELETE FROM conversations WHERE conversation_id = ?`
	_, err := m.DB.Exec(stmt, conversationID)	
	if err != nil {
		return err
	}	
	return nil
}