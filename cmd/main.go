package main

import (
	"flag"
	"fmt"
	"github.com/Vadim992/clinicAPI/internal/database/postgres"
	"log"
	"net/http"
	"os"
)

type clinicAPI struct {
	errLog  *log.Logger
	infoLog *log.Logger
	DB      *postgres.DB
}

// лучше сделать через флаги ?
const (
	host       = "localhost"
	username   = "clinicapi_user"
	password   = "mypass"
	dbname     = "clinicapi"
	port       = 5432
	sslmode    = "disable"
	searchPath = "clinicapi"
)

func main() {
	conn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s search_path=%s",
		host, port, username, password, dbname, sslmode, searchPath)

	addr := flag.String("addr", ":3000", "HTTP network address")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.LstdFlags)
	errLog := log.New(os.Stderr, "ERROR\t", log.LstdFlags|log.Lshortfile)

	db, err := postgres.ConnectDB(conn)

	if err != nil {
		errLog.Fatal(err)
	}

	c := &clinicAPI{
		errLog:  errLog,
		infoLog: infoLog,
		DB:      &postgres.DB{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		Handler:  c.routes(),
		ErrorLog: errLog,
	}

	infoLog.Printf("Starting server on port %s\n", *addr)

	errLog.Fatal(srv.ListenAndServe())

}
