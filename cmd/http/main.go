package main

import (
	"flag"
	"fmt"
	"github.com/Vadim992/clinicAPI/internal/clinicapi"
	"github.com/Vadim992/clinicAPI/pkg/database/postgres"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Set env variables err (or have no file .env): %e", err)
	}
}

func main() {

	host_db, port_db, username_db, password_db, dbname_db, searchPath_db := clinicapi.GetProjectEnv()

	conn := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable search_path=%s",
		host_db, port_db, username_db, password_db, dbname_db, searchPath_db)

	addr := flag.String("addr", ":3000", "HTTP network address")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.LstdFlags)
	errLog := log.New(os.Stderr, "ERROR\t", log.LstdFlags|log.Lshortfile)

	db, err := postgres.ConnectDB(conn)

	if err != nil {
		errLog.Fatal(err)
	}

	c := &clinicapi.ClinicAPI{
		ErrLog:  errLog,
		InfoLog: infoLog,
		DB:      &postgres.DB{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		Handler:  c.Routes(),
		ErrorLog: c.ErrLog,
	}

	infoLog.Printf("Starting server on port %s\n", *addr)

	errLog.Fatal(srv.ListenAndServe())

}
