package models

import(
	"database/sql"
)

type Participants struct {
	ParticipantID int
	Role string
	ConversationID int
	UserID int
	LastReadMessageID int
}
type ParticipantsModel struct{
	DB *sql.DB
}

func(m *ParticipantsModel) Insert(role string, conversationID, userID int) error {
	stmt := `INSERT INTO participants(role, conversation_id, user_id) VALUES (?, ?, ?)`
	_, err := m.DB.Exec(stmt, role, conversationID, userID)	
	if err != nil {
		return err
	}
	return nil
}
func(m *ParticipantsModel) GetParticipants(conversationID int) ([]int, error) {
	return nil, nil
}
func(m *ParticipantsModel) RemoveParticipants(participantID int) error {
	return nil
}