package models

import(
	"database/sql"
)

type Contact struct {
	ContactID int
	User1ID int
	User2ID int
}

type ContactModel struct{
	DB *sql.DB
}

func(m *ContactModel) Add(user_id1, user_id2 int) error {
	stmt := `INSERT INTO contacts(user1_id, user2_id) VALUES (?, ?)`
	_, err := m.DB.Exec(stmt, user_id1, user_id2)
	if err != nil {
		return err
	}
	return nil
}
func(m *ContactModel) GetContacts(user_id int) ([]int, error) {
	stmt:= `SELECT user2_id FROM contacts WHERE user1_id = ?`
	rows, err := m.DB.Query(stmt, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []int
	for rows.Next() {
		var contactID int
		if err := rows.Scan(&contactID); err != nil {
			return nil, err
		}
		contacts = append(contacts, contactID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}
func(m *ContactModel) RemoveContacts(contact_id, user_id2 int) error {
	stmt := `DELETE FROM contacts WHERE contact_id = ? AND user2_id = ?`
	_, err := m.DB.Exec(stmt, contact_id, user_id2)
	if err != nil {
		return err
	}
	return nil
}