package web

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"snippetbox/pkg/models/postgres"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type app struct {
	errLog   *log.Logger
	infoLog  *log.Logger
	snippets *postgres.SnippetModel
}


func Run() {
	addr := flag.String("addr", ":4000", "HTTP network address")

	err := godotenv.Load(".env")

	var (
		host     = os.Getenv("POSTGRES_HOST")
		port     = os.Getenv("POSTGRES_PORT")
		user     = os.Getenv("POSTGRES_USER")
		password = os.Getenv("POSTGRES_PASSWORD")
		dbname   = os.Getenv("POSTGRES_DB")
	)

	connStr := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	flag.Parse()

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
