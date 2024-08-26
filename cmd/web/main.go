package web

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"snippetbox/config"
	"snippetbox/pkg/models/postgres"

	_ "github.com/lib/pq"
)

type app struct {
	errLog   *log.Logger
	infoLog  *log.Logger
	snippets *postgres.SnippetModel
}

func Run() {
	// flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// config
	dbConfig := config.New()
	connStr := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName)

	infoLog := log.New(os.Stdout, "[INFO]\t", log.LstdFlags)
	errLog := log.New(os.Stderr, "[ERROR]\t", log.LstdFlags|log.Lshortfile)

	db, err := openDB(connStr)
	if err == nil {
		infoLog.Println("Connected to DB successfully!")
	} else {
		errLog.Fatal(err)
	}
	defer db.Close()

	app := &app{
		infoLog:  infoLog,
		errLog:   errLog,
		snippets: &postgres.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)

	err = srv.ListenAndServe()
	errLog.Fatal(err)
}

func openDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
