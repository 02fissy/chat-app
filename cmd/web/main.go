package main

import (
	"database/sql"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
	"chat.fisayo.net/internal/database"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)
var addr = flag.String("addr", ":8000", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
	 	return
	 }
	 if r.Method != http.MethodGet {
	 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	 	return
	 }
	 http.ServeFile(w, r, "cmd/web/home.html")
}
func main() {
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
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	logger.Info("starting server", "addr", *addr)

	err = http.ListenAndServe(*addr, nil)
	logger.Error(err.Error())

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 12 * time.Hour
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
	// schema, err := os.ReadFile("internal/database/chat.sql")
	// if err != nil {
	// 	db.Close()
	//  	return nil, err
	//  }
	//  _, err = db.Exec(string(schema))
	//  if err != nil {
	// 	db.Close()
	// 	return nil, err
	//  }

	return db, nil

}

