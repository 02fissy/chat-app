package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"html/template"
	"time"
	"chat.fisayo.net/internal/database"
	"chat.fisayo.net/internal/models"
	"github.com/go-playground/form/v4"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)
type application struct{
	logger *slog.Logger
	templateCache map[string]*template.Template
	users models.UserModelInterface
	sessionManager *scs.SessionManager
	hub *Hub
	formDecoder *form.Decoder
}

func main() {
	addr := flag.String("addr", ":8000", "http service address")
	flag.Parse()
	logger:= slog.New(slog.NewTextHandler(os.Stdout, nil))
	dsn := "chatapp.db"
	db,err := openDB(dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}	
	defer db.Close()
	hub := newHub()
	go hub.run()
	
	
	 templateCache, err := newTemplateCache()
     if err != nil {
         logger.Error(err.Error())
         os.Exit(1)
     }
	
	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	
	formDecoder := form.NewDecoder()
	app := &application{
		logger:  logger,
		templateCache:  templateCache,
		sessionManager:  sessionManager,
		hub: hub,
		formDecoder: formDecoder,
		users: &models.UserModel{DB: db},
	}
		
	logger.Info("starting server", "addr", *addr)

	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())

}
func runMigrations(db *sql.DB) error {
	goose.SetBaseFS(database.MigrationsFS)

	if err:= goose.SetDialect("sqlite"); err != nil {
		return err
	}
	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}
	return nil
}
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	err = runMigrations(db)
	if err != nil {
		db.Close()
		return nil, err
	}
	db.SetMaxOpenConns(1)
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil

}

