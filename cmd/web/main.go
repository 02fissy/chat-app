package main

import (
	"database/sql"
	"log/slog"
	"os"
	"time"
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/sqlite3store"
	_ "modernc.org/sqlite"

)

func main() {
	logger:= slog.New(slog.NewTextHandler(os.Stdout, nil))
	dsn := "chatapp.db"
	db,err := openDB(dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}	
	defer db.Close()

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 12 * time.Hour
}
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	schema, err := os.ReadFile("internal/database/chat.sql")
	if err != nil {
		db.Close()
	 	return nil, err
	 }
	 _, err = db.Exec(string(schema))
	 if err != nil {
		db.Close()
		return nil, err
	 }

	return db, nil

}

