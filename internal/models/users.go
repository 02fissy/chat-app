package models

import (
	"database/sql" 
	"errors"
	"strings"
	"golang.org/x/crypto/bcrypt"
	"modernc.org/sqlite"
)
type Users struct {
	ID int
	Username string
	PhoneNumber string
	PasswordHash  []byte
}
type Attachment struct {
	AttachmentID int
	MessageID int
	FileURL string
	FileType string
}

type UserModel struct {
	DB *sql.DB
}
 func(m *UserModel) Insert(name, phone_no, password_hash string) error{
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password_hash), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users(username, phone_number, password_hash) 
	VALUES (?, ?, ?)`
	_, err = m.DB.Exec(stmt, name, phone_no, passwordHash)
	if err != nil{
		   var sqliteError *sqlite.Error
        if errors.As(err, &sqliteError) {
            if strings.Contains(sqliteError.Error(), "users_uc_phone") {
                return ErrDuplicatePhone
            }
        }
		return err
	}
	return nil
 }
 func(m *UserModel) Authenticate(name, phone_no, password_hash string) (int64, error) {
	var id int64
	var passwordHash []byte
	stmt := `SELECT user_id, password_hash FROM users WHERE username = ? AND phone_number = ?`
	err := m.DB.QueryRow(stmt, name, phone_no).Scan(&id, &passwordHash)
	if err != nil{
		if errors.Is(err, sql.ErrNoRows) {
            return 0, ErrInvalidCredentials
        } else {
            return 0, err
        }
	}
	err = bcrypt.CompareHashAndPassword(passwordHash, []byte(password_hash))
    if err != nil {
        if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
            return 0, ErrInvalidCredentials
        } else {
            return 0, err
        }
    }
	
	return id, nil
 }
func (m *UserModel) Exists(id int) (bool, error) {

var exists bool
stmt := "SELECT EXISTS(SELECT true FROM users WHERE user_id = ?)"
err := m.DB.QueryRow(stmt, id).Scan(&exists)
return exists, err
}

//Attachments
type AttachmentModel struct{
	DB *sql.DB
}

func(m *AttachmentModel) Insert(messageID int, fileURL, fileType string) error {
	stmt := `INSERT INTO attachments(message_id, file_url, file_type) 
	VALUES (?, ?, ?)`
	_, err := m.DB.Exec(stmt, messageID, fileURL, fileType)
	return err
}
func(m *AttachmentModel) GetAttachments(messageID int) ([]Attachment, error) {
	stmt := `SELECT attachment_id, message_id, file_url, file_type FROM attachments WHERE message_id = ?`
	rows, err := m.DB.Query(stmt, messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var attachments []Attachment

	for rows.Next() {
		var a Attachment
		if err := rows.Scan(&a.AttachmentID, &a.MessageID, &a.FileURL, &a.FileType); err != nil {
			return nil, err
		}
		attachments = append(attachments, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return attachments, nil
}
func(m *AttachmentModel) DeleteAttachments(messageID int) error {
	stmt := `DELETE FROM attachments WHERE message_id = ?`
	_, err := m.DB.Exec(stmt, messageID)
	if err !=nil {
		return err
	}
	return nil
}