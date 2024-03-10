package main

import (
	"flag"
	"fmt"
	"github.com/Vadim992/clinicAPI/internal/clinicapi"
	"github.com/Vadim992/clinicAPI/internal/config"
	"github.com/Vadim992/clinicAPI/internal/database/postgres"
	"github.com/Vadim992/clinicAPI/pkg/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("cannot set env variables (err or have no file \".env\"): %e", err)
	}

	cfg := config.NewConfig()
	cfg.SetFromEnv()

	conn := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable search_path=%s",
		cfg.HostDB, cfg.PortDB, cfg.UsernameDB, cfg.PasswordDB, cfg.DbNameDB, cfg.SearchPathDB)

	addr := flag.String("addr", ":3000", "HTTP network address")

	flag.Parse()

	db, err := postgres.InitDB(conn)
	if err != nil {
		logger.ErrLog.Fatalf("cannot connect to database: %v", err)
	}

	postgres.DataBase.DB = db

	srv := &http.Server{
		Addr:     *addr,
		Handler:  clinicapi.Routes(),
		ErrorLog: logger.ErrLog,
	}

	logger.InfoLog.Printf("Starting server on port %s\n", *addr)

	logger.ErrLog.Fatal(srv.ListenAndServe())
}
