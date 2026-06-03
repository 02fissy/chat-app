package models

import (
	"testing"
	"os"
	"database/sql"
    "golang.org/x/crypto/bcrypt"
    "chat.fisayo.net/internal/assert"
    "github.com/pressly/goose/v3"
)
func newTestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "chat_test_*.db")
	if err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("sqlite", tmpFile.Name())
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		db.Close()
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		t.Fatal(err)
	}

	err = goose.SetDialect("sqlite")
	if err != nil {
		db.Close()
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		t.Fatal(err)
	}

	err = goose.Up(db, "../database/migrations")
	if err != nil {
		db.Close()
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		t.Fatal(err)
	}

	cleanup := func() {
		db.Close()
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}

	return db, cleanup
}
func TestUserInsert(t *testing.T) {
	tests := []struct {
		name      string
		username  string
		phone     string
		password  string
		wantError bool
	}{
		{
			name:      "valid user",
			username:  "john",
			phone:     "+2348012345678",
			password:  "password123",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, cleanup := newTestDB(t)
			defer cleanup()

			m := UserModel{DB: db}

			err := m.Insert(
				tt.username,
				tt.phone,
				tt.password,
			)

			if tt.wantError {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUserAuthenticate(t *testing.T) {

    db, cleanup := newTestDB(t)
    defer cleanup()

    model := UserModel{
        DB: db,
    }
    hash, err := bcrypt.GenerateFromPassword(
        []byte("password123"),
        bcrypt.DefaultCost,
    )

    if err != nil {
        t.Fatal(err)
    }
    _, err = db.Exec(`
        INSERT INTO users (
            username,
            phone_number,
            password_hash
        )
        VALUES (?, ?, ?)
    `,
        "fissy02",
        "+2348012345678",
        string(hash),
    )

    if err != nil {
        t.Fatal(err)
    }

    tests := []struct {
        name        string
        username    string
        phone       string
        password    string
        wantErr     bool
        wantUserID  int64
    }{
        {
            name:       "valid credentials",
            username:   "fissy02",
            phone:      "+2348012345678",
            password:   "password123",
            wantErr:    false,
            wantUserID: 1,
        },
        {
            name:       "wrong password",
            username:   "fissy02",
            phone:      "+2348012345678",
            password:   "wrongpassword",
            wantErr:    true,
            wantUserID: 0,
        },
        {
            name:       "non-existent user",
            username:   "ghost",
            phone:      "+2348012345679",
            password:   "password123",
            wantErr:    true,
            wantUserID: 0,
        },
    }

    for _, tt := range tests {

        t.Run(tt.name, func(t *testing.T) {

            userID, err := model.Authenticate(
                tt.username,
                tt.phone,
                tt.password,
            )

            if tt.wantErr {

                if err == nil {
                    t.Fatal("expected error but got nil")
                }

                return
            }

            if err != nil {
                t.Fatal(err)
            }

            if userID != tt.wantUserID {
                t.Fatalf(
                    "got %d want %d",
                    userID,
                    tt.wantUserID,
                )
            }
        })
    }
}
func TestUserModelExists(t *testing.T) {
    tests := []struct {
        name   string
        userID int
        want   bool
    }{
        {
            name:   "Valid ID",
            userID: 1,
            want:   true,
        },
        {
            name:   "Zero ID",
            userID: 0,
            want:   false,
        },
        {
            name:   "Non-existent ID",
            userID: 2,
            want:   false,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            db, cleanup := newTestDB(t)
            defer cleanup()
            m := UserModel{DB: db}
            err := m.Insert(
				"john",
				"+2348012345678",
				"password123",
			)
			if err != nil {
				t.Fatal(err)
			}

            exists, err := m.Exists(tt.userID)
            assert.Equal(t, exists, tt.want)
            assert.NilError(t, err)
        })
    }
}